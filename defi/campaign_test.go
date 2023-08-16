package defi

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func createCampaign(t *testing.T) (*ecdsa.PrivateKey, string, error) {
	t.Parallel()
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	password := "passphase"
	filepath, address, err := GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, filepath)
	require.NotEmpty(t, address)

	private, public, err := DecryptPrivateKey("UTC--2023-08-05T08-49-36.197726000Z--9616c35e6042a3c008c0f2badedcdc84fd7eb8b0", password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	title := "Fund Ikorudu Child Education"
	description := "13 years ago I undertook the restoration of a former Nuclear Monitoring Post.   Our aim is to teach visitors about the Cold War and how a Nuclear War would have affected the island of Ireland.   We do not charge visitors an entrance fee and rely on donations to keep our museum totally free.   Moving forward we want to reach out to Schools and other institutions and bring our collection to them. A successful campaign will allow us to purchase a trailer which means we can bring our collection anywhere in the country. Imagine that!!  Any donation, big or small, will help us keep the museum free for years to come and allow us to teach as many people as possible about the dangers of nuclear weapons."
	image := "https://www.qgiv.com/blog/wp-content/uploads/2023/01/C_pexels-rodnae-productions-7551758-1-1-1-1-1-1-1-1-1-1-1-1-1-300x200.jpg"
	goal := 0.005
	deadline := time.Now().AddDate(0, 0, 1)
	campaignType := "Education"

	campaign, err := CreateCampaign(title, campaignType, description, goal, deadline, image, private, "0x9616c35e6042a3c008c0f2badedcdc84fd7eb8b0")
	if err != nil {
		return nil, "", err
	}

	return private, campaign, nil
}

func TestCreateCampaign(t *testing.T) {

	c, rx, err := createCampaign(t)

	require.NoError(t, err)

	require.NotEmpty(t, c)
	require.NotEmpty(t, rx)

}

func TestGetCampaign(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	campaign, err := GetCampaign(0)
	require.NoError(t, err)
	require.NotEmpty(t, campaign)
}

func TestGetCampaigns(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	campaigns, err := GetCampaigns("0x0e0c554d2b37105838b45e8b5a49d0edc9b00a8f")
	require.NoError(t, err)
	require.NotEmpty(t, campaigns)
}

func TestDonate(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	private, tr, err := createCampaign(t)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, tr)

	donate, err := Donate(0.1, 1, private, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
	require.NoError(t, err)
	require.NotEmpty(t, donate)
}

func TestGetDonations(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	donations, err := GetDonations(1)
	t.Log(donations)
	require.NoError(t, err)
	require.NotEmpty(t, donations)
}

func TestGetDonorsAddressesAndAmounts(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	donors, amounts, total, err := GetDonorsAddressesAndAmounts(0)
	require.NoError(t, err)
	require.NotEmpty(t, donors)
	require.NotEmpty(t, amounts)
	require.NotEmpty(t, total)
}

func TestGetCampaignTypes(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	types, err := GetCampaignTypes("0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
	require.NoError(t, err)
	require.NotEmpty(t, types)
}

func TestGetCampaignsByOwner(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	campaigns, err := GetCampaignsByOwner("0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
	require.NoError(t, err)
	require.NotEmpty(t, campaigns)
}

func TestPayOut(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	password := "password"

	private, public, err := DecryptPrivateKey("UTC--2023-06-14T06-29-35.797400000Z--a487ff39ac2de30c0105b60dc3e51e377ae95985", password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	payout, err := PayOut(1, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985", private)
	require.NoError(t, err)
	require.NotEmpty(t, payout)
}

func TestSendBackDonations(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	password := "password"

	private, public, err := DecryptPrivateKey("UTC--2023-06-14T06-29-35.797400000Z--a487ff39ac2de30c0105b60dc3e51e377ae95985", password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	sendback, err := SendBackDonations(1, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985", private)
	require.NoError(t, err)
	require.NotEmpty(t, sendback)
}

func TestCreateCategory(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	password := "password"
	private, public, err := DecryptPrivateKey("UTC--2023-06-14T06-29-35.797400000Z--a487ff39ac2de30c0105b60dc3e51e377ae95985", password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	category, err := CreateCategories("Education", "Donate to Sponsor a child Education", "", private, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
	require.NoError(t, err)
	require.NotEmpty(t, category)
}

func TestGetCategories(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	categories, err := GetCategories()

	t.Log(err)
	// t.Log(categories)
	require.NoError(t, err)
	require.NotEmpty(t, categories)
}

func TestSearchCampaignByName(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	campaigns, err := SearchCampaigns("Test Campaign")
	require.NoError(t, err)
	require.NotEmpty(t, campaigns)
}
