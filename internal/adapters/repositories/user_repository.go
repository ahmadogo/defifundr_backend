package repositories

import (
	"context"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository struct implements the repository interface for users
type UserRepository struct {
	store db.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(store db.Queries) *UserRepository {
	return &UserRepository{store: store}
}

// Helper function to safely handle string pointers
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// RegisterUser implements the user registration functionality
func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	// Convert domain user to database parameters
	params := db.CreateUserParams{
		ID:                  user.ID,
		Email:               user.Email,
		PasswordHash:        pgtype.Text{String: user.PasswordHash, Valid: user.PasswordHash != ""},
		ProfilePicture:      pgtype.Text{String: getStringValue(user.ProfilePicture), Valid: user.ProfilePicture != nil},
		AccountType:         user.AccountType,
		Gender:              pgtype.Text{String: getStringValue(user.Gender), Valid: user.Gender != nil},
		PersonalAccountType: user.PersonalAccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Nationality:         user.Nationality,
		ResidentialCountry:  pgtype.Text{String: getStringValue(user.ResidentialCountry), Valid: user.ResidentialCountry != nil},
		JobRole:             pgtype.Text{String: getStringValue(user.JobRole), Valid: user.JobRole != nil},
		CompanyName:         pgtype.Text{String: user.CompanyName, Valid: user.CompanyName != ""},
		CompanyAddress:      pgtype.Text{String: user.CompanyAddress, Valid: user.CompanyAddress != ""},
		CompanyCity:         pgtype.Text{String: user.CompanyCity, Valid: user.CompanyCity != ""},
		CompanyPostalCode:   pgtype.Text{String: user.CompanyPostalCode, Valid: user.CompanyPostalCode != ""},
		CompanyCountry:      pgtype.Text{String: user.CompanyCountry, Valid: user.CompanyCountry != ""},
		AuthProvider:        pgtype.Text{String: user.AuthProvider, Valid: user.AuthProvider != ""},
		ProviderID:          pgtype.Text{String: user.ProviderID, Valid: user.ProviderID != ""},
		EmployeeType:        pgtype.Text{String: user.EmployeeType, Valid: user.EmployeeType != ""},
		CompanyWebsite:      pgtype.Text{String: getStringValue(user.CompanyWebsite), Valid: user.CompanyWebsite != nil},
		EmploymentType:      pgtype.Text{String: getStringValue(user.EmploymentType), Valid: user.EmploymentType != nil},
		CreatedAt:           pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:           pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	// Call the database to create user
	dbUser, err := r.store.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	// Map database user back to domain user
	return mapDBUserToDomainUser(dbUser), nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	dbUser, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdateUser updates a user's information
func (r *UserRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	params := db.UpdateUserParams{
		ID:                  user.ID,
		Email:               user.Email,
		ProfilePicture:      pgtype.Text{String: getStringValue(user.ProfilePicture), Valid: user.ProfilePicture != nil},
		AccountType:         user.AccountType,
		Gender:              pgtype.Text{String: getStringValue(user.Gender), Valid: user.Gender != nil},
		PersonalAccountType: user.PersonalAccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Nationality:         user.Nationality,
		ResidentialCountry:  pgtype.Text{String: getStringValue(user.ResidentialCountry), Valid: user.ResidentialCountry != nil},
		JobRole:             pgtype.Text{String: getStringValue(user.JobRole), Valid: user.JobRole != nil},
		CompanyWebsite:      pgtype.Text{String: getStringValue(user.CompanyWebsite), Valid: user.CompanyWebsite != nil},
		EmploymentType:      pgtype.Text{String: getStringValue(user.EmploymentType), Valid: user.EmploymentType != nil},
		CompanyName:         pgtype.Text{String: user.CompanyName, Valid: user.CompanyName != ""},
		CompanyAddress:      pgtype.Text{String: user.CompanyAddress, Valid: user.CompanyAddress != ""},
		CompanyCity:         pgtype.Text{String: user.CompanyCity, Valid: user.CompanyCity != ""},
		CompanyPostalCode:   pgtype.Text{String: user.CompanyPostalCode, Valid: user.CompanyPostalCode != ""},
		CompanyCountry:      pgtype.Text{String: user.CompanyCountry, Valid: user.CompanyCountry != ""},
		AuthProvider:        pgtype.Text{String: user.AuthProvider, Valid: user.AuthProvider != ""},
		ProviderID:          pgtype.Text{String: user.ProviderID, Valid: user.ProviderID != ""},
	}

	dbUser, err := r.store.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	params := db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	}

	return r.store.UpdateUserPassword(ctx, params)
}

// CheckEmailExists checks if an email already exists
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.store.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Helper function to get a string pointer
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function to get a string from pgtype.Text
func getTextString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// Helper function to map database user to domain user
func mapDBUserToDomainUser(dbUser db.Users) *domain.User {
	var passwordPtr *string
	if dbUser.PasswordHash.Valid {
		passwordPtr = &dbUser.PasswordHash.String
	}

	return &domain.User{
		ID:                  dbUser.ID,
		Email:               dbUser.Email,
		Password:            passwordPtr,
		PasswordHash:        getTextString(dbUser.PasswordHash),
		ProfilePicture:      strPtr(getTextString(dbUser.ProfilePicture)),
		AccountType:         dbUser.AccountType,
		Gender:              strPtr(getTextString(dbUser.Gender)),
		PersonalAccountType: dbUser.PersonalAccountType,
		FirstName:           dbUser.FirstName,
		LastName:            dbUser.LastName,
		Nationality:         dbUser.Nationality,
		ResidentialCountry:  strPtr(getTextString(dbUser.ResidentialCountry)),
		JobRole:             strPtr(getTextString(dbUser.JobRole)),
		CompanyName:         getTextString(dbUser.CompanyName),
		CompanyAddress:      getTextString(dbUser.CompanyAddress),
		CompanyCity:         getTextString(dbUser.CompanyCity),
		CompanyPostalCode:   getTextString(dbUser.CompanyPostalCode),
		CompanyCountry:      getTextString(dbUser.CompanyCountry),
		AuthProvider:        getTextString(dbUser.AuthProvider),
		ProviderID:          getTextString(dbUser.ProviderID),
		EmployeeType:        getTextString(dbUser.EmployeeType),
		CompanyWebsite:      strPtr(getTextString(dbUser.CompanyWebsite)),
		EmploymentType:      strPtr(getTextString(dbUser.EmploymentType)),
		// Fill in missing fields with empty values
		Address:      "",
		City:         "",
		PostalCode:   "",
		WebAuthToken: "",
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}
