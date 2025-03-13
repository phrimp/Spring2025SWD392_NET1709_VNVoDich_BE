package models

import (
	"time"

	"gorm.io/gorm"
)

// RefundStatus represents the status of a refund request
type RefundStatus string

const (
	RefundStatusPending  RefundStatus = "PENDING"
	RefundStatusApproved RefundStatus = "APPROVED"
	RefundStatusRejected RefundStatus = "REJECTED"
)

// RefundRequest represents a user's request for a refund
type RefundRequest struct {
	gorm.Model
	UserID           uint         `json:"user_id" gorm:"not null"`
	Username         string       `json:"username" gorm:"not null"`
	Email            string       `json:"email" gorm:"not null"`
	OrderID          string       `json:"order_id" gorm:"not null"`
	Amount           float64      `json:"amount" gorm:"not null"`
	CardNumber       string       `json:"card_number" gorm:"not null"` // Last 4 digits only
	Reason           string       `json:"reason" gorm:"type:text"`
	Status           RefundStatus `json:"status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	ProcessedBy      *uint        `json:"processed_by"`
	ProcessedAt      *time.Time   `json:"processed_at"`
	AdminNote        string       `json:"admin_note" gorm:"type:text"`
	NotificationSent bool         `json:"notification_sent" gorm:"default:false"`
}

// TableName specifies the table name for RefundRequest
func (RefundRequest) TableName() string {
	return "RefundRequests"
}

// RefundRequestInput represents the input for creating a refund request
type RefundRequestInput struct {
	OrderID    string  `json:"order_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	CardNumber string  `json:"card_number" validate:"required"`
	Reason     string  `json:"reason"`
}

// RefundProcessInput represents the input for processing a refund request
type RefundProcessInput struct {
	Status    RefundStatus `json:"status" validate:"required,oneof=approved rejected"`
	AdminNote string       `json:"admin_note"`
}

// PaginatedRefundResponse represents a paginated list of refund requests
type PaginatedRefundResponse struct {
	Data       []RefundRequest `json:"data"`
	Pagination Pagination      `json:"pagination"`
}
