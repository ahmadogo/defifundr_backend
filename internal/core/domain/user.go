package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleEmployer   UserRole = "employer"
	RoleEmployee   UserRole = "employee"
	RoleFreelancer UserRole = "freelancer"
	RoleAdmin      UserRole = "admin"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           UserRole  `json:"role"`
	Name           string    `json:"name"`
	JobRole        *string   `json:"job_role,omitempty"`
	CompanyWebsite *string   `json:"company_website,omitempty"`
	Country        *string   `json:"country,omitempty"`
	EmploymentType *string   `json:"employment_type,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Wallets        []Wallet  `json:"wallets,omitempty"`
}

type CreateUserParams struct {
	Email          string   `json:"email" validate:"required,email"`
	Password       string   `json:"password" validate:"required,min=8"`
	Role           UserRole `json:"role" validate:"required"`
	Name           string   `json:"name" validate:"required"`
	JobRole        *string  `json:"job_role,omitempty"`
	CompanyWebsite *string  `json:"company_website,omitempty"`
	Country        *string  `json:"country,omitempty"`
	EmploymentType *string  `json:"employment_type,omitempty"`
}

type UpdateUserParams struct {
	Name           *string `json:"name,omitempty"`
	JobRole        *string `json:"job_role,omitempty"`
	CompanyWebsite *string `json:"company_website,omitempty"`
	Country        *string `json:"country,omitempty"`
	EmploymentType *string `json:"employment_type,omitempty"`
}

type LoginUserParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
