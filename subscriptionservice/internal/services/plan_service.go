// internal/services/plan_service.go
package services

import (
	"subscription/internal/models"
	"subscription/internal/repository"
)

type PlanService interface {
	CreatePlan(req models.PlanRequest) (*models.SubscriptionPlan, error)
	GetPlanByID(id uint) (*models.SubscriptionPlan, error)
	GetAllPlans(activeOnly bool) ([]models.SubscriptionPlan, error)
	UpdatePlan(id uint, req models.PlanRequest) (*models.SubscriptionPlan, error)
	DeletePlan(id uint) error
}

type planService struct {
	planRepo repository.PlanRepository
}

func NewPlanService(planRepo repository.PlanRepository) PlanService {
	return &planService{
		planRepo: planRepo,
	}
}

func (s *planService) CreatePlan(req models.PlanRequest) (*models.SubscriptionPlan, error) {
	plan := &models.SubscriptionPlan{
		Name:           req.Name,
		Description:    req.Description,
		PriceMonthly:   req.PriceMonthly,
		PriceAnnually:  req.PriceAnnually,
		MaxCourses:     req.MaxCourses,
		CommissionRate: req.CommissionRate,
		FeaturesJSON:   req.Features,
		IsActive:       req.IsActive,
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *planService) GetPlanByID(id uint) (*models.SubscriptionPlan, error) {
	return s.planRepo.GetByID(id)
}

func (s *planService) GetAllPlans(activeOnly bool) ([]models.SubscriptionPlan, error) {
	return s.planRepo.GetAll(activeOnly)
}

// UpdatePlan updates an existing subscription plan
func (s *planService) UpdatePlan(id uint, req models.PlanRequest) (*models.SubscriptionPlan, error) {
	// Check if the plan exists
	existingPlan, err := s.planRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update the plan
	existingPlan.Name = req.Name
	existingPlan.Description = req.Description
	existingPlan.PriceMonthly = req.PriceMonthly
	existingPlan.PriceAnnually = req.PriceAnnually
	existingPlan.MaxCourses = req.MaxCourses
	existingPlan.CommissionRate = req.CommissionRate
	existingPlan.FeaturesJSON = req.Features
	existingPlan.IsActive = req.IsActive

	if err := s.planRepo.Update(existingPlan); err != nil {
		return nil, err
	}

	return existingPlan, nil
}

// DeletePlan soft-deletes a subscription plan
func (s *planService) DeletePlan(id uint) error {
	// Check if the plan exists
	if _, err := s.planRepo.GetByID(id); err != nil {
		return err
	}

	return s.planRepo.Delete(id)
}
