package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Full_name string `json:"fullname"`
}

func (User) TableName() string {
	return "users"
}

