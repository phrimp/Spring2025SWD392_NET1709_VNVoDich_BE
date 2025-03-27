package com.example.payment_service.controller;

import com.example.payment_service.entity.Payment;
import com.example.payment_service.repository.PaymentRepo;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpStatus;

import java.io.IOException;
import java.io.InputStream;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.HashMap;

@RestController
@RequestMapping("/api/payment/admin")
public class PaymentImportController {

  @Autowired
  private PaymentRepo paymentRepo;

  /**
   * Import payments from a JSON file uploaded by the user
   */
  @PostMapping("/import")
  public ResponseEntity<Map<String, Object>> importPaymentsFromFile(
      @RequestParam("file") MultipartFile file,
      @RequestHeader("API_KEY") String apiKey) {

    Map<String, Object> response = new HashMap<>();

    try {
      // Validate API key
      if (!apiKey.equals(System.getenv("API_KEY"))) {
        response.put("error", "Unauthorized: Invalid API key");
        return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
      }

      // Configure ObjectMapper with JavaTimeModule for LocalDateTime support
      ObjectMapper mapper = new ObjectMapper();
      mapper.registerModule(new JavaTimeModule());
      mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);

      // Read payment data from uploaded file
      List<Payment> payments = mapper.readValue(
          file.getInputStream(),
          new TypeReference<List<Payment>>() {
          });

      // Save all payments to the database
      List<Payment> savedPayments = paymentRepo.saveAll(payments);

      response.put("status", "success");
      response.put("message", "Successfully imported " + savedPayments.size() + " payment records");
      response.put("count", savedPayments.size());

      return ResponseEntity.ok(response);

    } catch (IOException e) {
      response.put("status", "error");
      response.put("message", "Failed to import payment data: " + e.getMessage());
      return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }
  }

  /**
   * Import payments from the default payment.json file in the resources folder
   */
  @PostMapping("/import/default")
  public ResponseEntity<Map<String, Object>> importDefaultPayments(
      @RequestHeader("API_KEY") String apiKey) {

    Map<String, Object> response = new HashMap<>();

    try {
      // Validate API key
      if (!apiKey.equals(System.getenv("API_KEY"))) {
        response.put("error", "Unauthorized: Invalid API key");
        return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
      }

      // Configure ObjectMapper with JavaTimeModule for LocalDateTime support
      ObjectMapper mapper = new ObjectMapper();
      mapper.registerModule(new JavaTimeModule());
      mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);

      // Load the JSON file from resources
      Resource resource = new ClassPathResource("payment.json");

      try (InputStream inputStream = resource.getInputStream()) {
        // Read payment data from JSON file
        List<Payment> payments = mapper.readValue(
            inputStream,
            new TypeReference<List<Payment>>() {
            });

        // Save all payments to the database
        List<Payment> savedPayments = paymentRepo.saveAll(payments);

        response.put("status", "success");
        response.put("message",
            "Successfully imported " + savedPayments.size() + " payment records from default payment.json");
        response.put("count", savedPayments.size());

        return ResponseEntity.ok(response);
      }

    } catch (IOException e) {
      response.put("status", "error");
      response.put("message", "Failed to import payment data: " + e.getMessage());
      return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }
  }

  /**
   * Clear all payment records from the database
   */
  @PostMapping("/clear")
  public ResponseEntity<Map<String, Object>> clearAllPayments(
      @RequestHeader("API_KEY") String apiKey) {

    Map<String, Object> response = new HashMap<>();

    // Validate API key
    if (!apiKey.equals(System.getenv("API_KEY"))) {
      response.put("error", "Unauthorized: Invalid API key");
      return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
    }

    try {
      long count = paymentRepo.count();
      paymentRepo.deleteAll();

      response.put("status", "success");
      response.put("message", "Successfully deleted all " + count + " payment records");

      return ResponseEntity.ok(response);

    } catch (Exception e) {
      response.put("status", "error");
      response.put("message", "Failed to clear payment data: " + e.getMessage());
      return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }
  }

  /**
   * Add a single payment using request body
   */
  @PostMapping("/add")
  public ResponseEntity<Map<String, Object>> addPayment(
      @RequestBody Payment payment,
      @RequestHeader("API_KEY") String apiKey) {

    Map<String, Object> response = new HashMap<>();

    // Validate API key
    if (!apiKey.equals(System.getenv("API_KEY"))) {
      response.put("error", "Unauthorized: Invalid API key");
      return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
    }

    try {
      // Set timestamps if not provided
      if (payment.getCreatedAt() == null) {
        payment.setCreatedAt(LocalDateTime.now());
      }
      if (payment.getUpdatedAt() == null) {
        payment.setUpdatedAt(LocalDateTime.now());
      }

      // Save the payment to the database
      Payment savedPayment = paymentRepo.save(payment);

      response.put("status", "success");
      response.put("message", "Payment added successfully");
      response.put("payment", savedPayment);

      return ResponseEntity.ok(response);

    } catch (Exception e) {
      response.put("status", "error");
      response.put("message", "Failed to add payment: " + e.getMessage());
      return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }
  }

  /**
   * Add multiple payments using request body
   */
  @PostMapping("/add/batch")
  public ResponseEntity<Map<String, Object>> addPayments(
      @RequestBody List<Payment> payments,
      @RequestHeader("API_KEY") String apiKey) {

    Map<String, Object> response = new HashMap<>();

    // Validate API key
    if (!apiKey.equals(System.getenv("API_KEY"))) {
      response.put("error", "Unauthorized: Invalid API key");
      return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
    }

    try {
      // Set timestamps if not provided
      LocalDateTime now = LocalDateTime.now();
      for (Payment payment : payments) {
        if (payment.getCreatedAt() == null) {
          payment.setCreatedAt(now);
        }
        if (payment.getUpdatedAt() == null) {
          payment.setUpdatedAt(now);
        }
      }

      // Save all payments to the database
      List<Payment> savedPayments = paymentRepo.saveAll(payments);

      response.put("status", "success");
      response.put("message", "Successfully added " + savedPayments.size() + " payment records");
      response.put("count", savedPayments.size());

      return ResponseEntity.ok(response);

    } catch (Exception e) {
      response.put("status", "error");
      response.put("message", "Failed to add payments: " + e.getMessage());
      return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }
  }
}
