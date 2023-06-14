package crypto

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/demola234/defiraise/gen"
	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

func CreateCampaign(title string, campaignType string, description string, goal int, deadline int, image string, privateKey *ecdsa.PrivateKey, address string) (*bind.TransactOpts, string, error, *Campaign) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		return nil, "", err, nil
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		return nil, "", err, nil
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, "", err, nil
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, "", err, nil
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGen(cAdd, client)
	if err != nil {
		return nil, "", err, nil
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, "", err, nil
	}

	auth.GasPrice = (gasPrice)
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))

	tsx, err := tx.CreateCampaign(auth, campaignType, title, description, big.NewInt(int64(goal)), big.NewInt(int64(deadline)), image)
	if err != nil {
		return nil, "", err, nil
	}

	campaign, err := tx.GetCampaign(&bind.CallOpts{},
		big.NewInt(int64(0)),
	)

	if err != nil {
		log.Err(err)
		return nil, "", err, nil
	}

	campaigns := Campaign{
		Title:        campaign.Title,
		CampaignType: campaign.CampaignType,
		Description:  campaign.Description,
		Goal:         campaign.Goal.Int64(),
		Deadline:     campaign.Deadline.Int64(),
		Image:        campaign.Image,
	}

	log.Info().Msgf("campaign: %+v", campaigns)

	fmt.Println("-----------------------------------")
	fmt.Println("tx view: ", tx)
	fmt.Println("............Loading............")
	fmt.Println("-----------------------------------")

	return auth, tsx.Hash().Hex(), nil, &campaigns

}

func GetCampaign(id int, address string) (*Campaign, error) {
	configs, err := utils.LoadConfig("./../")
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

	campaign, err := tx.GetCampaign(&bind.CallOpts{},
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
