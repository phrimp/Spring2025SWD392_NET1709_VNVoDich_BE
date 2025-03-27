package com.example.payment_service.config;

import com.example.payment_service.entity.Payment;
import com.example.payment_service.repository.PaymentRepo;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.orm.ObjectOptimisticLockingFailureException;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;

@Configuration
public class DataLoader {

  private static final Logger logger = LoggerFactory.getLogger(DataLoader.class);

  @Autowired
  private PaymentRepo paymentRepo;

  // In DataLoader.java
  @Bean
  public CommandLineRunner loadData() {
    return args -> {
      // Check if we already have payment data
      if (paymentRepo.count() > 0) {
        logger.info("Database already has payment data, skipping initialization");
        return;
      }

      logger.info("Initializing payment data from payment.json");

      try {
        // Configure ObjectMapper with JavaTimeModule for LocalDateTime support
        ObjectMapper objectMapper = new ObjectMapper();
        objectMapper.registerModule(new JavaTimeModule());
        objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        // In DataLoader.java - Add this to the ObjectMapper configuration
        objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        // Load the JSON file from resources
        ClassPathResource resource = new ClassPathResource("payment.json");

        try (InputStream inputStream = resource.getInputStream()) {
          // Read payment data from JSON file
          List<Payment> payments = objectMapper.readValue(
              inputStream,
              new TypeReference<List<Payment>>() {
              });

          // Save each payment individually with retry logic
          for (Payment payment : payments) {
            try {
              paymentRepo.save(payment);
            } catch (ObjectOptimisticLockingFailureException e) {
              logger.warn("Optimistic locking failed for payment {}, skipping", payment.getId());
            }
          }

          logger.info("Successfully loaded payment records");
        }
      } catch (IOException e) {
        logger.error("Failed to load payment data: {}", e.getMessage(), e);
      }
    };
  }
}
