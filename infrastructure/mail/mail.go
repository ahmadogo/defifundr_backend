package mail

import (
	"time"

	"github.com/demola234/defifundr/internal/core/ports"
)

// EmailMessage represents an email message in the queue
type EmailMessage struct {
	ID           string                 `json:"id"`
	Recipient    string                 `json:"recipient"`
	Subject      string                 `json:"subject"`
	TemplateName string                 `json:"template_name"`
	Data         map[string]interface{} `json:"data"`
	Attachments  []ports.EmailAttachment `json:"attachments,omitempty"`
	Priority     ports.EmailPriority    `json:"priority"`
	CreatedAt    time.Time              `json:"created_at"`
}