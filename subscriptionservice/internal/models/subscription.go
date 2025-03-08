package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	SubscriptionActive     SubscriptionStatus = "active"
	SubscriptionCanceled   SubscriptionStatus = "canceled"
	SubscriptionPastDue    SubscriptionStatus = "past_due"
	SubscriptionTrialing   SubscriptionStatus = "trialing"
	SubscriptionIncomplete SubscriptionStatus = "incomplete"
)

type TutorSubscription struct {
	gorm.Model
	TutorID            uint               `gorm:"not null" json:"tutor_id"`
	PlanID             uint               `gorm:"not null" json:"plan_id"`
	Status             SubscriptionStatus `gorm:"type:enum('active','canceled','past_due','trialing','incomplete');not null" json:"status"`
	CurrentPeriodStart time.Time          `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `gorm:"not null" json:"current_period_end"`
	CancelAtPeriodEnd  bool               `gorm:"default:false" json:"cancel_at_period_end"`
	BillingCycle       BillingCycle       `gorm:"type:enum('monthly','annually');not null" json:"billing_cycle"`
	PaymentOrderID     string             `gorm:"size:255" json:"payment_order_id"`

	// Relations
	Plan SubscriptionPlan `gorm:"foreignKey:PlanID" json:"plan"`
}

func (TutorSubscription) TableName() string {
	return "TutorSubscriptions"
}

// SubscriptionRequest is used for creating a new subscription
type SubscriptionRequest struct {
	TutorID      uint         `json:"tutor_id" validate:"required"`
	PlanID       uint         `json:"plan_id" validate:"required"`
	BillingCycle BillingCycle `json:"billing_cycle" validate:"required,oneof=monthly annually"`
}

// SubscriptionResponse is the response returned after subscription operations
type SubscriptionResponse struct {
	ID                 uint               `json:"id"`
	TutorID            uint               `json:"tutor_id"`
	PlanName           string             `json:"plan_name"`
	Status             SubscriptionStatus `json:"status"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	CancelAtPeriodEnd  bool               `json:"cancel_at_period_end"`
	BillingCycle       BillingCycle       `json:"billing_cycle"`
	Price              float64            `json:"price"`
	Features           []string           `json:"features"`
	MaxCourses         int                `json:"max_courses"`
	CommissionRate     float64            `json:"commission_rate"`
	PaymentOrderID     string             `json:"payment_order_id,omitempty"`
	PaymentURL         string             `json:"payment_url,omitempty"` // URL to redirect for payment
}

// SubscriptionStatusUpdateRequest is used to update the status of a subscription
type SubscriptionStatusUpdateRequest struct {
	Status SubscriptionStatus `json:"status" validate:"required,oneof=active canceled past_due trialing incomplete"`
}

// ChangePlanRequest is used when a tutor wants to change their subscription plan
type ChangePlanRequest struct {
	NewPlanID    uint         `json:"new_plan_id" validate:"required"`
	BillingCycle BillingCycle `json:"billing_cycle" validate:"required,oneof=monthly annually"`
}

// PaymentConfirmationRequest is used to confirm payment completion
type PaymentConfirmationRequest struct {
	OrderID   string `json:"order_id" validate:"required"`
	PaymentID string `json:"payment_id"`
	PayerID   string `json:"payer_id"`
}

// SubscriptionEvent represents an event in the subscription lifecycle
type SubscriptionEvent struct {
	gorm.Model
	SubscriptionID uint               `gorm:"not null" json:"subscription_id"`
	EventType      string             `gorm:"size:50;not null" json:"event_type"` // created, renewed, canceled, plan_changed, payment_failed, ...
	PreviousStatus SubscriptionStatus `json:"previous_status"`
	CurrentStatus  SubscriptionStatus `json:"current_status"`
	Notes          string             `gorm:"type:text" json:"notes"`
}

// TableName specifies the table name for the SubscriptionEvent model
func (SubscriptionEvent) TableName() string {
	return "SubscriptionEvents"
}

// PaymentWebhookPayload represents the structure of payment webhooks
type PaymentWebhookPayload struct {
	Event   string `json:"event"`
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}
