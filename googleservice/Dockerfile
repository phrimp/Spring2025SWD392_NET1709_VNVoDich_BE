FROM golang:1.23-alpine

WORKDIR /app

# Install required dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy .env file (optional if required during build)
COPY .env .env

# Build the application
RUN go build -o main ./cmd

# Run the application
CMD ["./main"]
