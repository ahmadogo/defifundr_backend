package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	EmployerID uuid.UUID `json:"employer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Members    []Member  `json:"members,omitempty"`
}

type CreateOrganizationParams struct {
	Name       string    `json:"name" validate:"required"`
	EmployerID uuid.UUID `json:"employer_id" validate:"required"`
}

type UpdateOrganizationParams struct {
	Name string `json:"name" validate:"required"`
}

type Member struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	EmployeeID     uuid.UUID `json:"employee_id"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	EmployeeName   string    `json:"employee_name,omitempty"`
	EmployeeEmail  string    `json:"employee_email,omitempty"`
}

type AddMemberParams struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	EmployeeID     uuid.UUID `json:"employee_id" validate:"required"`
	Role           string    `json:"role" validate:"required"`
}

type UpdateMemberRoleParams struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	EmployeeID     uuid.UUID `json:"employee_id" validate:"required"`
	Role           string    `json:"role" validate:"required"`
}