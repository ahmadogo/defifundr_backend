package defi

import (
	"context"
	"math"
	"math/big"

	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// Stablecoin represents a stablecoin contract address and name
type Stablecoin struct {
	Name    string
	Address string
}

// GetStablecoinBalances retrieves the balances of stablecoins for a given address
func GetStablecoinBalances(userAddress string, stablecoins []Stablecoin) (map[string]string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	accountAddress := common.HexToAddress(userAddress)
	balances := make(map[string]string)

	for _, coin := range stablecoins {
		tokenAddress := common.HexToAddress(coin.Address)

		// Load the token contract
		token, err := NewDefi(tokenAddress, client)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to load token contract for %s", coin.Name)
			continue
		}

		// Get the balance of the user
		balance, err := token.BalanceOf(&bind.CallOpts{Context: context.Background()}, accountAddress)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get balance for %s", coin.Name)
			continue
		}

		// Convert from wei to eth
		fbBalance := new(big.Float).SetInt(balance)
		ethValue := new(big.Float).Quo(fbBalance, big.NewFloat(math.Pow10(18)))

		// Store the balance in the map
		balances[coin.Name] = ethValue.String()
	}

	return balances, nil
}

func GetBalance(address string) (string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		return "", err
	}

	defer client.Close()
	accountAt := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), accountAt, nil)
	if err != nil {
		return "", err
	}
	fbBalance := new(big.Float)
	fbBalance.SetString(balance.String())

	// Convert from wei to eth
	ethValue := new(big.Float).Quo(fbBalance, big.NewFloat(math.Pow10(18)))
	return ethValue.String(), nil
}
