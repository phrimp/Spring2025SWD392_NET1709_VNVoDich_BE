# Stage 1: Build
FROM maven:3.9.4-eclipse-temurin-21-alpine AS build

WORKDIR /app

# Copy pom.xml for dependency resolution first (better caching)
COPY pom.xml .
RUN mvn dependency:go-offline -B

# Copy source code and build
COPY src ./src
RUN mvn clean package -DskipTests

# Stage 2: Production
FROM eclipse-temurin:21-jre-alpine

WORKDIR /app
#Update TZ data
RUN apk add --no-cache tzdata
# Create non-root user
RUN addgroup -g 1001 -S javauser && \
    adduser -S javauser -G javauser -u 1001
# Copy the jar from build stage
COPY --from=build /app/target/*.jar app.jar

# Set ownership
RUN chown -R javauser:javauser /app

# Use non-root user
USER javauser

# Run the application
CMD ["java", "-jar", "app.jar"]
