package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEthPrice(t *testing.T) {
	price, err := GetEthPrice()
	require.NoError(t, err)
	require.NotEmpty(t, price)
}
