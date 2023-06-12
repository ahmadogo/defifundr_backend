package crypto

import (
	"context"
	"math/big"

	"github.com/demola234/defiraise/gen"
	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

func CreateCampaign(title string, campaignType string, description string, goal int, deadline int, image string, privateKey string, address string) (string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		return "", err
	}

	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	cAdd := common.HexToAddress(configs.ContractAddress)

	tx, err := gen.NewGen(cAdd, client)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return "", err
	}

	auth.GasPrice = (gasPrice)
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce) + 100)

	tsx, err := tx.CreateCampaign(auth, campaignType, title, description, big.NewInt(int64(goal)), big.NewInt(int64(deadline)), image)
	if err != nil {
		return "", err
	}

	return tsx.Hash().Hex(), nil

}

func GetCampaign(id int, address string) (*Campaign, error) {
	configs, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return &Campaign{}, err
	}

	cAdd := common.HexToAddress(configs.ContractAddress)

	tx, err := gen.NewGen(cAdd, client)
	if err != nil {
		log.Err(err)
		return &Campaign{}, err
	}

	campaign, err := tx.GetCampaign(&bind.CallOpts{
		From: common.HexToAddress(address),
	},
		big.NewInt(int64(id)),
	)

	if err != nil {
		log.Err(err)
		return &Campaign{}, err
	}

	campaigns := Campaign{
		Title:        campaign.Title,
		CampaignType: campaign.CampaignType,
		Description:  campaign.Description,
		Goal:         campaign.Goal.Int64(),
		Deadline:     campaign.Deadline.Int64(),
		Image:        campaign.Image,
	}

	return &campaigns, nil
}

type Campaign struct {
	Title        string
	CampaignType string
	Description  string
	Goal         int64
	Deadline     int64
	Image        string
	ID           int64
}

func GetCampaigns(address string) ([]Campaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		return []Campaign{}, err
	}

	defer client.Close()

	cAdd := common.HexToAddress(configs.ContractAddress)

	tx, err := gen.NewGen(cAdd, client)
	if err != nil {
		return []Campaign{}, err
	}

	campaigns, err := tx.GetCampaigns(&bind.CallOpts{
		Pending: true,
		From:    common.HexToAddress(address),
	})

	if err != nil {
		return []Campaign{}, err
	}

	var campaignList []Campaign
	for _, campaign := range campaigns {
		campaignList = append(campaignList, Campaign{
			Title:        campaign.Title,
			CampaignType: campaign.CampaignType,
			Description:  campaign.Description,
			Goal:         campaign.Goal.Int64(),
			Deadline:     campaign.Deadline.Int64(),
			Image:        campaign.Image,
		})
	}

	return campaignList, nil
}
