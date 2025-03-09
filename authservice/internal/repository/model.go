package repository

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}
