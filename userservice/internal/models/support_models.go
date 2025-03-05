package models

import (
	"time"

)

type Course struct {
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

func (Course) TableName() string {
	return "Course"
}

type CourseReview struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Rating        int
	ReviewContent string `gorm:"type:text;null"`
	CreatedAt     *time.Time

	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`
	ParentID uint
	Parent   Parent `gorm:"foreignKey:ParentID"`
}

func (CourseReview) TableName() string {
	return "CourseReview"
}

type Availability struct {
	ID      uint  `gorm:"primaryKey;autoIncrement"`
	TutorID uint  `gorm:"uniqueIndex"`
	Tutor   Tutor `gorm:"foreignKey:TutorID"`

	TimeGap int // Minimum gap between bookings in minutes

	Days []DayAvailability `gorm:"foreignKey:AvailabilityID"`
}

func (Availability) TableName() string {
	return "Availability"
}

type CourseSubscription struct {
	ID                uint   `gorm:"primaryKey;autoIncrement"`
	Status            string `gorm:"type:varchar(20)"`
	SessionsRemaining int    `gorm:"null"`

	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`

	ChildrenID uint
	Children   Children `gorm:"foreignKey:ChildrenID"`
}

func (CourseSubscription) TableName() string {
	return "CourseSubscription"
}

type DayAvailability struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	AvailabilityID uint
	Day            string `gorm:"type:varchar(20)"`
	StartTime      time.Time
	EndTime        time.Time

	Availability Availability `gorm:"foreignKey:AvailabilityID"`
}

func (DayAvailability) TableName() string {
	return "DayAvailability"
}
