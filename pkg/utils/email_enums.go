package utils

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// EmailPriority represents the priority level of an email
type EmailPriority int

const (
	LowPriority      EmailPriority = 1
	NormalPriority   EmailPriority = 2
	HighPriority     EmailPriority = 3
	CriticalPriority EmailPriority = 4
)
