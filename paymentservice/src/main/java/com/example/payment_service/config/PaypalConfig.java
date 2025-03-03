package com.example.payment_service.config;

import com.paypal.base.rest.APIContext;
import com.paypal.base.rest.OAuthTokenCredential;
import com.paypal.base.rest.PayPalRESTException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.HashMap;
import java.util.Map;

@Configuration
public class PaypalConfig {

  @Value("${payment.paypal.clientId}")
  private String clientId;

  @Value("${payment.paypal.clientSecret}")
  private String clientSecret;

  @Value("${payment.paypal.mode}")
  private String mode;

  @Bean
  public Map<String, String> paypalSdkConfig() {
    Map<String, String> sdkConfig = new HashMap<>();
    sdkConfig.put("mode", mode);
    return sdkConfig;
  }

  @Bean
  public OAuthTokenCredential oAuthTokenCredential() {
    return new OAuthTokenCredential(clientId, clientSecret, paypalSdkConfig());
  }

  @Bean
  public APIContext apiContext() throws PayPalRESTException {
    APIContext context = new APIContext(clientId, clientSecret, mode);
    context.setConfigurationMap(paypalSdkConfig());
    return context;
  }
}
