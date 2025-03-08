package com.example.payment_service.service.payment;

import com.example.payment_service.entity.Payment;

public interface IPayment {

  default Payment getPaymentByOrderID(String orderId) {
    throw new UnsupportedOperationException("Not implemented");
  }
}
