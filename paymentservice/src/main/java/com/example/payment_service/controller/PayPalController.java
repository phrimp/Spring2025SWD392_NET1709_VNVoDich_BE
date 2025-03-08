// com/example/payment_service/controller/PayPalController.java
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

@RestController
@RequestMapping("/api/payment/paypal")
public class PayPalController {

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
      @RequestHeader("API_KEY") String api_key) {

    try {
      System.out.println("DEBUGGGGGGG" + api_key);
      com.paypal.api.payments.Payment payment = payPalService.createPayment(
          amount,
          "USD",
          "paypal",
          "sale",
          description,
          orderId);

      System.out.println("DEBUGGGGGGG22222");
      System.out.println(payment);
      // Temp Save Payment
      Payment paymentRecord = new Payment();
      paymentRecord.setOrderId(orderId);
      paymentRecord.setAmount(amount);
      paymentRecord.setStatus("PENDING");
      paymentRecord.setPaymentMethod("PAYPAL");
      paymentRecord.setCreatedAt(LocalDateTime.now());
      paymentRepo.save(paymentRecord);

      // Get approval URL to redirect user
      for (Links link : payment.getLinks()) {
        if (link.getRel().equals("approval_url")) {
          Map<String, String> response = new HashMap<>();
          response.put("redirectUrl", link.getHref());
          response.put("paymentId", payment.getId());

          return ResponseEntity.ok(response);
        }
      }

      return ResponseEntity.badRequest().body("Failed to get PayPal approval URL");

    } catch (PayPalRESTException e) {
      return ResponseEntity.badRequest().body(e.getMessage());
    }
  }

  @GetMapping("/success")
  public ResponseEntity<?> completePayment(
      @RequestParam("paymentId") String paymentId,
      @RequestParam("PayerID") String payerId,
      @RequestParam("orderId") String orderId) {

    try {
      com.paypal.api.payments.Payment payment = payPalService.executePayment(paymentId, payerId);

      if (payment.getState().equals("approved")) {
        // Update payment record
        Payment paymentRecord = paymentRepo.findByOrderId(orderId)
            .orElseThrow(() -> new RuntimeException("Payment record not found"));

        paymentRecord.setStatus("COMPLETED");
        paymentRecord.setTransactionId(paymentId);
        paymentRecord.setPayerId(payerId);
        paymentRecord.setUpdatedAt(LocalDateTime.now());
        paymentRepo.save(paymentRecord);

        // Send webhook notification
        webhookService.sendPaymentEvent("payment.completed", orderId, "COMPLETED");

        // Return success response
        Map<String, Object> response = new HashMap<>();
        response.put("status", "success");
        response.put("paymentId", paymentId);
        response.put("orderId", orderId);

        return ResponseEntity.ok(response);
      }

      return ResponseEntity.badRequest().body("Payment not approved");

    } catch (PayPalRESTException e) {
      // Send webhook notification about failure
      webhookService.sendPaymentEvent("payment.failed", orderId, "FAILED");
      return ResponseEntity.badRequest().body(e.getMessage());
    }
  }

  @GetMapping("/cancel")
  public ResponseEntity<?> cancelPayment(@RequestParam("orderId") String orderId) {
    // Update payment record to cancelled
    Payment paymentRecord = paymentRepo.findByOrderId(orderId)
        .orElseThrow(() -> new RuntimeException("Payment record not found"));

    paymentRecord.setStatus("CANCELLED");
    paymentRecord.setUpdatedAt(LocalDateTime.now());
    paymentRepo.save(paymentRecord);

    // Send webhook notification
    webhookService.sendPaymentEvent("payment.canceled", orderId, "CANCELLED");

    Map<String, String> response = new HashMap<>();
    response.put("status", "cancelled");
    response.put("message", "Payment was cancelled");

    return ResponseEntity.ok(response);
  }
}
