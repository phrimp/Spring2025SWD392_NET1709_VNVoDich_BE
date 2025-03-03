package models

type Parent struct {
	ID                   uint   `gorm:"primaryKey"`
	PreferredLanguage    string `gorm:"type:varchar(50)"`
	NotificationsEnabled bool   `gorm:"default:false"`

	// Relationship with User model (Parent is linked to a User)
	User User `gorm:"foreignKey:ID"`

	// Relationships with other models
	Children      []Children     `gorm:"foreignKey:ParentID"`
	TutorReviews  []TutorReview  `gorm:"foreignKey:ParentID"`
	CourseReviews []CourseReview `gorm:"foreignKey:ParentID"`
}

func (Parent) TableName() string {
	return "Parent"
}
