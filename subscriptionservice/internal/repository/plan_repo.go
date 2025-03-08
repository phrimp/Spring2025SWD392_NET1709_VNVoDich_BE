package repository

import (
	"errors"
	"subscription/internal/models"
	"time"

	"gorm.io/gorm"
)

type PlanRepository interface {
	Create(plan *models.SubscriptionPlan) error
	GetByID(id uint) (*models.SubscriptionPlan, error)
	GetAll(activeOnly bool) ([]models.SubscriptionPlan, error)
	Update(plan *models.SubscriptionPlan) error
	Delete(id uint) error
}

type planRepository struct {
	db *gorm.DB
}

// NewPlanRepository creates a new instance of PlanRepository
func NewPlanRepository(db *gorm.DB) PlanRepository {
	return &planRepository{
		db: db,
	}
}

// Create adds a new subscription plan
func (r *planRepository) Create(plan *models.SubscriptionPlan) error {
	return r.db.Create(plan).Error
}

// GetByID retrieves a subscription plan by its ID
func (r *planRepository) GetByID(id uint) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	result := r.db.First(&plan, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan not found")
		}
		return nil, result.Error
	}
	return &plan, nil
}

// GetAll retrieves all subscription plans
func (r *planRepository) GetAll(activeOnly bool) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	query := r.db

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	result := query.Find(&plans)
	if result.Error != nil {
		return nil, result.Error
	}
	return plans, nil
}

// Update updates an existing subscription plan
func (r *planRepository) Update(plan *models.SubscriptionPlan) error {
	return r.db.Save(plan).Error
}

// Delete soft-deletes a subscription plan by setting is_active to false
func (r *planRepository) Delete(id uint) error {
	return r.db.Model(&models.SubscriptionPlan{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  false,
			"deleted_at": time.Now(),
		}).
		Error
}
