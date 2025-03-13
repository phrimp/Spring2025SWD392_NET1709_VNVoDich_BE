package services

import (
	"adminservice/internal/models"
	"adminservice/internal/repository"
	"adminservice/utils"
	"fmt"
	"math"
	"net/url"

	"github.com/valyala/fasthttp"
)

type RefundService interface {
	CreateRefundRequest(userID uint, username, email string, input models.RefundRequestInput) (*models.RefundRequest, error)
	GetRefundRequestByID(id uint) (*models.RefundRequest, error)
	GetAllRefundRequests(page, pageSize int, filters map[string]interface{}) (*models.PaginatedRefundResponse, error)
	ProcessRefundRequest(id, adminID uint, input models.RefundProcessInput) error
	GetRefundStatistics() (map[string]interface{}, error)
}

type refundService struct {
	refundRepo       repository.RefundRepository
	googleServiceURL string
	apiKey           string
}

// NewRefundService creates a new instance of RefundService
func NewRefundService(refundRepo repository.RefundRepository, googleServiceURL, apiKey string) RefundService {
	return &refundService{
		refundRepo:       refundRepo,
		googleServiceURL: googleServiceURL,
		apiKey:           apiKey,
	}
}

// CreateRefundRequest creates a new refund request
func (s *refundService) CreateRefundRequest(userID uint, username, email string, input models.RefundRequestInput) (*models.RefundRequest, error) {
	// Mask the card number to only keep the last 4 digits
	var maskedCardNumber string
	if len(input.CardNumber) > 4 {
		maskedCardNumber = "XXXX-XXXX-XXXX-" + input.CardNumber[len(input.CardNumber)-4:]
	} else {
		maskedCardNumber = input.CardNumber
	}

	refund := &models.RefundRequest{
		UserID:     userID,
		Username:   username,
		Email:      email,
		OrderID:    input.OrderID,
		Amount:     input.Amount,
		CardNumber: maskedCardNumber,
		Reason:     input.Reason,
		Status:     models.RefundStatusPending,
	}

	if err := s.refundRepo.Create(refund); err != nil {
		return nil, fmt.Errorf("failed to create refund request: %w", err)
	}

	return refund, nil
}

// GetRefundRequestByID retrieves a refund request by its ID
func (s *refundService) GetRefundRequestByID(id uint) (*models.RefundRequest, error) {
	return s.refundRepo.GetByID(id)
}

// GetAllRefundRequests retrieves all refund requests with pagination and filters
func (s *refundService) GetAllRefundRequests(page, pageSize int, filters map[string]interface{}) (*models.PaginatedRefundResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	refunds, total, err := s.refundRepo.GetAll(page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.PaginatedRefundResponse{
		Data: refunds,
		Pagination: models.Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}, nil
}

// ProcessRefundRequest processes a refund request (approve or reject)
func (s *refundService) ProcessRefundRequest(id, adminID uint, input models.RefundProcessInput) error {
	// Get the refund request to make sure it exists and check its current status
	refund, err := s.refundRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Validate status transition
	if refund.Status != models.RefundStatusPending {
		return fmt.Errorf("cannot process refund request with status %s", refund.Status)
	}

	// Update the refund request status
	if err := s.refundRepo.UpdateStatus(id, input.Status, adminID, input.AdminNote); err != nil {
		return err
	}

	// If approved, send an email notification
	if input.Status == models.RefundStatusApproved {
		if err := s.sendRefundNotification(refund); err != nil {
			return fmt.Errorf("refund status updated but failed to send notification: %w", err)
		}
	}

	return nil
}

// sendRefundNotification sends an email notification about the refund
func (s *refundService) sendRefundNotification(refund *models.RefundRequest) error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Prepare email content
	emailTitle := "Your Refund Request Has Been Approved"
	emailBody := fmt.Sprintf(
		"Dear %s,\n\nYour refund request for order %s in the amount of %.2f has been approved.\n\nThe refund will be processed to the card ending in %s within 5-10 business days.\n\nIf you have any questions, please contact our customer support.\n\nThank you for your patience.\n\nBest regards,\nOnline Tutoring Platform Team",
		refund.Username,
		refund.OrderID,
		refund.Amount,
		refund.CardNumber[len(refund.CardNumber)-4:],
	)

	// Build the request to the email service
	url := fmt.Sprintf("%s/api/email/send?to=%s&title=%s&body=%s",
		s.googleServiceURL,
		url.QueryEscape(refund.Email),
		url.QueryEscape(emailTitle),
		url.QueryEscape(emailBody),
	)

	fmt.Println(url)
	utils.BuildRequest(req, "POST", nil, s.apiKey, url)

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("email service unavailable: %w", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("failed to send email: %s", resp.Body())
	}

	// Mark notification as sent
	if err := s.refundRepo.MarkNotificationSent(refund.ID); err != nil {
		return fmt.Errorf("failed to mark notification as sent: %w", err)
	}

	return nil
}

// GetRefundStatistics gets statistics about refund requests
func (s *refundService) GetRefundStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get total counts by status
	var pending, approved, rejected int64

	filters := map[string]interface{}{
		"status": models.RefundStatusPending,
	}
	_, pending, err := s.refundRepo.GetAll(1, 1, filters)
	if err != nil {
		return nil, err
	}

	filters["status"] = models.RefundStatusApproved
	_, approved, err = s.refundRepo.GetAll(1, 1, filters)
	if err != nil {
		return nil, err
	}

	filters["status"] = models.RefundStatusRejected
	_, rejected, err = s.refundRepo.GetAll(1, 1, filters)
	if err != nil {
		return nil, err
	}

	// Get total amount refunded
	var totalRefunded struct {
		TotalAmount float64
	}

	db := repository.DB
	err = db.Table("RefundRequests").
		Select("SUM(amount) as total_amount").
		Where("status = ?", models.RefundStatusApproved).
		Scan(&totalRefunded).Error
	if err != nil {
		return nil, err
	}

	// Get recent refund requests (last 5)
	var recentRefunds []models.RefundRequest
	err = db.Order("created_at DESC").Limit(5).Find(&recentRefunds).Error
	if err != nil {
		return nil, err
	}

	// Prepare statistics
	stats["total_pending"] = pending
	stats["total_approved"] = approved
	stats["total_rejected"] = rejected
	stats["total_refunded_amount"] = totalRefunded.TotalAmount
	stats["recent_refunds"] = recentRefunds

	return stats, nil
}
