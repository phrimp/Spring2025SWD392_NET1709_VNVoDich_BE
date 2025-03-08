package services

import (
	"encoding/json"
	"fmt"
	"subscription/utils"
	"time"

	"github.com/valyala/fasthttp"
)

type PaymentService interface {
	CreatePayment(tutorID uint, planID uint, amount float64, billingCycle string) (string, string, error)
	GetPaymentByOrderID(orderID string) (*PaymentInfo, error)
	ValidatePaymentStatus(orderID string) (bool, error)
}

// PaymentInfo represents payment information returned from payment service
type PaymentInfo struct {
	OrderID       string    `json:"order_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type paymentService struct {
	paymentServiceURL string
	apiKey            string
}

// NewPaymentService creates a new instance of PaymentService
func NewPaymentService(paymentServiceURL, apiKey string) PaymentService {
	return &paymentService{
		paymentServiceURL: paymentServiceURL,
		apiKey:            apiKey,
	}
}

// CreatePayment creates a new payment for a subscription
func (s *paymentService) CreatePayment(tutorID uint, planID uint, amount float64, billingCycle string) (string, string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Generate order ID, format SUB-{tutorID}-{planID}-{timestamp}
	orderID := fmt.Sprintf("SUB-%d-%d-%d", tutorID, planID, time.Now().Unix())

	query := fmt.Sprintf("?amount=%f&description=%s&orderId=%s",
		amount,
		fmt.Sprintf("Subscription payment - %s plan (%s)", getPlanName(planID), billingCycle),
		orderID)

	// Send request to payment service
	utils.BuildRequest(req, "POST", nil, s.apiKey, fmt.Sprintf("%s/api/payment/paypal/create%s", s.paymentServiceURL, query))

	if err := fasthttp.Do(req, resp); err != nil {
		return "", "", fmt.Errorf("payment service unavailable: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return "", "", fmt.Errorf("failed to create payment: %s", resp.Body())
	}

	// Parse response to get redirect URL
	var response struct {
		RedirectURL string `json:"redirectUrl"`
		PaymentID   string `json:"paymentId"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Return the order ID and redirect URL
	return orderID, response.RedirectURL, nil
}

// GetPaymentByOrderID gets payment details by order ID
func (s *paymentService) GetPaymentByOrderID(orderID string) (*PaymentInfo, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Send request to payment service
	utils.BuildRequest(req, "GET", nil, s.apiKey, fmt.Sprintf("%s/api/payment/order/%s", s.paymentServiceURL, orderID))

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("payment service unavailable: %w", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("failed to get payment: %s", resp.Body())
	}

	// Parse response
	var payment PaymentInfo
	if err := json.Unmarshal(resp.Body(), &payment); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &payment, nil
}

// ValidatePaymentStatus checks if a payment has been completed successfully
func (s *paymentService) ValidatePaymentStatus(orderID string) (bool, error) {
	payment, err := s.GetPaymentByOrderID(orderID)
	if err != nil {
		return false, err
	}

	// Check if payment status is "COMPLETED"
	return payment.Status == "COMPLETED", nil
}

// Helper function to map plan IDs to names
func getPlanName(planID uint) string {
	planNames := map[uint]string{
		1: "Basic",
		2: "Premium",
		3: "Professional",
	}

	if name, ok := planNames[planID]; ok {
		return name
	}
	return "Unknown"
}
