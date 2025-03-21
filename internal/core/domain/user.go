package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `json:"id"`
	Email               string    `json:"email"`
	Password            *string    `json:"-"`
	ProfilePicture      *string   `json:"profile_picture,omitempty"`
	AccountType         string    `json:"account_type"`
	Gender              *string   `json:"gender,omitempty"`
	PersonalAccountType string    `json:"personal_account_type"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Nationality         string    `json:"nationality"`
	ResidentialCountry  *string   `json:"residential_country,omitempty"`
	JobRole             *string   `json:"job_role,omitempty"`
	CompanyWebsite      *string   `json:"company_website,omitempty"`
	EmploymentType      *string   `json:"employment_type,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
