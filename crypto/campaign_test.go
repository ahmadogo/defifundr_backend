package crypto

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

	password := "password"
	filepath, address, err := GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, filepath)
	require.NotEmpty(t, address)

	private, public, err := DecryptPrivateKey("UTC--2023-06-14T06-29-35.797400000Z--a487ff39ac2de30c0105b60dc3e51e377ae95985", password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	title := "Test Campaign"
	description := "Test Campaign Description"
	image := "Test Campaign Image"
	goal := int(1000000000000000000)
	deadline := time.Now().AddDate(0, 0, 7)
	campaignType := "Test Campaign Type"

	campaign, err := CreateCampaign(title, campaignType, description, goal, deadline, image, private, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
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

	campaign, err := GetCampaign(0, "0x0e0c554d2b37105838b45e8b5a49d0edc9b00a8f")
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

	donate, err := Donate(100000000000, 1, private, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
	require.NoError(t, err)
	require.NotEmpty(t, donate)
}

func TestGetDonations(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	donations, err := GetDonations(1, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
	t.Log(donations)
	require.NoError(t, err)
	require.NotEmpty(t, donations)
}

func TestGetDonorsAddressesAndAmounts(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	donors, amounts, total, err := GetDonorsAddressesAndAmounts(1, "0xa487ff39ac2de30c0105b60dc3e51e377ae95985")
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

func TestGetCampaignByType(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	campaigns, err := GetCampaignByType("Test Campaign Type")
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
