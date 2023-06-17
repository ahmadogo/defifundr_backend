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

func CreateCampaign(title string, campaignType string, description string, goal int, deadline int, image string, privateKey *ecdsa.PrivateKey, address string) (*bind.TransactOpts, string, *Campaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		return nil, "", nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		return nil, "", nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, "", nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, "", nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGen(cAdd, client)
	if err != nil {
		return nil, "", nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, "", nil, err
	}

	auth.GasPrice = (gasPrice)
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce) + 100)

	tsx, err := tx.CreateCampaign(auth, campaignType, title, description, big.NewInt(int64(goal)), big.NewInt(int64(deadline)), image)
	if err != nil {
		return nil, "", nil, err
	}

	campaign, err := tx.GetCampaign(&bind.CallOpts{},
		big.NewInt(int64(3)),
	)

	if err != nil {
		log.Err(err)
		return nil, "", nil, err
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

	return auth, tsx.Hash().Hex(),&campaigns, nil

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

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

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
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	campaigns, err := tx.GetCampaigns(&bind.CallOpts{})
	if err != nil {
		log.Err(err)
		return nil, err
	}

	var campaignList []Campaign

	for _, campaign := range campaigns {
		campaigns := Campaign{
			Title:        campaign.Title,
			CampaignType: campaign.CampaignType,
			Description:  campaign.Description,
			Goal:         campaign.Goal.Int64(),
			Deadline:     campaign.Deadline.Int64(),
			Image:        campaign.Image,
		}

		campaignList = append(campaignList, campaigns)
	}

	return campaignList, nil
}

func Donate(amount float32, id int, privateKey *ecdsa.PrivateKey, address string) (string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
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

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenTransactor(cAdd, client)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}

	// convert amount to wei
	amount = amount * 1000000000000000000

	

	auth.GasPrice = (gasPrice)
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(int64(amount))
	auth.From = common.HexToAddress(address)

	tsx, err := tx.Donate(auth, big.NewInt(int64(id)))
	if err != nil {
		return "", err
	}

	fmt.Println("-----------------------------------")
	fmt.Println("tx view: ", tx)
	fmt.Println("............Loading............")
	fmt.Println("-----------------------------------")

	return tsx.Hash().Hex(), nil
}

func GetDonations(id int, address string) ([]common.Address, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	donations, err := tx.GetCampaignDonators(&bind.CallOpts{},
		big.NewInt(int64(id)),
	)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	var donationList []common.Address

	for _, donation := range donations {
		donations := []common.Address{
			donation,
		}

		donationList = append(donationList, donations...)
	}

	return donationList, nil
}

func GetDonorsAddressesAndAmounts(id int, address string) ([]common.Address, []*big.Int, *big.Int, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, nil, nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, nil, nil, err
	}

	donators, donations, totalFunds, err := tx.GetDonorsAddressesAndAmounts(&bind.CallOpts{},
		big.NewInt(int64(id)),
	)
	if err != nil {
		log.Err(err)
		return nil, nil, nil, err
	}

	return donators, donations, totalFunds, nil
}

func GetCampaignTypes(address string) ([]gen.CrowdFundingCampaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	campaignTypes, err := tx.GetCampaigns(&bind.CallOpts{})
	if err != nil {
		log.Err(err)
		return nil, err
	}

	var campaignTypeList []gen.CrowdFundingCampaign

	for _, campaignType := range campaignTypes {
		campaignTypes := campaignType

		campaignTypeList = append(campaignTypeList, campaignTypes)
	}

	return campaignTypeList, nil
}

func GetCampaignsByOwner(address string) ([]gen.CrowdFundingCampaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	campaigns, err := tx.GetCampaignsByOwner(&bind.CallOpts{},
		common.HexToAddress(address),
	)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	var campaignList []gen.CrowdFundingCampaign

	for _, campaign := range campaigns {
		campaigns := campaign

		campaignList = append(campaignList, campaigns)
	}

	return campaignList, nil
}

func GetCampaignByType(campaignName string) ([]gen.CrowdFundingCampaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	campaign, err := tx.GetCampaignsByType(&bind.CallOpts{},
		campaignName,
	)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	return campaign, nil
}

func PayOut(id int, address string, privateKey *ecdsa.PrivateKey) (string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return "", err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Err(err)
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Err(err)
		return "", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Err(err)
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(1000000)
	auth.GasPrice = gasPrice

	tx, err := gen.NewGenTransactor(cAdd, client)
	if err != nil {
		log.Err(err)
		return "", err
	}

	tsx, err := tx.PayOut(auth,
		big.NewInt(int64(id)),
	)
	if err != nil {
		log.Err(err)
		return "", err
	}

	fmt.Println("-----------------------------------")
	fmt.Println("tx view: ", tx)
	fmt.Println("............Loading............")
	fmt.Println("-----------------------------------")

	return tsx.Hash().Hex(), nil
}

func SendBackDonations(id int, address string, privateKey *ecdsa.PrivateKey) (string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return "", err
	}

	cAdd := common.HexToAddress("0xd9d4b660f51eb66b3f8b3829012424e46186857f")

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Err(err)
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Err(err)
		return "", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Err(err)
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(1000000)
	auth.GasPrice = gasPrice

	tx, err := gen.NewGenTransactor(cAdd, client)
	if err != nil {
		log.Err(err)
		return "", err
	}

	tsx, err := tx.SendBackDonations(auth,
		big.NewInt(int64(id)),
	)
	if err != nil {
		log.Err(err)
		return "", err
	}

	fmt.Println("-----------------------------------")
	fmt.Println("tx view: ", tx)
	fmt.Println("............Loading............")
	fmt.Println("-----------------------------------")

	return tsx.Hash().Hex(), nil
}
