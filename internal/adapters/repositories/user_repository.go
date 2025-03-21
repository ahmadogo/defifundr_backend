// internal/adapters/repositories/user_repository.go
package repositories

import (
	"context"

	db "github.com/demola234/defifundr/db/sqlc"
	utils "github.com/demola234/defifundr/infrastructure/hash"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	store db.Queries
}

func NewUserRepository(store db.Queries) *UserRepository {
	return &UserRepository{store: store}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(*user.Password)
	if err != nil {
		return nil, err
	}

	_, err = r.store.CreateUser(ctx, db.CreateUserParams{
		Email:               user.Email,
		PasswordHash:        pgtype.Text{String: hashedPassword, Valid: true},
		AccountType:         user.AccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PersonalAccountType: user.PersonalAccountType,
		Gender:              pgtype.Text{String: *user.Gender, Valid: user.Gender != nil},
		Nationality:         user.Nationality,
		ResidentialCountry:  pgtype.Text{String: *user.ResidentialCountry, Valid: user.ResidentialCountry != nil},
		JobRole:             pgtype.Text{String: *user.JobRole, Valid: user.JobRole != nil},
		CompanyWebsite:      pgtype.Text{String: *user.CompanyWebsite, Valid: user.CompanyWebsite != nil},
		EmploymentType:      pgtype.Text{String: *user.EmploymentType, Valid: user.EmploymentType != nil},
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	returnedUser := &domain.User{
		ID:                  user.ID,
		Email:               user.Email,
		Password:            &user.PasswordHash.String,
		AccountType:         user.AccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PersonalAccountType: user.PersonalAccountType,
		Gender:              &user.Gender.String,
		Nationality:         user.Nationality,
		ResidentialCountry:  &user.ResidentialCountry.String,
		JobRole:             &user.JobRole.String,
		CompanyWebsite:      &user.CompanyWebsite.String,
		EmploymentType:      &user.EmploymentType.String,
	}

	return returnedUser, nil

}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	returnedUser := &domain.User{
		ID:                  user.ID,
		Email:               user.Email,
		Password:            &user.PasswordHash.String,
		AccountType:         user.AccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PersonalAccountType: user.PersonalAccountType,
		Gender:              &user.Gender.String,
		Nationality:         user.Nationality,
		ResidentialCountry:  &user.ResidentialCountry.String,
		JobRole:             &user.JobRole.String,
		CompanyWebsite:      &user.CompanyWebsite.String,
		EmploymentType:      &user.EmploymentType.String,
	}

	return returnedUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	updatedUser, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:                  user.ID,
		Email:               user.Email,
		AccountType:         user.AccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PersonalAccountType: user.PersonalAccountType,
		Gender:              pgtype.Text{String: *user.Gender, Valid: user.Gender != nil},
		Nationality:         user.Nationality,
		ResidentialCountry:  pgtype.Text{String: *user.ResidentialCountry, Valid: user.ResidentialCountry != nil},
		JobRole:             pgtype.Text{String: *user.JobRole, Valid: user.JobRole != nil},
		CompanyWebsite:      pgtype.Text{String: *user.CompanyWebsite, Valid: user.CompanyWebsite != nil},
		EmploymentType:      pgtype.Text{String: *user.EmploymentType, Valid: user.EmploymentType != nil},
	})

	if err != nil {
		return nil, err
	}

	returnedUser := &domain.User{
		ID:    user.ID,
		Email: user.Email,

		AccountType:         user.AccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PersonalAccountType: user.PersonalAccountType,
		Nationality:         user.Nationality,
		ResidentialCountry:  user.ResidentialCountry,
		JobRole:             user.JobRole,
		CompanyWebsite:      user.CompanyWebsite,
		EmploymentType:      user.EmploymentType,
		CreatedAt:           updatedUser.CreatedAt,
		UpdatedAt:           updatedUser.UpdatedAt,
	}

	return returnedUser, nil

}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	err := r.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}
