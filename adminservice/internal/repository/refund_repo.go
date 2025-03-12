package repository

import (
	"adminservice/internal/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RefundRepository interface {
	Create(refund *models.RefundRequest) error
	GetByID(id uint) (*models.RefundRequest, error)
	GetAll(page, pageSize int, filters map[string]interface{}) ([]models.RefundRequest, int64, error)
	Update(refund *models.RefundRequest) error
	UpdateStatus(id uint, status models.RefundStatus, processedBy uint, adminNote string) error
	MarkNotificationSent(id uint) error
}

type refundRepository struct {
	db *gorm.DB
}

// NewRefundRepository creates a new instance of RefundRepository
func NewRefundRepository(db *gorm.DB) RefundRepository {
	return &refundRepository{
		db: db,
	}
}

// Create adds a new refund request
func (r *refundRepository) Create(refund *models.RefundRequest) error {
	return r.db.Create(refund).Error
}

// GetByID retrieves a refund request by its ID
func (r *refundRepository) GetByID(id uint) (*models.RefundRequest, error) {
	var refund models.RefundRequest
	result := r.db.First(&refund, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("refund request not found")
		}
		return nil, result.Error
	}
	return &refund, nil
}

// GetAll retrieves all refund requests with pagination and filters
func (r *refundRepository) GetAll(page, pageSize int, filters map[string]interface{}) ([]models.RefundRequest, int64, error) {
	var refunds []models.RefundRequest
	var count int64

	query := r.db.Model(&models.RefundRequest{})

	// Apply filters
	if filters != nil {
		if status, ok := filters["status"].(models.RefundStatus); ok && status != "" {
			query = query.Where("status = ?", status)
		}

		if userID, ok := filters["user_id"].(uint); ok && userID > 0 {
			query = query.Where("user_id = ?", userID)
		}

		if orderID, ok := filters["order_id"].(string); ok && orderID != "" {
			query = query.Where("order_id = ?", orderID)
		}

		if username, ok := filters["username"].(string); ok && username != "" {
			query = query.Where("username LIKE ?", "%"+username+"%")
		}

		if email, ok := filters["email"].(string); ok && email != "" {
			query = query.Where("email LIKE ?", "%"+email+"%")
		}

		if fromDate, ok := filters["from_date"].(time.Time); ok {
			query = query.Where("created_at >= ?", fromDate)
		}

		if toDate, ok := filters["to_date"].(time.Time); ok {
			query = query.Where("created_at <= ?", toDate)
		}
	}

	// Get total count for pagination
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	result := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&refunds)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return refunds, count, nil
}

// Update updates an existing refund request
func (r *refundRepository) Update(refund *models.RefundRequest) error {
	return r.db.Save(refund).Error
}

// UpdateStatus updates the status of a refund request
func (r *refundRepository) UpdateStatus(id uint, status models.RefundStatus, processedBy uint, adminNote string) error {
	now := time.Now()
	result := r.db.Model(&models.RefundRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"processed_by": processedBy,
			"processed_at": now,
			"admin_note":   adminNote,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("refund request with ID %d not found", id)
	}

	return nil
}

// MarkNotificationSent marks a refund request notification as sent
func (r *refundRepository) MarkNotificationSent(id uint) error {
	result := r.db.Model(&models.RefundRequest{}).
		Where("id = ?", id).
		Update("notification_sent", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("refund request with ID %d not found", id)
	}

	return nil
}
