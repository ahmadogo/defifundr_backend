package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypePayroll TransactionType = "payroll"
	TransactionTypeInvoice TransactionType = "invoice"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	TxHash    string            `json:"tx_hash"`
	Amount    decimal.Decimal   `json:"amount"`
	Currency  PayrollCurrency   `json:"currency"`
	Type      TransactionType   `json:"type"`
	Status    TransactionStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}

type CreateTransactionParams struct {
	UserID   uuid.UUID         `json:"user_id" validate:"required"`
	TxHash   string            `json:"tx_hash" validate:"required"`
	Amount   decimal.Decimal   `json:"amount" validate:"required"`
	Currency PayrollCurrency   `json:"currency" validate:"required"`
	Type     TransactionType   `json:"type" validate:"required"`
	Status   TransactionStatus `json:"status" validate:"required"`
}

type UpdateTransactionStatusParams struct {
	ID     uuid.UUID         `json:"id" validate:"required"`
	Status TransactionStatus `json:"status" validate:"required"`
}