// internal/models/plan.go
package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type PlanType string

const (
	PlanBasic        PlanType = "basic"
	PlanPremium      PlanType = "premium"
	PlanProfessional PlanType = "professional"
)

type BillingCycle string

const (
	BillingMonthly  BillingCycle = "monthly"
	BillingAnnually BillingCycle = "annually"
)

// SubscriptionPlan represents a plan that tutors can subscribe to
type SubscriptionPlan struct {
	gorm.Model
	Name           string   `gorm:"size:100;not null" json:"name"`
	Description    string   `gorm:"type:text" json:"description"`
	PriceMonthly   float64  `gorm:"type:decimal(10,2);not null" json:"price_monthly"`
	PriceAnnually  float64  `gorm:"type:decimal(10,2);not null" json:"price_annually"`
	MaxCourses     int      `gorm:"not null" json:"max_courses"`
	CommissionRate float64  `gorm:"type:decimal(5,2);not null" json:"commission_rate"`
	Features       string   `gorm:"type:json" json:"-"`
	FeaturesJSON   []string `gorm:"-" json:"features"`
	IsActive       bool     `gorm:"default:true" json:"is_active"`
}

// TableName specifies the table name for the SubscriptionPlan model
func (SubscriptionPlan) TableName() string {
	return "SubscriptionPlan"
}

// BeforeSave hook converts FeaturesJSON to Features JSON string before saving
func (p *SubscriptionPlan) BeforeSave(tx *gorm.DB) error {
	if len(p.FeaturesJSON) > 0 {
		featuresBytes, err := json.Marshal(p.FeaturesJSON)
		if err != nil {
			return err
		}
		p.Features = string(featuresBytes)
	}
	return nil
}

// AfterFind hook converts Features JSON string to FeaturesJSON after retrieval
func (p *SubscriptionPlan) AfterFind(tx *gorm.DB) error {
	if p.Features != "" {
		var features []string
		if err := json.Unmarshal([]byte(p.Features), &features); err != nil {
			return err
		}
		p.FeaturesJSON = features
	}
	return nil
}

type PlanRequest struct {
	Name           string   `json:"name" validate:"required"`
	Description    string   `json:"description"`
	PriceMonthly   float64  `json:"price_monthly" validate:"required,gt=0"`
	PriceAnnually  float64  `json:"price_annually" validate:"required,gt=0"`
	MaxCourses     int      `json:"max_courses" validate:"required,gt=0"`
	CommissionRate float64  `json:"commission_rate" validate:"required,gte=0,lte=100"`
	Features       []string `json:"features"`
	IsActive       bool     `json:"is_active"`
}
