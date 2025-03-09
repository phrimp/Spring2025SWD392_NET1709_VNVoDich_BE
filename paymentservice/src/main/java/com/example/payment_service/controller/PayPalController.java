package com.example.payment_service.controller;

import com.example.payment_service.entity.Payment;
import com.example.payment_service.repository.PaymentRepo;
import com.example.payment_service.service.WebhookService;
import com.example.payment_service.service.payment.PayPalService;
import com.paypal.api.payments.Links;
import com.paypal.base.rest.PayPalRESTException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;

@RestController
@RequestMapping("/api/payment/paypal")
public class PayPalController {

  private static final Logger logger = Logger.getLogger(PayPalController.class.getName());

  @Autowired
  private PayPalService payPalService;

  @Autowired
  private PaymentRepo paymentRepo;

  @Autowired
  private WebhookService webhookService;

  @PostMapping("/create")
  public ResponseEntity<?> createPayment(
      @RequestParam("amount") Double amount,
      @RequestParam("description") String description,
      @RequestParam("orderId") String orderId,
      @RequestHeader("API_KEY") String apiKey) {

    logger.info("Creating payment for order: " + orderId + ", amount: " + amount + ", description: " + description);

    try {
      if (amount == null || amount <= 0) {
        return ResponseEntity.badRequest().body("Invalid amount value");
      }

      if (description == null || description.isEmpty()) {
        description = "Payment for order: " + orderId;
      }

      com.paypal.api.payments.Payment payment = payPalService.createPayment(
          amount,
          "USD",
          "paypal",
          "sale",
          description,
          orderId);

      // Temp Save Payment
      Payment paymentRecord = new Payment();
      paymentRecord.setOrderId(orderId);
      paymentRecord.setAmount(amount);
      paymentRecord.setStatus("PENDING");
      paymentRecord.setPaymentMethod("PAYPAL");
      paymentRecord.setCreatedAt(LocalDateTime.now());
      paymentRepo.save(paymentRecord);

      // Get approval URL to redirect user
      String approvalUrl = payment.getLinks().stream()
          .filter(link -> "approval_url".equals(link.getRel()))
          .findFirst()
          .map(Links::getHref)
          .orElse(null);

      if (approvalUrl != null) {
        Map<String, String> response = new HashMap<>();
        response.put("redirectUrl", approvalUrl);
        response.put("paymentId", payment.getId());
        logger.info("Payment created successfully for order: " + orderId + ", redirecting to: " + approvalUrl);
        return ResponseEntity.ok(response);
      }

      logger.warning("Failed to get PayPal approval URL for order: " + orderId);
      return ResponseEntity.badRequest().body("Failed to get PayPal approval URL");

    } catch (PayPalRESTException e) {
      logger.log(Level.SEVERE, "PayPal error for order: " + orderId, e);
      return ResponseEntity.badRequest().body(e.getMessage());
    } catch (Exception e) {
      logger.log(Level.SEVERE, "Unexpected error for order: " + orderId, e);
      return ResponseEntity.badRequest().body("An unexpected error occurred: " + e.getMessage());
    }
  }

  @GetMapping("/success")
  public ResponseEntity<?> completePayment(
      @RequestParam("paymentId") String paymentId,
      @RequestParam("PayerID") String payerId,
      @RequestParam("orderId") String orderId) {

    logger.info("Completing payment for order: " + orderId + ", paymentId: " + paymentId + ", payerId: " + payerId);

    try {
      com.paypal.api.payments.Payment payment = payPalService.executePayment(paymentId, payerId);

      if (payment.getState().equals("approved")) {
        // Update payment record
        Payment paymentRecord = paymentRepo.findByOrderId(orderId)
            .orElseThrow(() -> new RuntimeException("Payment record not found for order: " + orderId));

        paymentRecord.setStatus("COMPLETED");
        paymentRecord.setTransactionId(paymentId);
        paymentRecord.setPayerId(payerId);
        paymentRecord.setUpdatedAt(LocalDateTime.now());
        paymentRepo.save(paymentRecord);

        // Send webhook notification
        webhookService.sendPaymentEvent("payment.completed", orderId, "COMPLETED");
        logger.info("Payment completed successfully for order: " + orderId);

        // Return success response
        Map<String, Object> response = new HashMap<>();
        response.put("status", "success");
        response.put("paymentId", paymentId);
        response.put("orderId", orderId);

        return ResponseEntity.ok(response);
      }

      logger.warning("Payment not approved for order: " + orderId + ", state: " + payment.getState());
      return ResponseEntity.badRequest().body("Payment not approved");

    } catch (PayPalRESTException e) {
      logger.log(Level.SEVERE, "PayPal error completing payment for order: " + orderId, e);
      // Send webhook notification about failure
      webhookService.sendPaymentEvent("payment.failed", orderId, "FAILED");
      return ResponseEntity.badRequest().body(e.getMessage());
    } catch (Exception e) {
      logger.log(Level.SEVERE, "Unexpected error completing payment for order: " + orderId, e);
      webhookService.sendPaymentEvent("payment.failed", orderId, "FAILED");
      return ResponseEntity.badRequest().body("An unexpected error occurred: " + e.getMessage());
    }
  }

  @GetMapping("/cancel")
  public ResponseEntity<?> cancelPayment(@RequestParam("orderId") String orderId) {
    logger.info("Cancelling payment for order: " + orderId);
    
    try {
      // Update payment record to cancelled
      Payment paymentRecord = paymentRepo.findByOrderId(orderId)
          .orElseThrow(() -> new RuntimeException("Payment record not found for order: " + orderId));

      paymentRecord.setStatus("CANCELLED");
      paymentRecord.setUpdatedAt(LocalDateTime.now());
      paymentRepo.save(paymentRecord);

      // Send webhook notification
      webhookService.sendPaymentEvent("payment.canceled", orderId, "CANCELLED");
      logger.info("Payment cancelled for order: " + orderId);

      Map<String, String> response = new HashMap<>();
      response.put("status", "cancelled");
      response.put("message", "Payment was cancelled");

      return ResponseEntity.ok(response);
    } catch (Exception e) {
      logger.log(Level.SEVERE, "Error cancelling payment for order: " + orderId, e);
      return ResponseEntity.badRequest().body("Error cancelling payment: " + e.getMessage());
    }
  }
}
