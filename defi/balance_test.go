package defi

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// Token represents a token contract address and its name
type Token struct {
	Name    string
	Address string
}

// GetTokenBalance retrieves the balance of a specific token for a given address
func GetTokenBalance(client *ethclient.Client, tokenAddress string, userAddress string) (*big.Float, error) {
	// Create a new instance of the token contract
	token, err := NewDefi(common.HexToAddress(tokenAddress), client)
	if err != nil {
		return nil, err
	}

	// Fetch the balance
	balance, err := token.BalanceOf(&bind.CallOpts{Context: context.Background()}, common.HexToAddress(userAddress))
	if err != nil {
		return nil, err
	}

	// Convert balance to a human-readable format using token decimals
	decimals, err := token.Decimals(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return nil, err
	}
	decimalsFactor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	humanReadableBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), decimalsFactor)

	return humanReadableBalance, nil
}

// TestGetTokenBalances tests retrieving balances for USDT and USDC
func TestGetTokenBalances(t *testing.T) {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	require.NoError(t, err)

	userAddress := "0xYourEthereumAddressHere" // Replace with your Ethereum address

	tokens := []Token{
		{Name: "USDT", Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7"},
		{Name: "USDC", Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
	}

	for _, token := range tokens {
		balance, err := GetTokenBalance(client, token.Address, userAddress)
		require.NoError(t, err)
		t.Logf("%s Balance for address %s: %f", token.Name, userAddress, balance)
	}
}
