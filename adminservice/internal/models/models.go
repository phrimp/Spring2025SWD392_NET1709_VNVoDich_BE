package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleParent UserRole = "Parent"
	RoleChild  UserRole = "Children"
	RoleTutor  UserRole = "Tutor"
	RoleAdmin  UserRole = "Admin"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	StatusActive    UserStatus = "Active"
	StatusSuspended UserStatus = "Suspended"
	StatusBanned    UserStatus = "Banned"
)

// User model represents the user table in the database
type User struct {
	gorm.Model
	Username          string     `json:"username" gorm:"uniqueIndex;not null"`
	Password          string     `json:"-" gorm:"not null"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null"`
	Role              UserRole   `json:"role" gorm:"type:ENUM('Parent', 'Children', 'Tutor', 'Admin');default:'Parent'"`
	Phone             *string    `json:"phone" gorm:"type:varchar(20)"`
	FullName          *string    `json:"full_name" gorm:"type:varchar(255)"`
	Picture           string     `json:"picture" gorm:"type:varchar(255)"`
	IsVerified        bool       `json:"is_verified" gorm:"default:false"`
	Status            UserStatus `json:"status" gorm:"type:ENUM('Active', 'Suspended', 'Banned');default:'Active'"`
	LastLoginAt       *int64     `json:"last_login_at,omitempty"`
	AccountLocked     bool       `json:"account_locked" gorm:"default:false"`
	PasswordChangedAt *int64     `json:"-"`

	// Relations
	Tutor *Tutor `json:"tutor,omitempty" gorm:"foreignKey:ID"`
}

// Tutor model represents a tutor profile
type Tutor struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	Bio            string  `json:"bio"`
	Qualifications string  `json:"qualifications"`
	TeachingStyle  string  `json:"teaching_style"`
	IsAvailable    bool    `json:"is_available" gorm:"default:false"`
	DemoVideoURL   *string `json:"demo_video_url"`
	Image          *string `json:"image"`

	// Relations
	ProfileID      uint             `json:"-" gorm:"column:id"`
	Profile        User             `json:"profile" gorm:"foreignKey:ProfileID;references:ID"`
	TutorSpecialty []TutorSpecialty `json:"specialties,omitempty" gorm:"foreignKey:TutorID"`
	Courses        []Course         `json:"courses,omitempty" gorm:"foreignKey:TutorID"`
	TutorReviews   []TutorReview    `json:"reviews,omitempty" gorm:"foreignKey:TutorID"`
	Availability   *Availability    `json:"availability,omitempty" gorm:"foreignKey:TutorID"`
}

// Parent model represents a parent profile
type Parent struct {
	ID                  uint   `json:"id" gorm:"primaryKey"`
	PreferredLanguage   string `json:"preferred_language"`
	NotificationsEnable bool   `json:"notifications_enable" gorm:"default:false"`

	// Relations
	Children      []Children     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	TutorReviews  []TutorReview  `json:"tutor_reviews,omitempty" gorm:"foreignKey:ParentID"`
	CourseReviews []CourseReview `json:"course_reviews,omitempty" gorm:"foreignKey:ParentID"`
}

// Children model represents a child profile
type Children struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Age           int    `json:"age"`
	GradeLevel    string `json:"grade_level"`
	LearningGoals string `json:"learning_goals"`
	FullName      string `json:"full_name"`
	Password      string `json:"-"`

	// Relations
	ParentID uint   `json:"parent_id"`
	Parent   Parent `json:"parent,omitempty" gorm:"foreignKey:ParentID"`

	CourseSubscriptions []CourseSubscription `json:"course_subscriptions,omitempty" gorm:"foreignKey:ChildrenID"`
}

// TutorSpecialty model represents a tutor's subject specialization
type TutorSpecialty struct {
	ID              uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Subject         string `json:"subject"`
	Level           string `json:"level"`
	Certification   string `json:"certification"`
	YearsExperience int    `json:"years_experience"`

	// Relations
	TutorID uint  `json:"tutor_id"`
	Tutor   Tutor `json:"-" gorm:"foreignKey:TutorID"`
}

// CourseStatus represents the status of a course
type CourseStatus string

const (
	CourseStatusDraft     CourseStatus = "Draft"
	CourseStatusPublished CourseStatus = "Published"
)

// Course model represents an educational course
type Course struct {
	ID           uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string       `json:"title"`
	Description  *string      `json:"description"`
	Price        float64      `json:"price"`
	TotalLessons int          `json:"total_lessons"`
	Grade        int          `json:"grade"`
	Subject      string       `json:"subject"`
	Status       CourseStatus `json:"status" gorm:"type:ENUM('Draft', 'Published');default:'Draft'"`
	Image        *string      `json:"image"`

	// Relations
	TutorID uint  `json:"tutor_id"`
	Tutor   Tutor `json:"tutor,omitempty" gorm:"foreignKey:TutorID"`

	Lessons             []Lesson             `json:"lessons,omitempty" gorm:"foreignKey:CourseID"`
	CourseSubscriptions []CourseSubscription `json:"subscriptions,omitempty" gorm:"foreignKey:CourseID"`
	CourseReviews       []CourseReview       `json:"reviews,omitempty" gorm:"foreignKey:CourseID"`
}

// Lesson model represents a lesson in a course
type Lesson struct {
	ID                 uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title              string  `json:"title"`
	Description        *string `json:"description"`
	LearningObjectives *string `json:"learning_objectives"`
	MaterialsNeeded    *string `json:"materials_needed"`

	// Relations
	CourseID uint   `json:"course_id"`
	Course   Course `json:"-" gorm:"foreignKey:CourseID"`
}

// Day type for days of the week
type Day string

const (
	Monday    Day = "MONDAY"
	Tuesday   Day = "TUESDAY"
	Wednesday Day = "WEDNESDAY"
	Thursday  Day = "THURSDAY"
	Friday    Day = "FRIDAY"
	Saturday  Day = "SATURDAY"
	Sunday    Day = "SUNDAY"
)

// Availability model represents a tutor's availability
type Availability struct {
	ID        uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	TutorID   uint              `json:"tutor_id" gorm:"uniqueIndex"`
	TimeGap   int               `json:"time_gap"` // Minimum gap between bookings in minutes
	Days      []DayAvailability `json:"days,omitempty" gorm:"foreignKey:AvailabilityID"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// DayAvailability model represents availability for a specific day
type DayAvailability struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AvailabilityID uint      `json:"availability_id"`
	Day            Day       `json:"day" gorm:"type:ENUM('MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY')"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
}

// CourseSubscription model represents a subscription to a course
type CourseSubscription struct {
	ID                uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Status            string `json:"status"`
	SessionsRemaining *int   `json:"sessions_remaining"`

	// Relations
	CourseID   uint     `json:"course_id"`
	Course     Course   `json:"course,omitempty" gorm:"foreignKey:CourseID"`
	ChildrenID uint     `json:"children_id"`
	Children   Children `json:"children,omitempty" gorm:"foreignKey:ChildrenID"`

	CourseSubscriptionSchedules []CourseSubscriptionSchedule `json:"schedules,omitempty" gorm:"foreignKey:SubscriptionID"`
	TeachingSessions            []TeachingSession            `json:"sessions,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// TeachingSession model represents a teaching session
type TeachingSession struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	GoogleMeetID     *string   `json:"google_meet_id"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Status           string    `json:"status"`
	TopicsCovered    *string   `json:"topics_covered"`
	HomeworkAssigned *string   `json:"homework_assigned"`

	// Relations
	SubscriptionID  uint               `json:"subscription_id"`
	Subscription    CourseSubscription `json:"-" gorm:"foreignKey:SubscriptionID"`
	SessionFeedback []SessionFeedback  `json:"feedback,omitempty" gorm:"foreignKey:SessionID"`
}

// SessionFeedback model represents feedback for a teaching session
type SessionFeedback struct {
	ID              uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating          int     `json:"rating"`
	Comment         *string `json:"comment"`
	TeachingQuality *string `json:"teaching_quality"`

	// Relations
	SessionID uint            `json:"session_id"`
	Session   TeachingSession `json:"-" gorm:"foreignKey:SessionID"`
}

// CourseSubscriptionSchedule model represents a scheduled time for a course
type CourseSubscriptionSchedule struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SubscriptionID uint      `json:"subscription_id"`
	Day            Day       `json:"day" gorm:"type:ENUM('MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY')"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`

	// Relations
	Subscription CourseSubscription `json:"-" gorm:"foreignKey:SubscriptionID"`
}

// TutorReview model represents a review for a tutor
type TutorReview struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating        int       `json:"rating"`
	ReviewContent *string   `json:"review_content"`
	CreateAt      time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	TutorID  uint   `json:"tutor_id"`
	Tutor    Tutor  `json:"-" gorm:"foreignKey:TutorID"`
	ParentID uint   `json:"parent_id"`
	Parent   Parent `json:"-" gorm:"foreignKey:ParentID"`
}

// CourseReview model represents a review for a course
type CourseReview struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating        int       `json:"rating"`
	ReviewContent *string   `json:"review_content"`
	CreateAt      time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	CourseID uint   `json:"course_id"`
	Course   Course `json:"-" gorm:"foreignKey:CourseID"`
	ParentID uint   `json:"parent_id"`
	Parent   Parent `json:"-" gorm:"foreignKey:ParentID"`
}

// DashboardStats represents summarized statistics for the admin dashboard
type DashboardStats struct {
	TotalUsers           int     `json:"total_users"`
	TotalTutors          int     `json:"total_tutors"`
	TotalParents         int     `json:"total_parents"`
	TotalChildren        int     `json:"total_children"`
	TotalCourses         int     `json:"total_courses"`
	TotalSubscriptions   int     `json:"total_subscriptions"`
	TotalSessions        int     `json:"total_sessions"`
	AverageRating        float64 `json:"average_rating"`
	ActiveSubscriptions  int     `json:"active_subscriptions"`
	PublishedCourses     int     `json:"published_courses"`
	AvailableTutors      int     `json:"available_tutors"`
	UpcomingSessions     int     `json:"upcoming_sessions"`
	TotalRevenue         float64 `json:"total_revenue"`
	MostPopularSubject   string  `json:"most_popular_subject"`
	MostActiveGradeLevel string  `json:"most_active_grade_level"`
}

// RecentActivity represents recent system activity for the dashboard
type RecentActivity struct {
	ID           uint      `json:"id"`
	ActivityType string    `json:"activity_type"` // e.g., "New Subscription", "Course Review", etc.
	Description  string    `json:"description"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	EntityID     uint      `json:"entity_id"`   // ID of the related entity
	EntityType   string    `json:"entity_type"` // Type of the related entity
	Timestamp    time.Time `json:"timestamp"`
}

// Pagination struct for handling paginated responses
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse is a generic struct for paginated responses
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
