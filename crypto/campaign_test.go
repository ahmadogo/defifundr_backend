package crypto

import (
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
)

func createCampaign(t *testing.T) (*bind.TransactOpts, string, error) {
	t.Parallel()
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	password := "pass"
	filepath, address, err := GenerateAccountKeyStone(password)
	require.NoError(t, err)
	require.NotEmpty(t, filepath)
	require.NotEmpty(t, address)

	private, public, err := DecryptPrivateKey("UTC--2023-06-13T20-30-13.183818000Z--0e0c554d2b37105838b45e8b5a49d0edc9b00a8f", password)
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)

	title := "Test Campaign"
	description := "Test Campaign Description"
	image := "Test Campaign Image"
	goal := int(1000000000000000000)
	deadline := int(1000000000000000000)
	campaignType := "Test Campaign Type"

	auth, campaign, err, c := CreateCampaign(title, campaignType, description, goal, deadline, image, private, address)
	if err != nil {
		return nil, "", err
	}
	t.Log(c)

	return auth, campaign, nil
}

func TestCreateCampaign(t *testing.T) {

	a, c, err := createCampaign(t)
	require.NoError(t, err)
	require.NotEmpty(t, c)
	require.NotEmpty(t, a)

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
