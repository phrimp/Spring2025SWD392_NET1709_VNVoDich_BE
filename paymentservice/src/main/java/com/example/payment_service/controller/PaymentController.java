package com.example.payment_service.controller;

import com.example.payment_service.entity.Payment;
import com.example.payment_service.repository.PaymentRepo;
import com.example.payment_service.service.payment.PaymentService;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/payment")
public class PaymentController {

  @Autowired
  private PaymentRepo paymentRepo;

  @Autowired
  private PaymentService paymentService;

  @GetMapping("/order/{order_id}")
  public Payment getPaymentByOrderID(@PathVariable String order_id) {
    return paymentService.getPaymentByOrderID(order_id);
  }

  @GetMapping("/all")
  public ResponseEntity<List<Payment>> getAllPayments() {
    List<Payment> payments = paymentRepo.findAll();
    return ResponseEntity.ok(payments);
  }

}
