package interfaces

import (
	"time"

	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/google/uuid"
)

var ErrUserNotFound = "user not found"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	Address           string    `json:"address"`
	Balance           string     `json:"balance"`
}

func NewUserResponse(user db.Users) UserResponse {
	return UserResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		Address:           user.Address,
		Balance:           user.Balance,
	}
}

type GetUserRequest struct {
	Username string `json:"username" binding:"required"`
}

type VerifyUserRequest struct {
	Username string `json:"username" binding:"required"`
	OtpCode  string `json:"otp_code" binding:"required"`
}

type ResendVerificationCodeRequest struct {
	Username string `json:"username" binding:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}


type ResetPasswordRequest struct {
	Username string `json:"username" binding:"required"`
}


type VerifyUserResetRequest struct {
	Username string `json:"username" binding:"required"`
	OtpCode  string `json:"otp_code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CheckUsernameExistsRequest struct {
	Username string `json:"username" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}