package db

import (
	"context"
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func createVerifyEmails(t *testing.T) VerifyEmails {
	users := CreateRandomUser(t)

	arg := CreateVerifyEmailParams{
		Email:      users.Email,
		Username:   users.Username,
		SecretCode: utils.RandomOtp(),
	}

	verifyEmail, err := testQueries.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, verifyEmail)

	require.Equal(t, arg.Email, verifyEmail.Email)
	require.Equal(t, arg.Username, verifyEmail.Username)
	require.Equal(t, arg.SecretCode, verifyEmail.SecretCode)

	require.NotZero(t, verifyEmail.Email)
	require.NotZero(t, verifyEmail.SecretCode)
	require.NotZero(t, verifyEmail.Username)

	return verifyEmail
}

func TestCreateVerifyEmail(t *testing.T) {
	createVerifyEmails(t)
}
