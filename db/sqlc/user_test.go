package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) Users {
	password := "passphase"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	name:= utils.RandomString(6)

	arg := CreateUserParams{
		Username:       name,
		HashedPassword: hashedPassword,
		FirstName:      name,
		Email:          utils.RandomEmail(),
		Avatar:         utils.RandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FirstName, user.FirstName)
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

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FirstName, user2.FirstName)
	require.Equal(t, user.Avatar, user2.Avatar)
	require.Equal(t, user.Email, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestUpdateAvatar(t *testing.T) {
	user := CreateRandomUser(t)

	arg := UpdateAvatarParams{
		Username: user.Username,
		Avatar:   utils.RandomString(6),
	}

	user2, err := testQueries.UpdateAvatar(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FirstName, user2.FirstName)
	require.Equal(t, arg.Avatar, user2.Avatar)
	require.Equal(t, user.Email, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestChangePassword(t *testing.T) {
	user := CreateRandomUser(t)

	arg := ChangePasswordParams{
		Username:       user.Username,
		HashedPassword: utils.RandomString(6),
	}

	user2, err := testQueries.ChangePassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FirstName, user2.FirstName)
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
		FirstName: sql.NullString{
			String: utils.RandomString(6),
			Valid:  true,
		},
		Email: sql.NullString{
			String: utils.RandomEmail(),
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, arg.FirstName.String, user2.FirstName)
	require.Equal(t, arg.Email.String, user2.Email)

	require.NotZero(t, user2.Username)
	require.NotZero(t, user2.CreatedAt)
}

func TestCheckUsernameExists(t *testing.T) {
	user := CreateRandomUser(t)

	exists, err := testQueries.CheckUsernameExists(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, exists)
	require.Equal(t, true, exists)
}
