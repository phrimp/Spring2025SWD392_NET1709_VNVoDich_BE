package services

import (
	"errors"
	"fmt"
	"subscription/internal/models"
	"subscription/internal/repository"
	"time"
)

type SubscriptionService interface {
	InitiateSubscription(req models.SubscriptionRequest) (*models.SubscriptionResponse, error)
	ConfirmSubscription(req models.PaymentConfirmationRequest) (*models.SubscriptionResponse, error)
	GetSubscriptionByID(id uint) (*models.SubscriptionResponse, error)
	GetTutorSubscription(tutorID uint) (*models.SubscriptionResponse, error)
	GetAllSubscriptions(page, pageSize int, filters map[string]interface{}) ([]models.SubscriptionResponse, int64, error)
	CancelSubscription(id uint) error
	ChangePlan(id uint, req models.ChangePlanRequest) (*models.SubscriptionResponse, error)
	UpdateSubscriptionStatus(id uint, status models.SubscriptionStatus) error
	ProcessPaymentWebhook(payload models.PaymentWebhookPayload) error
	GetExpiringSoonSubscriptions(days int) ([]models.SubscriptionResponse, error)
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
	planRepo         repository.PlanRepository
	paymentService   PaymentService
}

func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
	paymentService PaymentService,
) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		paymentService:   paymentService,
	}
}

// InitiateSubscription starts the subscription process and returns payment info
func (s *subscriptionService) InitiateSubscription(req models.SubscriptionRequest) (*models.SubscriptionResponse, error) {
	// Get the plan
	plan, err := s.planRepo.GetByID(req.PlanID)
	if err != nil {
		return nil, err
	}

	// Check if tutor already has an active subscription
	existingSub, err := s.subscriptionRepo.GetByTutorID(req.TutorID)
	if err == nil && existingSub != nil {
		// If there's an existing subscription but it's not active or it's canceled,
		// we can allow creating a new one
		if existingSub.Status != models.SubscriptionCanceled &&
			existingSub.Status != models.SubscriptionPastDue {
			return nil, errors.New("tutor already has an active subscription")
		}
	}

	// Calculate subscription period
	now := time.Now()
	var endDate time.Time
	var amount float64

	if req.BillingCycle == models.BillingMonthly {
		endDate = now.AddDate(0, 1, 0) // Add 1 month
		amount = plan.PriceMonthly
	} else {
		endDate = now.AddDate(1, 0, 0) // Add 1 year
		amount = plan.PriceAnnually
	}

	// Create payment order
	orderID := fmt.Sprintf("SUB-%d-%d-%d", req.TutorID, req.PlanID, time.Now().Unix())

	// Create subscription in pending state
	subscription := &models.TutorSubscription{
		TutorID:            req.TutorID,
		PlanID:             req.PlanID,
		Status:             models.SubscriptionIncomplete, // Start as incomplete until payment confirmed
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   endDate,
		CancelAtPeriodEnd:  false,
		BillingCycle:       req.BillingCycle,
		PaymentOrderID:     orderID, // Reference to payment order ID
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	// Log subscription creation event
	event := &models.SubscriptionEvent{
		SubscriptionID: subscription.ID,
		EventType:      "initiated",
		CurrentStatus:  models.SubscriptionIncomplete,
		Notes:          fmt.Sprintf("Subscription initiated with %s billing cycle", req.BillingCycle),
	}
	if err := s.subscriptionRepo.LogEvent(event); err != nil {
		// Just log the error but continue
		fmt.Printf("Failed to log subscription event: %v\n", err)
	}

	// Create response
	response := &models.SubscriptionResponse{
		ID:                 subscription.ID,
		TutorID:            subscription.TutorID,
		PlanName:           plan.Name,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		BillingCycle:       subscription.BillingCycle,
		Price:              amount,
		Features:           plan.FeaturesJSON,
		MaxCourses:         plan.MaxCourses,
		CommissionRate:     plan.CommissionRate,
		PaymentOrderID:     orderID,
	}

	return response, nil
}

// ConfirmSubscription finalizes a subscription after successful payment
func (s *subscriptionService) ConfirmSubscription(req models.PaymentConfirmationRequest) (*models.SubscriptionResponse, error) {
	// Find subscription by payment order ID
	var subscriptions []models.TutorSubscription
	var count int64
	filters := map[string]interface{}{
		"payment_order_id": req.OrderID,
	}

	subscriptions, count, err := s.subscriptionRepo.GetAll(1, 1, filters)
	if err != nil {
		return nil, err
	}

	if count == 0 || len(subscriptions) == 0 {
		return nil, errors.New("subscription not found for this payment order ID")
	}

	subscription := subscriptions[0]

	// Make sure it's not already active
	if subscription.Status == models.SubscriptionActive {
		plan, err := s.planRepo.GetByID(subscription.PlanID)
		if err != nil {
			return nil, err
		}

		var price float64
		if subscription.BillingCycle == models.BillingMonthly {
			price = plan.PriceMonthly
		} else {
			price = plan.PriceAnnually
		}

		return &models.SubscriptionResponse{
			ID:                 subscription.ID,
			TutorID:            subscription.TutorID,
			PlanName:           plan.Name,
			Status:             subscription.Status,
			CurrentPeriodStart: subscription.CurrentPeriodStart,
			CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
			CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
			BillingCycle:       subscription.BillingCycle,
			Price:              price,
			Features:           plan.FeaturesJSON,
			MaxCourses:         plan.MaxCourses,
			CommissionRate:     plan.CommissionRate,
		}, nil
	}

	// Update subscription status to active
	oldStatus := subscription.Status
	subscription.Status = models.SubscriptionActive

	if err := s.subscriptionRepo.Update(&subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// Log status change event
	event := &models.SubscriptionEvent{
		SubscriptionID: subscription.ID,
		EventType:      "payment_confirmed",
		PreviousStatus: oldStatus,
		CurrentStatus:  models.SubscriptionActive,
		Notes:          fmt.Sprintf("Payment confirmed with PaymentID: %s, PayerID: %s", req.PaymentID, req.PayerID),
	}

	if err := s.subscriptionRepo.LogEvent(event); err != nil {
		fmt.Printf("Failed to log subscription event: %v\n", err)
	}

	// Get plan details for response
	plan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil {
		return nil, err
	}

	// Determine price based on billing cycle
	var price float64
	if subscription.BillingCycle == models.BillingMonthly {
		price = plan.PriceMonthly
	} else {
		price = plan.PriceAnnually
	}

	// Create response
	response := &models.SubscriptionResponse{
		ID:                 subscription.ID,
		TutorID:            subscription.TutorID,
		PlanName:           plan.Name,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		BillingCycle:       subscription.BillingCycle,
		Price:              price,
		Features:           plan.FeaturesJSON,
		MaxCourses:         plan.MaxCourses,
		CommissionRate:     plan.CommissionRate,
	}

	return response, nil
}

// GetSubscriptionByID retrieves a subscription by its ID
func (s *subscriptionService) GetSubscriptionByID(id uint) (*models.SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var price float64
	if subscription.BillingCycle == models.BillingMonthly {
		price = subscription.Plan.PriceMonthly
	} else {
		price = subscription.Plan.PriceAnnually
	}

	return &models.SubscriptionResponse{
		ID:                 subscription.ID,
		TutorID:            subscription.TutorID,
		PlanName:           subscription.Plan.Name,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		BillingCycle:       subscription.BillingCycle,
		Price:              price,
		Features:           subscription.Plan.FeaturesJSON,
		MaxCourses:         subscription.Plan.MaxCourses,
		CommissionRate:     subscription.Plan.CommissionRate,
	}, nil
}

// GetTutorSubscription retrieves a tutor's active subscription
func (s *subscriptionService) GetTutorSubscription(tutorID uint) (*models.SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetByTutorID(tutorID)
	if err != nil {
		return nil, err
	}

	var price float64
	if subscription.BillingCycle == models.BillingMonthly {
		price = subscription.Plan.PriceMonthly
	} else {
		price = subscription.Plan.PriceAnnually
	}

	return &models.SubscriptionResponse{
		ID:                 subscription.ID,
		TutorID:            subscription.TutorID,
		PlanName:           subscription.Plan.Name,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		BillingCycle:       subscription.BillingCycle,
		Price:              price,
		Features:           subscription.Plan.FeaturesJSON,
		MaxCourses:         subscription.Plan.MaxCourses,
		CommissionRate:     subscription.Plan.CommissionRate,
	}, nil
}

// GetAllSubscriptions retrieves all subscriptions with pagination and filters
func (s *subscriptionService) GetAllSubscriptions(page, pageSize int, filters map[string]interface{}) ([]models.SubscriptionResponse, int64, error) {
	subscriptions, total, err := s.subscriptionRepo.GetAll(page, pageSize, filters)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		var price float64
		if subscription.BillingCycle == models.BillingMonthly {
			price = subscription.Plan.PriceMonthly
		} else {
			price = subscription.Plan.PriceAnnually
		}

		responses[i] = models.SubscriptionResponse{
			ID:                 subscription.ID,
			TutorID:            subscription.TutorID,
			PlanName:           subscription.Plan.Name,
			Status:             subscription.Status,
			CurrentPeriodStart: subscription.CurrentPeriodStart,
			CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
			CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
			BillingCycle:       subscription.BillingCycle,
			Price:              price,
			Features:           subscription.Plan.FeaturesJSON,
			MaxCourses:         subscription.Plan.MaxCourses,
			CommissionRate:     subscription.Plan.CommissionRate,
		}
	}

	return responses, total, nil
}

// CancelSubscription cancels a subscription
func (s *subscriptionService) CancelSubscription(id uint) error {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// If the subscription is still in trial or incomplete, cancel it immediately
	if subscription.Status == models.SubscriptionTrialing || subscription.Status == models.SubscriptionIncomplete {
		oldStatus := subscription.Status
		subscription.Status = models.SubscriptionCanceled

		if err := s.subscriptionRepo.Update(subscription); err != nil {
			return err
		}

		// Log cancellation event
		event := &models.SubscriptionEvent{
			SubscriptionID: subscription.ID,
			EventType:      "canceled",
			PreviousStatus: oldStatus,
			CurrentStatus:  models.SubscriptionCanceled,
			Notes:          "Subscription canceled immediately",
		}

		if err := s.subscriptionRepo.LogEvent(event); err != nil {
			fmt.Printf("Failed to log subscription event: %v\n", err)
		}

		return nil
	}

	// Otherwise, mark it to be canceled at the end of the current period
	if err := s.subscriptionRepo.SetCancelAtPeriodEnd(id, true); err != nil {
		return err
	}

	// Log cancellation at period end event
	event := &models.SubscriptionEvent{
		SubscriptionID: subscription.ID,
		EventType:      "cancel_scheduled",
		PreviousStatus: subscription.Status,
		CurrentStatus:  subscription.Status,
		Notes:          fmt.Sprintf("Subscription scheduled to cancel at the end of current period (%s)", subscription.CurrentPeriodEnd.Format("2006-01-02")),
	}

	if err := s.subscriptionRepo.LogEvent(event); err != nil {
		fmt.Printf("Failed to log subscription event: %v\n", err)
	}

	return nil
}

// ChangePlan changes a subscription's plan
func (s *subscriptionService) ChangePlan(id uint, req models.ChangePlanRequest) (*models.SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Only active subscriptions can be changed
	if subscription.Status != models.SubscriptionActive {
		return nil, errors.New("only active subscriptions can be changed")
	}

	// Get the new plan
	newPlan, err := s.planRepo.GetByID(req.NewPlanID)
	if err != nil {
		return nil, err
	}

	// Calculate the new end date and price
	var newEndDate time.Time
	var amount float64

	if req.BillingCycle == models.BillingMonthly {
		newEndDate = time.Now().AddDate(0, 1, 0) // Add 1 month
		amount = newPlan.PriceMonthly
	} else {
		newEndDate = time.Now().AddDate(1, 0, 0) // Add 1 year
		amount = newPlan.PriceAnnually
	}

	fmt.Println(newEndDate)

	// Create a new payment for the plan change
	orderID := fmt.Sprintf("SUB-%d-%d-%d", subscription.TutorID, req.NewPlanID, time.Now().Unix())

	// Store the old plan ID for event logging
	oldPlanID := subscription.PlanID
	oldBillingCycle := subscription.BillingCycle

	// The actual plan change will happen after payment confirmation
	// Update the subscription with new payment order ID
	subscription.PaymentOrderID = orderID

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// Log plan change initiated event
	event := &models.SubscriptionEvent{
		SubscriptionID: subscription.ID,
		EventType:      "plan_change_initiated",
		PreviousStatus: subscription.Status,
		CurrentStatus:  subscription.Status,
		Notes: fmt.Sprintf("Plan change initiated from ID %d to ID %d with billing cycle from %s to %s. Order ID: %s",
			oldPlanID, req.NewPlanID, req.BillingCycle, oldBillingCycle, orderID),
	}

	if err := s.subscriptionRepo.LogEvent(event); err != nil {
		fmt.Printf("Failed to log subscription event: %v\n", err)
	}

	// Create a response with the current plan details but including payment URL
	response := &models.SubscriptionResponse{
		ID:                 subscription.ID,
		TutorID:            subscription.TutorID,
		PlanName:           subscription.Plan.Name,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		BillingCycle:       subscription.BillingCycle,
		Price:              amount,
		Features:           subscription.Plan.FeaturesJSON,
		MaxCourses:         subscription.Plan.MaxCourses,
		CommissionRate:     subscription.Plan.CommissionRate,
		PaymentOrderID:     orderID,
	}

	return response, nil
}

// UpdateSubscriptionStatus changes a subscription's status
func (s *subscriptionService) UpdateSubscriptionStatus(id uint, status models.SubscriptionStatus) error {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return err
	}

	oldStatus := subscription.Status

	if err := s.subscriptionRepo.UpdateStatus(id, status); err != nil {
		return err
	}

	// Log status change event
	event := &models.SubscriptionEvent{
		SubscriptionID: subscription.ID,
		EventType:      "status_changed",
		PreviousStatus: oldStatus,
		CurrentStatus:  status,
		Notes:          fmt.Sprintf("Status manually changed from %s to %s", oldStatus, status),
	}

	if err := s.subscriptionRepo.LogEvent(event); err != nil {
		fmt.Printf("Failed to log subscription event: %v\n", err)
	}

	return nil
}

// ProcessPaymentWebhook handles payment webhooks from the payment service
func (s *subscriptionService) ProcessPaymentWebhook(payload models.PaymentWebhookPayload) error {
	// Handle different event types
	switch payload.Event {
	case "payment.completed", "payment.success":
		// When payment completes, activate the subscription
		_, err := s.ConfirmSubscription(models.PaymentConfirmationRequest{
			OrderID: payload.OrderID,
		})
		return err

	case "payment.failed", "payment.error":
		// Find subscription by payment order ID
		var subscriptions []models.TutorSubscription
		var count int64
		filters := map[string]interface{}{
			"payment_order_id": payload.OrderID,
		}

		subscriptions, count, err := s.subscriptionRepo.GetAll(1, 1, filters)
		if err != nil {
			return err
		}

		if count == 0 || len(subscriptions) == 0 {
			return errors.New("subscription not found for this payment")
		}

		subscription := subscriptions[0]

		// Update subscription status to failed
		oldStatus := subscription.Status
		subscription.Status = models.SubscriptionPastDue

		if err := s.subscriptionRepo.Update(&subscription); err != nil {
			return err
		}

		// Log payment failure event
		event := &models.SubscriptionEvent{
			SubscriptionID: subscription.ID,
			EventType:      "payment_failed",
			PreviousStatus: oldStatus,
			CurrentStatus:  models.SubscriptionPastDue,
			Notes:          "Payment failed, subscription marked as past due",
		}

		if err := s.subscriptionRepo.LogEvent(event); err != nil {
			fmt.Printf("Failed to log subscription event: %v\n", err)
		}

		return nil

	case "payment.cancelled", "payment.canceled":
		// Find subscription by payment order ID
		var subscriptions []models.TutorSubscription
		var count int64
		filters := map[string]interface{}{
			"payment_order_id": payload.OrderID,
		}

		subscriptions, count, err := s.subscriptionRepo.GetAll(1, 1, filters)
		if err != nil {
			return err
		}

		if count == 0 || len(subscriptions) == 0 {
			return errors.New("subscription not found for this payment")
		}

		subscription := subscriptions[0]

		// If the subscription is still incomplete, cancel it
		if subscription.Status == models.SubscriptionIncomplete {
			oldStatus := subscription.Status
			subscription.Status = models.SubscriptionCanceled

			if err := s.subscriptionRepo.Update(&subscription); err != nil {
				return err
			}

			// Log cancellation event
			event := &models.SubscriptionEvent{
				SubscriptionID: subscription.ID,
				EventType:      "payment_canceled",
				PreviousStatus: oldStatus,
				CurrentStatus:  models.SubscriptionCanceled,
				Notes:          "Payment canceled, subscription canceled",
			}

			if err := s.subscriptionRepo.LogEvent(event); err != nil {
				fmt.Printf("Failed to log subscription event: %v\n", err)
			}
		}

		return nil

	default:
		return nil // Ignore unsupported event types
	}
}

// GetExpiringSoonSubscriptions finds subscriptions that will expire in the given number of days
func (s *subscriptionService) GetExpiringSoonSubscriptions(days int) ([]models.SubscriptionResponse, error) {
	subscriptions, err := s.subscriptionRepo.GetExpiringSoon(days)
	if err != nil {
		return nil, err
	}

	responses := make([]models.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		var price float64
		if subscription.BillingCycle == models.BillingMonthly {
			price = subscription.Plan.PriceMonthly
		} else {
			price = subscription.Plan.PriceAnnually
		}

		responses[i] = models.SubscriptionResponse{
			ID:                 subscription.ID,
			TutorID:            subscription.TutorID,
			PlanName:           subscription.Plan.Name,
			Status:             subscription.Status,
			CurrentPeriodStart: subscription.CurrentPeriodStart,
			CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
			CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
			BillingCycle:       subscription.BillingCycle,
			Price:              price,
			Features:           subscription.Plan.FeaturesJSON,
			MaxCourses:         subscription.Plan.MaxCourses,
			CommissionRate:     subscription.Plan.CommissionRate,
		}
	}

	return responses, nil
}
