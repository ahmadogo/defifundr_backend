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
