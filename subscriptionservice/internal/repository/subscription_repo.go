// internal/repository/subscription_repository.go
package repository

import (
	"errors"
	"subscription/internal/models"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.TutorSubscription) error
	GetByID(id uint) (*models.TutorSubscription, error)
	GetByTutorID(tutorID uint) (*models.TutorSubscription, error)
	GetAll(page, pageSize int, filters map[string]interface{}) ([]models.TutorSubscription, int64, error)
	Update(subscription *models.TutorSubscription) error
	UpdateStatus(id uint, status models.SubscriptionStatus) error
	SetCancelAtPeriodEnd(id uint, cancel bool) error
	LogEvent(event *models.SubscriptionEvent) error
	GetExpiringSoon(days int) ([]models.TutorSubscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

// Create adds a new tutor subscription
func (r *subscriptionRepository) Create(subscription *models.TutorSubscription) error {
	return r.db.Create(subscription).Error
}

// GetByID retrieves a subscription by its ID
func (r *subscriptionRepository) GetByID(id uint) (*models.TutorSubscription, error) {
	var subscription models.TutorSubscription
	result := r.db.Preload("Plan").First(&subscription, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription not found")
		}
		return nil, result.Error
	}
	return &subscription, nil
}

// GetByTutorID retrieves a subscription by tutor ID
func (r *subscriptionRepository) GetByTutorID(tutorID uint) (*models.TutorSubscription, error) {
	var subscription models.TutorSubscription
	result := r.db.Preload("Plan").Where("tutor_id = ? AND status != ?", tutorID, models.SubscriptionCanceled).First(&subscription)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription not found")
		}
		return nil, result.Error
	}
	return &subscription, nil
}

// GetAll retrieves all subscriptions with pagination and filters
func (r *subscriptionRepository) GetAll(page, pageSize int, filters map[string]interface{}) ([]models.TutorSubscription, int64, error) {
	var subscriptions []models.TutorSubscription
	var count int64

	query := r.db.Model(&models.TutorSubscription{})

	// Apply filters
	if filters != nil {
		if status, ok := filters["status"].(models.SubscriptionStatus); ok {
			query = query.Where("status = ?", status)
		}

		if tutorID, ok := filters["tutor_id"].(uint); ok {
			query = query.Where("tutor_id = ?", tutorID)
		}

		if planID, ok := filters["plan_id"].(uint); ok {
			query = query.Where("plan_id = ?", planID)
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
	result := query.Preload("Plan").
		Limit(pageSize).
		Offset(offset).
		Find(&subscriptions)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return subscriptions, count, nil
}

// Update updates an existing subscription
func (r *subscriptionRepository) Update(subscription *models.TutorSubscription) error {
	return r.db.Save(subscription).Error
}

// UpdateStatus updates the status of a subscription
func (r *subscriptionRepository) UpdateStatus(id uint, status models.SubscriptionStatus) error {
	result := r.db.Model(&models.TutorSubscription{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// SetCancelAtPeriodEnd marks a subscription to be canceled at the end of the period
func (r *subscriptionRepository) SetCancelAtPeriodEnd(id uint, cancel bool) error {
	result := r.db.Model(&models.TutorSubscription{}).
		Where("id = ?", id).
		Update("cancel_at_period_end", cancel)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// LogEvent records a subscription lifecycle event
func (r *subscriptionRepository) LogEvent(event *models.SubscriptionEvent) error {
	return r.db.Create(event).Error
}

// GetExpiringSoon finds subscriptions that will expire in the given number of days
func (r *subscriptionRepository) GetExpiringSoon(days int) ([]models.TutorSubscription, error) {
	var subscriptions []models.TutorSubscription

	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	result := r.db.Preload("Plan").
		Where("status = ? AND current_period_end BETWEEN ? AND ?",
			models.SubscriptionActive, now, futureDate).
		Find(&subscriptions)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
