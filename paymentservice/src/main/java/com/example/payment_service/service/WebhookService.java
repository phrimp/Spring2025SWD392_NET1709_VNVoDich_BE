package com.example.payment_service.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.Map;

@Service
public class WebhookService {

  @Autowired
  private RestTemplate restTemplate;

  @Value("${webhook.subscription.url}")
  private String subscriptionWebhookUrl;

  @Value("${api.key}")
  private String apiKey;

  public void sendPaymentEvent(String event, String orderId, String status) {
    Map<String, String> payload = new HashMap<>();
    payload.put("event", event);
    payload.put("orderId", orderId);
    payload.put("status", status);

    HttpHeaders headers = new HttpHeaders();
    headers.set("API_KEY", apiKey);
    headers.set("Content-Type", "application/json");

    HttpEntity<Map<String, String>> request = new HttpEntity<>(payload, headers);

    try {
      restTemplate.postForEntity(subscriptionWebhookUrl, request, String.class);
    } catch (Exception e) {
      System.err.println("Failed to send webhook: " + e.getMessage());
    }
  }
}
