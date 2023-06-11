package interfaces

import (
	"time"

	db "github.com/demola234/defiraise/db/sqlc"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FirstName         string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewUserResponse(user db.Users) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FirstName:         user.FirstName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
