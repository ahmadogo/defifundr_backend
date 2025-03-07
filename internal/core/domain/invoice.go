package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvoiceStatus string

const (
	InvoiceStatusPending  InvoiceStatus = "pending"
	InvoiceStatusApproved InvoiceStatus = "approved"
	InvoiceStatusRejected InvoiceStatus = "rejected"
	InvoiceStatusPaid     InvoiceStatus = "paid"
)

type Invoice struct {
	ID               uuid.UUID       `json:"id"`
	FreelancerID     uuid.UUID       `json:"freelancer_id"`
	EmployerID       uuid.UUID       `json:"employer_id"`
	Amount           decimal.Decimal `json:"amount"`
	Currency         PayrollCurrency `json:"currency"`
	Status           InvoiceStatus   `json:"status"`
	ContractAddress  *string         `json:"contract_address,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	FreelancerName   string          `json:"freelancer_name,omitempty"`
	EmployerName     string          `json:"employer_name,omitempty"`
}

type CreateInvoiceParams struct {
	FreelancerID uuid.UUID       `json:"freelancer_id" validate:"required"`
	EmployerID   uuid.UUID       `json:"employer_id" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
	Currency     PayrollCurrency `json:"currency" validate:"required"`
}

type UpdateInvoiceStatusParams struct {
	ID     uuid.UUID     `json:"id" validate:"required"`
	Status InvoiceStatus `json:"status" validate:"required"`
}

type UpdateInvoiceContractParams struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	ContractAddress string    `json:"contract_address" validate:"required"`
}