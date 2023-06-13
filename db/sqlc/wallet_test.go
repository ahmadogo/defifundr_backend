package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/demola234/defiraise/crypto"
	"github.com/stretchr/testify/require"
)

func CreateWallet(t *testing.T) Wallet {
	users := CreateRandomUser(t)
	require.NotEmpty(t, users)
	password := "passphase"
	// Get Users Password from HashPassword

	fileName, address, err := crypto.GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, fileName)
	require.NotEmpty(t, address)

	arg := CreateWalletParams{
		Address:  address,
		Owner:    users.Username,
		Balance:  0,
		FilePath: fileName,
	}

	wallet, err := testQueries.CreateWallet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet)
	return wallet
}

func TestCreateWallet(t *testing.T) {
	CreateWallet(t)
}

func TestGetWallet(t *testing.T) {
	wallet := CreateWallet(t)
	require.NotEmpty(t, wallet)

	wallet2, err := testQueries.GetWallet(context.Background(), wallet.Owner)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	// Get Public Key from Private Key
	privateKey, public, err := crypto.DecryptPrivateKey(wallet2.FilePath, "passphase")
	require.NoError(t, err)
	require.NotEmpty(t, privateKey)
	require.NotEmpty(t, public)

	require.Equal(t, wallet.Address, wallet2.Address)
	require.Equal(t, wallet.Owner, wallet2.Owner)
	require.Equal(t, wallet.Balance, wallet2.Balance)
	require.Equal(t, wallet.FilePath, wallet2.FilePath)
	require.NotZero(t, wallet2.Address)
	require.NotZero(t, wallet2.Owner)
}

func TestGetInvalidUser(t *testing.T) {
	_, err := testQueries.GetWallet(context.Background(), "invalid")
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestInvalidPhase(t *testing.T) {
	wallet := CreateWallet(t)
	require.NotEmpty(t, wallet)

	wallet2, err := testQueries.GetWallet(context.Background(), wallet.Owner)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	invalidPhase := "invalid"

	// Get Public Key from Private Key
	privateKey, public, err := crypto.DecryptPrivateKey(wallet2.FilePath, invalidPhase)

	require.Error(t, err)
	require.EqualError(t, err, "could not decrypt key with given password")
	require.Empty(t, privateKey)
	require.Empty(t, public)

}
