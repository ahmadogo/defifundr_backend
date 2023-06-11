package db

import (
	"context"
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
		Avatar:         utils.RandomString(6),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.Avatar, user.Avatar)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user
}


func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}