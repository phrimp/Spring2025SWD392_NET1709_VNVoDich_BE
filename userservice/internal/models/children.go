package models

type Children struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Age           int
	GradeLevel    string `gorm:"type:varchar(50)"`
	LearningGoals string `gorm:"type:text"`
	FullName      string `gorm:"type:varchar(255)"`
	Password      string `gorm:"type:varchar(255)"`

	ParentID uint
	Parent   Parent `gorm:"foreignKey:ParentID"`

	CourseSubscriptions []CourseSubscription `gorm:"foreignKey:ChildrenID"`
}
