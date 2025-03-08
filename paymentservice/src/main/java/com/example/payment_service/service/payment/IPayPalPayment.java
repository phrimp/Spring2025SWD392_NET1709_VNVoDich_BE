package com.example.payment_service.service.payment;

import com.paypal.api.payments.Payment;
import com.paypal.base.rest.PayPalRESTException;

public interface IPayPalPayment {
  // Existing methods

  // PayPal specific methods
  default Payment createPayment(Double total, String currency, String method, String intent,
      String description, String orderId) throws PayPalRESTException {
    throw new UnsupportedOperationException("Not implemented");
  }

  default Payment executePayment(String paymentId, String payerId) throws PayPalRESTException {
    throw new UnsupportedOperationException("Not implemented");
  }
}
