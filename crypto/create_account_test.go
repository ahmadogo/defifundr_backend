package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	address, privateKey, err := CreateAddress()
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotEmpty(t, privateKey)
}

func TestGenerateAccountKeyStone(t *testing.T) {
	password := "password"
	filepath, address, err := GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotEmpty(t, filepath)

	private, public, err := DecryptPrivateKey(filepath, password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	// Check balance
	balance, err := GetBalance(address)
	require.NoError(t, err)
	require.NotEmpty(t, balance)
}
