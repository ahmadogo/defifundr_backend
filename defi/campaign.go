package defi

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/demola234/defiraise/gen"
	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

var (
	// DeployedContractAddress is the address of the deployed contract
	Address = "0x574Bc33136180f0734fc3fa55379e9e28701395E"
)

func CreateCampaign(title string, campaignType string, description string, goal float64, deadline time.Time, image string, privateKey *ecdsa.PrivateKey, address string) (string, error) {
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

	cAdd := common.HexToAddress(Address)

	tx, err := gen.NewGen(cAdd, client)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}

	auth.GasPrice = (gasPrice)
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))

	goals := big.NewInt(int64(goal * 1e18))

	tsx, err := tx.CreateCampaign(auth, campaignType, title, description, goals, big.NewInt(deadline.Unix()), image)
	if err != nil {
		return "", err
	}

	fmt.Println("-----------------------------------")
	fmt.Println("tx view: ", tx)
	fmt.Println("............Loading............")
	fmt.Println("-----------------------------------")

	return tsx.Hash().Hex(), nil

}

func GetCampaign(id int) (*Campaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return &Campaign{}, err
	}

	cAdd := common.HexToAddress(Address)

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
		ID:           campaign.Id.Int64(),
		TotalFunds:   campaign.TotalFunds.Int64(),
		Owner:        campaign.Owner.Hex(),
	}

	return &campaigns, nil
}

type Campaign struct {
	Title                  string
	CampaignType           string
	Description            string
	Goal                   int64
	Deadline               int64
	Image                  string
	ID                     int64
	TotalFunds             int64
	Owner                  string
	TotalNumberOfDonations int64
}

type CampaignCategory struct {
	Name        string
	Image       string
	Description string
	ID          string
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

	cAdd := common.HexToAddress(Address)

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
			TotalFunds:   campaign.TotalFunds.Int64(),
			Owner:        campaign.Owner.Hex(),
			ID:           campaign.Id.Int64(),
		}

		campaignList = append(campaignList, campaigns)
	}

	return campaignList, nil
}

func Donate(amount float64, id int, privateKey *ecdsa.PrivateKey, address string) (string, error) {
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

	cAdd := common.HexToAddress(Address)

	tx, err := gen.NewGenTransactor(cAdd, client)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}

	// convert amount to wei
	amount = amount * 1e18

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

func GetDonations(id int) ([]common.Address, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress(Address)

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

func GetDonorsAddressesAndAmounts(id int) ([]string, []int64, *big.Int, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, nil, nil, err
	}

	cAdd := common.HexToAddress(Address)

	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, nil, nil, err
	}

	donators, amounts, totalFunds, err := tx.GetDonorsAddressesAndAmounts(&bind.CallOpts{},
		big.NewInt(int64(id)),
	)
	if err != nil {
		log.Err(err)
		return nil, nil, nil, err
	}

	var donatorsList []string

	for _, donator := range donators {
		donators := donator.Hex()

		donatorsList = append(donatorsList, donators)
	}

	var amountsList []int64

	for _, amount := range amounts {
		amounts := amount.Int64()

		amountsList = append(amountsList, amounts)
	}

	return donatorsList, amountsList, totalFunds, nil
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

	cAdd := common.HexToAddress(Address)

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

func GetCampaignsByOwner(address string) ([]Campaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	cAdd := common.HexToAddress(Address)

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

	var campaignList []Campaign

	for _, campaign := range campaigns {
		campaigns := Campaign{
			Title:        campaign.Title,
			CampaignType: campaign.CampaignType,
			Description:  campaign.Description,
			Goal:         campaign.Goal.Int64(),
			Deadline:     campaign.Deadline.Int64(),
			Image:        campaign.Image,
			TotalFunds:   campaign.TotalFunds.Int64(),
			Owner:        campaign.Owner.Hex(),
			ID:           campaign.Id.Int64(),
		}

		campaignList = append(campaignList, campaigns)
	}

	return campaignList, nil
}

func GetTotalDonationsByCampaignId(id int) (*big.Int, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	cAdd := common.HexToAddress(Address)
	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	totalDonations, err := tx.GetTotalDonationsByCampaignId(&bind.CallOpts{},
		big.NewInt(int64(id)),
	)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	return totalDonations, nil
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

	cAdd := common.HexToAddress(Address)

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

	cAdd := common.HexToAddress(Address)

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

func CreateCategories(categoryName string, description string, image string, privateKey *ecdsa.PrivateKey, address string) (string, error) {
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
	cAdd := common.HexToAddress(Address)
	tx, err := gen.NewGenTransactor(cAdd, client)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.GasPrice = (gasPrice)
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))
	tsx, err := tx.CreateCategory(auth, categoryName, description, image)
	if err != nil {
		return "", err
	}
	fmt.Println("-----------------------------------")
	fmt.Println("tx view: ", tx)
	fmt.Println("............Loading............")
	fmt.Println("-----------------------------------")
	return tsx.Hash().Hex(), nil
}

func GetCategories() ([]CampaignCategory, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	cAdd := common.HexToAddress(Address)
	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	categories, err := tx.GetCategories(&bind.CallOpts{})
	if err != nil {
		log.Err(err)
		return nil, err
	}
	var categoryList []CampaignCategory
	for _, category := range categories {
		campaigns := CampaignCategory{
			Name:        category.Name,
			Image:       category.Image,
			Description: category.Description,
			ID:          category.Id.String(),
		}

		categoryList = append(categoryList, campaigns)

	}

	return categoryList, nil
}

func SearchCampaigns(name string) ([]Campaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	cAdd := common.HexToAddress(Address)
	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	campaigns, err := tx.SearchCampaignByName(&bind.CallOpts{}, name)
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
			TotalFunds:   campaign.TotalFunds.Int64(),
			Owner:        campaign.Owner.Hex(),
			ID:           campaign.Id.Int64(),
		}

		campaignList = append(campaignList, campaigns)
	}

	return campaignList, nil
}

func GetCampaignByCategory(categoryId int64) ([]Campaign, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	client, err := ethclient.DialContext(context.Background(), configs.CryptoDeployURL)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	cAdd := common.HexToAddress(Address)
	tx, err := gen.NewGenCaller(cAdd, client)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	campaigns, err := tx.GetCampaignsByCategory(&bind.CallOpts{}, big.NewInt(categoryId))
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
			TotalFunds:   campaign.TotalFunds.Int64(),
			Owner:        campaign.Owner.Hex(),
			ID:           campaign.Id.Int64(),
		}

		campaignList = append(campaignList, campaigns)
	}

	return campaignList, nil
}
