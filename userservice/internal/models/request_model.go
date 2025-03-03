package models

type UserCreationParams struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type UserUpdateParams struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Picture  string `json:"picture"`
	Status   string `json:"status"`
}
