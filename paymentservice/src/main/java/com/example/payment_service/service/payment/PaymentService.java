package com.example.payment_service.service.payment;

import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.example.payment_service.entity.Payment;
import com.example.payment_service.repository.PaymentRepo;

@Service
public class PaymentService implements IPayment {
  @Autowired
  private PaymentRepo paymentRepo;

  @Override
  public Payment getPaymentByOrderID(String orderId) {
    Optional<Payment> payment = paymentRepo.findByOrderId(orderId);
    if (payment == null) {
      System.out.println("Payment not found with order ID: " + orderId);
      return null;
    }
    return payment.get();
  }

}
