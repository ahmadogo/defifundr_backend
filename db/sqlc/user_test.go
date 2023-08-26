package db

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) Users {
	password := "passphase"

	filepath, address, err := GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotEmpty(t, filepath)
	otpCode := utils.RandomOtp()

	name := utils.RandomString(6)

	arg := CreateUserParams{

		Username:   name,
		Email:      utils.RandomEmail(),
		Avatar:     utils.RandomString(6),
		Balance:    "0",
		Address:    address,
		SecretCode: otpCode,
		FilePath:   filepath,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
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
		Username: sql.NullString{
			String: user.Username,
			Valid:  true,
		},
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

		Username: sql.NullString{
			String: user.Username,
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
		Username: sql.NullString{
			String: oldUser.Username,
			Valid:  true,
		},
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
		Username: sql.NullString{
			String: oldUser.Username,
			Valid:  true,
		},
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
		Username: sql.NullString{
			String: oldUser.Username,
			Valid:  true,
		},
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

func GenerateAccountKeyStone(password string) (string, string, error) {
	// Generate a new random private key
	key := keystore.NewKeyStore("./../tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	// Create a new account with the specified encryption passphrase
	passwordKey := password
	account, err := key.NewAccount(passwordKey)

	if err != nil {
		log.Error().Err(err).Msg("cannot create account")
	}

	filename := account.URL.Path[strings.LastIndex(account.URL.Path, "/")+1:]

	accountName := account.Address.Hex()

	return filename, accountName, nil
}




// func TestCreateUserPassword(t *testing.T) {
// 	user := CreateRandomUser(t)

// 	arg := CreateUserPasswordParams{
// 		Username:          user.Username,
// 		HashedPassword:    utils.RandomString(6),
// 		PasswordChangedAt: time.Now(),
// 		IsUsed:            true,
// 	}

// 	user2, err := testQueries.CreateUserPassword(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, user2)

// 	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
// 	require.Equal(t, arg.Username, user2.Username)
// 	require.Equal(t, arg.IsUsed, user2.IsUsed)

// 	require.NotZero(t, user2.Username)
// 	require.NotZero(t, user2.CreatedAt)
// }
