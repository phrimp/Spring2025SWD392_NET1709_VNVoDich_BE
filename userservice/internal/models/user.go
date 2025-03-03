package models

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleParent UserRole = "Parent"
	RoleKid    UserRole = "Kid"
	RoleTutor  UserRole = "Tutor"
	RoleAdmin  UserRole = "Admin"
)

type UserStatus string

const (
	StatusActive  UserStatus = "Active"
	StatusDeleted UserStatus = "Deleted"
	StatusBanned  UserStatus = "Banned"
)

type User struct {
	gorm.Model
	Username string   `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string   `gorm:"not null" json:"-" validate:"required,min=8"`
	Email    string   `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Role     UserRole `gorm:"type:varchar(50);not null;default:'user'" json:"role"`
	Phone    *string  `gorm:"type:varchar(20)" json:"phone" validate:"omitempty,e164"`
	FullName *string  `gorm:"type:varchar(255)" json:"full_name" validate:"omitempty,min=2,max=100"`
	Picture  string   `gorm:"type:varchar(255)" json:"picture"`

	IsVerified bool       `gorm:"default:false" json:"is_verified"`
	Status     UserStatus `gorm:"type:varchar(20);default:'Active'" json:"status"`

	LastLoginAt       *int64 `json:"last_login_at,omitempty"`
	AccountLocked     bool   `gorm:"default:false" json:"account_locked"`
	PasswordChangedAt *int64 `json:"-"`
}

func (User) TableName() string {
	return "User"
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
		log.Printf("Password changed at %s:%v ", u.Username, u.PasswordChangedAt)
	} else {
		log.Println("Password field has not changed")
	}
	return nil
}

func (u *User) EmailValidation() (bool, error) {
	basicEmailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	email := strings.ToLower(u.Email)
	if len(email) < 3 || len(email) > 254 {
		return false, fmt.Errorf("email length is invalid")
	}
	return basicEmailRegex.MatchString(strings.ToLower(email)), fmt.Errorf("validation failed")
}

func (u *User) Validate() error {
	// Username validation
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !usernameRegex.MatchString(u.Username) {
		basicEmailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !basicEmailRegex.MatchString(u.Username) {
			return errors.New("username can only contain alphanumeric characters and underscores")
		}
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
	if !validRoles[u.Role] || u.Role == RoleAdmin {
		return errors.New("invalid user role")
	}

	// Status validation
	validStatuses := map[UserStatus]bool{
		StatusActive:  true,
		StatusDeleted: true,
		StatusBanned:  true,
	}
	if !validStatuses[u.Status] {
		return errors.New("invalid user status")
	}

	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
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

func (u *User) Lock(lock_time int) error {
	u.AccountLocked = true
	fmt.Println("Lock account", u.Username, "for", lock_time, "minutes")
	go func() {
		time.Sleep(time.Minute * time.Duration(lock_time))
		if u.AccountLocked {
			u.AccountLocked = false
		}
		fmt.Println("Account", u.Username, "is unlocked")
	}()
	return nil
}

func (u *User) Unlock() error {
	if u.AccountLocked {
		u.AccountLocked = false
		fmt.Println("Account unlocked")
		return nil
	}
	return fmt.Errorf("account is already unlocked")
}
