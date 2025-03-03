package com.example.payment_service.entity;

import jakarta.persistence.*;
import lombok.Data;

import java.time.LocalDateTime;

@Entity
@Table(name = "payments")
@Data
public class Payment {
  @Id
  @GeneratedValue(strategy = GenerationType.IDENTITY)
  private Long id;

  private String orderId;
  private Double amount;
  private String currency;
  private String status;
  private String paymentMethod;
  private String transactionId;
  private String payerId;
  private LocalDateTime createdAt;
  private LocalDateTime updatedAt;
}
