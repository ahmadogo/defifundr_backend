package crypto

import (
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)

	balance, err := GetBalance(configs.ContractAddress)
	require.NoError(t, err)
	require.NotEmpty(t, balance)

}
