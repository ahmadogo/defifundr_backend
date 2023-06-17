package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/demola234/defiraise/crypto"
	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) Users {
	password := "passphase"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	filepath, address, err := crypto.GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotEmpty(t, filepath)
	otpCode := utils.RandomOtp()

	name := utils.RandomString(6)

	arg := CreateUserParams{
		HashedPassword: hashedPassword,
		Username:       name,
		Email:          utils.RandomEmail(),
		Avatar:         utils.RandomString(6),
		Balance:        "0",
		Address:        address,
		SecretCode:     otpCode,
		FilePath:       filepath,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Avatar, user.Avatar)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.Avatar, user2.Avatar)
	require.Equal(t, user.Email, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestUpdateAvatar(t *testing.T) {
	user := CreateRandomUser(t)

	arg := UpdateUserParams{Avatar: sql.NullString{
		String: utils.RandomEmail(),
		Valid:  true,
	},
		Username: user.Username,
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, arg.Avatar.String, user2.Avatar)
	require.Equal(t, user.Email, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestChangePassword(t *testing.T) {
	user := CreateRandomUser(t)

	arg := ChangePasswordParams{
		Username:          user.Username,
		HashedPassword:    utils.RandomString(6),
		PasswordChangedAt: time.Now(),
	}

	user2, err := testQueries.ChangePassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.Avatar, user2.Avatar)
	require.Equal(t, user.Email, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestDeleteUser(t *testing.T) {
	user := CreateRandomUser(t)

	deleted, err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, deleted)

	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	user := CreateRandomUser(t)

	arg := UpdateUserParams{

		Username: user.Username,
		Email: sql.NullString{
			String: utils.RandomEmail(),
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Email.String, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newPassword := utils.RandomString(6)
	newHashedPassword, err := utils.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newEmail := utils.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyAvatar(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newAvatar := utils.RandomString(6)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Avatar: sql.NullString{
			String: newAvatar,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Avatar, updatedUser.Avatar)
	require.Equal(t, newAvatar, updatedUser.Avatar)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestCheckUsernameExist(t *testing.T) {
	user := CreateRandomUser(t)

	exist, err := testQueries.CheckUsernameExists(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, true, exist)
}
