package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeploy(t *testing.T) {
	address, err := Deploy()
	require.NoError(t, err)
	require.NotEmpty(t, address)
	if err != nil {
		require.Error(t, err)
	}
}
