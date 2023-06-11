package crypto

import (
	"context"
	"fmt"
	"math/big"

	"github.com/demola234/defiraise/gen"
	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

func Deploy() {
	configs, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.Dial(configs.CryptoDeployURL)
	if err != nil {
		log.Fatal().Msg("cannot connect to ethereum network with the given url")
	}
	defer client.Close()
	account := common.HexToAddress(configs.ContractAddress)

	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		log.Fatal().Msg("cannot get nonce")
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal().Msg("cannot get gas price")
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal().Msg("cannot get chain id")
	}

	key, err := crypto.HexToECDSA(configs.ContractPrivateKey)
	if err != nil {
		log.Fatal().Msg("cannot get private key")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Fatal().Msg("cannot create auth")
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))

	hotel, ts, _, err := gen.DeployGen(auth, client)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Msg("cannot deploy contract")
	}

	fmt.Println("-----------------------------------")
	fmt.Println(hotel.Hex())
	fmt.Println(ts.Hash().Hex())
	fmt.Println("-----------------------------------")
}
