package com.example.payment_service.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.example.payment_service.entity.Payment;
import com.example.payment_service.service.payment.PaymentService;

@RestController
@RequestMapping("/api/payment")
public class PaymentController {

  @Autowired
  private PaymentService paymentService;

  @GetMapping("/order/{order_id}")
  public Payment getPaymentByOrderID(@PathVariable String order_id) {
    return paymentService.getPaymentByOrderID(order_id);
  }

}
