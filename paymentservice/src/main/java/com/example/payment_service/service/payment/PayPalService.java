package com.example.payment_service.service.payment;

import com.paypal.api.payments.*;
import com.paypal.base.rest.APIContext;
import com.paypal.base.rest.PayPalRESTException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.ArrayList;
import java.util.List;

@Service
public class PayPalService implements IPayPalPayment {

  @Autowired
  private APIContext apiContext;

  @Value("${payment.paypal.successUrl}")
  private String successUrl;

  @Value("${payment.paypal.cancelUrl}")
  private String cancelUrl;

  public Payment createPayment(
      Double total,
      String currency,
      String method,
      String intent,
      String description,
      String orderId) throws PayPalRESTException {

    // Convert VND to USD if necessary (PayPal doesn't support VND)
    BigDecimal usdAmount = convertToUSD(total);
    System.out.println("VND to USD: " + usdAmount);

    Amount amount = new Amount();
    amount.setCurrency("USD");
    amount.setTotal(usdAmount.toString());

    Transaction transaction = new Transaction();
    transaction.setDescription(description);
    transaction.setAmount(amount);
    transaction.setInvoiceNumber(orderId);

    List<Transaction> transactions = new ArrayList<>();
    transactions.add(transaction);

    Payer payer = new Payer();
    payer.setPaymentMethod(method);

    Payment payment = new Payment();
    payment.setIntent(intent);
    payment.setPayer(payer);
    payment.setTransactions(transactions);

    RedirectUrls redirectUrls = new RedirectUrls();
    redirectUrls.setCancelUrl(cancelUrl);
    redirectUrls.setReturnUrl(successUrl + "?orderId=" + orderId);
    payment.setRedirectUrls(redirectUrls);

    return payment.create(apiContext);
  }

  public Payment executePayment(String paymentId, String payerId) throws PayPalRESTException {
    Payment payment = new Payment();
    payment.setId(paymentId);

    PaymentExecution paymentExecution = new PaymentExecution();
    paymentExecution.setPayerId(payerId);

    return payment.execute(apiContext, paymentExecution);
  }

  private BigDecimal convertToUSD(Double vndAmount) {
    // Use a proper exchange rate API in production
    // This is just a simplified example with a fixed rate
    double exchangeRate = 0.000042; // Example rate: 1 VND = 0.000042 USD
    return new BigDecimal(vndAmount * exchangeRate)
        .setScale(2, RoundingMode.HALF_UP);
  }
}
