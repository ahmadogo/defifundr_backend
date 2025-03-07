package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentFrequency string

const (
	FrequencyWeekly   PaymentFrequency = "weekly"
	FrequencyBiWeekly PaymentFrequency = "bi-weekly"
	FrequencyMonthly  PaymentFrequency = "monthly"
)

type PayrollStatus string

const (
	PayrollStatusPending   PayrollStatus = "pending"
	PayrollStatusActive    PayrollStatus = "active"
	PayrollStatusCompleted PayrollStatus = "completed"
)

type PayrollCurrency string

const (
	CurrencyUSDC PayrollCurrency = "USDC"
	CurrencySOL  PayrollCurrency = "SOL"
	CurrencyETH  PayrollCurrency = "ETH"
)

type Payroll struct {
	ID                uuid.UUID       `json:"id"`
	EmployerID        uuid.UUID       `json:"employer_id"`
	OrganizationID    *uuid.UUID      `json:"organization_id,omitempty"`
	PaymentFrequency  PaymentFrequency `json:"payment_frequency"`
	SalaryAmount      decimal.Decimal `json:"salary_amount"`
	Currency          PayrollCurrency `json:"currency"`
	ContractAddress   *string         `json:"contract_address,omitempty"`
	Status            PayrollStatus   `json:"status"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	Employees         []PayrollEmployee `json:"employees,omitempty"`
}

type PayrollEmployee struct {
	ID           uuid.UUID `json:"id"`
	PayrollID    uuid.UUID `json:"payroll_id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	CreatedAt    time.Time `json:"created_at"`
	EmployeeName string    `json:"employee_name,omitempty"`
	EmployeeEmail string   `json:"employee_email,omitempty"`
}

type CreatePayrollParams struct {
	EmployerID       uuid.UUID       `json:"employer_id" validate:"required"`
	OrganizationID   *uuid.UUID      `json:"organization_id,omitempty"`
	PaymentFrequency PaymentFrequency `json:"payment_frequency" validate:"required"`
	SalaryAmount     decimal.Decimal `json:"salary_amount" validate:"required"`
	Currency         PayrollCurrency `json:"currency" validate:"required"`
	EmployeeIDs      []uuid.UUID     `json:"employee_ids,omitempty"`
}

type UpdatePayrollStatusParams struct {
	ID     uuid.UUID     `json:"id" validate:"required"`
	Status PayrollStatus `json:"status" validate:"required"`
}

type UpdatePayrollContractParams struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	ContractAddress string    `json:"contract_address" validate:"required"`
}

type AddEmployeeToPayrollParams struct {
	PayrollID  uuid.UUID `json:"payroll_id" validate:"required"`
	EmployeeID uuid.UUID `json:"employee_id" validate:"required"`
}