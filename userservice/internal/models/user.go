package models

import (
	"errors"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleParent UserRole = "parent"
	RoleKid    UserRole = "kid"
	RoleTutor  UserRole = "tutor"
	RoleAdmin  UserRole = "admin"
)

type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusBanned    UserStatus = "banned"
)

type User struct {
	gorm.Model
	Username string   `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string   `gorm:"not null" json:"-" validate:"required,min=8"`
	Email    string   `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Role     UserRole `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	Phone    *string  `gorm:"type:varchar(20)" json:"phone" validate:"omitempty,e164"`
	FullName *string  `gorm:"type:varchar(255)" json:"full_name" validate:"omitempty,min=2,max=100"`

	IsVerified             bool       `gorm:"default:false" json:"is_verified"`
	Status                 UserStatus `gorm:"type:varchar(20);default:'inactive'" json:"status"`
	EmailVerificationToken *string    `json:"-"`

	LastLoginAt       *int64 `json:"last_login_at,omitempty"`
	LoginAttempts     uint   `gorm:"default:0" json:"-"`
	AccountLocked     bool   `gorm:"default:false" json:"account_locked"`
	PasswordChangedAt *int64 `json:"-"`

	PasswordResetToken   *string `json:"-"`
	PasswordResetExpires *int64  `json:"-"`
}

func (User) TableName() string {
	return "users"
}

// BeforeSave handles any necessary modifications before saving to database
func (u *User) BeforeSave(tx *gorm.DB) error {
	log.Printf("BeforeSave hook called for user: %s", u.Username)
	log.Printf("Password before hash: %s", u.Password)

	if tx.Statement.Changed("Password") || tx.Statement.ReflectValue.FieldByName("ID").IsZero() {
		log.Println("Password field has changed, hashing...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return err
		}
		u.Password = string(hashedPassword)
		now := time.Now().Unix()
		u.PasswordChangedAt = &now
		log.Printf("Password after hash: %s", u.Password)
	} else {
		log.Println("Password field has not changed")
	}
	return nil
}

func (u *User) Validate() error {
	// Username validation
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("username can only contain alphanumeric characters and underscores")
	}

	// Phone validation
	if u.Phone != nil {
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !phoneRegex.MatchString(*u.Phone) {
			return errors.New("invalid phone number format")
		}
	}

	// Full name validation
	if u.FullName != nil {
		nameRegex := regexp.MustCompile(`^[a-zA-Z .-]+$`)
		if !nameRegex.MatchString(*u.FullName) {
			return errors.New("full name can only contain letters, spaces, dots, and hyphens")
		}
	}

	// Role validation
	validRoles := map[UserRole]bool{
		RoleParent: true,
		RoleKid:    true,
		RoleTutor:  true,
		RoleAdmin:  true,
	}
	if !validRoles[u.Role] {
		return errors.New("invalid user role")
	}

	// Status validation
	validStatuses := map[UserStatus]bool{
		StatusActive:    true,
		StatusInactive:  true,
		StatusSuspended: true,
		StatusBanned:    true,
	}
	if !validStatuses[u.Status] {
		return errors.New("invalid user status")
	}

	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) IncrementLoginAttempts() {
	u.LoginAttempts++
	if u.LoginAttempts >= 5 {
		u.AccountLocked = true
	}
}

func (u *User) IsAccountValid() error {
	if u.AccountLocked {
		return errors.New("account is locked due to too many failed login attempts")
	}

	if u.Status != StatusActive {
		return errors.New("account is not active")
	}

	return nil
}
