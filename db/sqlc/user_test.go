package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomString(6),
		HashedPassword: hashedPassword,
		FirstName:      utils.RandomString(6),
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
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

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
	user := createRandomUser(t)

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
	user := createRandomUser(t)

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
	user := createRandomUser(t)

	deleted, err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, deleted)

	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

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
	user := createRandomUser(t)

	exists, err := testQueries.CheckUsernameExists(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, exists)
	require.Equal(t, true, exists)
}
