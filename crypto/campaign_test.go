package crypto

import (
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func createCampaign() (string, error) {
	configs, err := utils.LoadConfig("./../")
	if err != nil {
		return "", err
	}
	title := "Test Campaign"
	description := "Test Campaign Description"
	image := "Test Campaign Image"
	goal := int(1000000000000000000)
	deadline := int(1000000000000000000)
	campaignType := "Test Campaign Type"

	campaign, err := CreateCampaign(title, campaignType, description, goal, deadline, image, configs.DeployKey, configs.DeployAddress)
	if err != nil {
		return "", err
	}

	return campaign, nil
}

func TestCreateCampaign(t *testing.T) {
	c, err := createCampaign()
	require.NoError(t, err)
	require.NotEmpty(t, c)
}

func TestGetCampaign(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	campaign, err := GetCampaign(0, configs.DeployAddress)
	require.NoError(t, err)
	require.NotEmpty(t, campaign)
}

func TestGetCampaigns(t *testing.T) {
	configs, err := utils.LoadConfig("./../")
	require.NoError(t, err)
	campaigns, err := GetCampaigns(configs.DeployAddress)
	require.NoError(t, err)
	require.NotEmpty(t, campaigns)
}
