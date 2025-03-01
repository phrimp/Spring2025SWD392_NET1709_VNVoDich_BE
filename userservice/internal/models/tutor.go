package models

import (
	"time"

	"gorm.io/gorm"
)

type Tutor struct {
	ID             uint   `gorm:"primaryKey"`
	Bio            string `gorm:"type:text"`
	Qualifications string `gorm:"type:text"`
	TeachingStyle  string `gorm:"type:text"`
	IsAvailable    bool   `gorm:"default:false"`
	DemoVideoURL   string `gorm:"type:varchar(255);null"`
	Image          string `gorm:"type:varchar(255);null"`

	// Relationship with User model (Tutor is linked to a User)
	User User `gorm:"foreignKey:ID"`

	// Relationships with other models
	TutorSpecialties []TutorSpecialty `gorm:"foreignKey:TutorID"`
	Courses          []Course         `gorm:"foreignKey:TutorID"`
	TutorReviews     []TutorReview    `gorm:"foreignKey:TutorID"`
	Availability     *Availability    `gorm:"foreignKey:TutorID"`
}

type TutorSpecialty struct {
	gorm.Model
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	Subject         string `gorm:"type:varchar(100)"`
	Level           string `gorm:"type:varchar(50)"`
	Certification   string `gorm:"type:varchar(255)"`
	YearsExperience int

	TutorID uint
	Tutor   Tutor `gorm:"foreignKey:TutorID"`
}

type TutorReview struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Rating        int
	ReviewContent string `gorm:"type:text;null"`
	CreatedAt     *time.Time

	TutorID  uint
	Tutor    Tutor `gorm:"foreignKey:TutorID"`
	ParentID uint
	Parent   Parent `gorm:"foreignKey:ParentID"`
}
