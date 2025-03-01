package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Title        string `gorm:"type:varchar(255)"`
	Description  string `gorm:"type:text;null"`
	Price        float64
	TotalLessons int
	Grade        int
	Subject      string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(20)"`
	Image        string `gorm:"type:varchar(255);null"`

	TutorID uint
	Tutor   Tutor `gorm:"foreignKey:TutorID"`
}

type CourseReview struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Rating        int
	ReviewContent string `gorm:"type:text;null"`
	CreatedAt     *time.Time

	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`
	ParentID uint
	Parent   Parent `gorm:"foreignKey:ParentID"`
}

type Availability struct {
	gorm.Model
	ID      uint  `gorm:"primaryKey;autoIncrement"`
	TutorID uint  `gorm:"uniqueIndex"`
	Tutor   Tutor `gorm:"foreignKey:TutorID"`

	TimeGap int // Minimum gap between bookings in minutes

	Days []DayAvailability `gorm:"foreignKey:AvailabilityID"`
}

type CourseSubscription struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey;autoIncrement"`
	Status            string `gorm:"type:varchar(20)"`
	SessionsRemaining int    `gorm:"null"`

	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`

	ChildrenID uint
	Children   Children `gorm:"foreignKey:ChildrenID"`
}

type DayAvailability struct {
	gorm.Model
	ID             uint `gorm:"primaryKey;autoIncrement"`
	AvailabilityID uint
	Day            string `gorm:"type:varchar(20)"`
	StartTime      time.Time
	EndTime        time.Time

	Availability Availability `gorm:"foreignKey:AvailabilityID"`
}
