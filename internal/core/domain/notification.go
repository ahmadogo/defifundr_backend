package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypePayroll     NotificationType = "payroll"
	NotificationTypeInvoice     NotificationType = "invoice"
	NotificationTypeTransaction NotificationType = "transaction"
)

type Notification struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"user_id"`
	Message   string           `json:"message"`
	Type      NotificationType `json:"type"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

type CreateNotificationParams struct {
	UserID  uuid.UUID        `json:"user_id" validate:"required"`
	Message string           `json:"message" validate:"required"`
	Type    NotificationType `json:"type" validate:"required"`
}

type NotificationListParams struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Limit  int32     `json:"limit" validate:"required,min=1,max=100"`
	Offset int32     `json:"offset" validate:"min=0"`
}