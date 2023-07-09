package defi
import (
	"context"
	"math"
	"math/big"

	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

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
