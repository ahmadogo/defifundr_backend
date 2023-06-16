package db

import (
	"context"
	"testing"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCampaignType(t *testing.T) Campaigns {

	campaignTypes := utils.RandomString(6)

	campaignType, err := testQueries.CreateCampaignType(context.Background(), campaignTypes)
	require.NoError(t, err)
	require.NotEmpty(t, campaignType)

	require.NotZero(t, campaignType.CampaignName)
	return campaignType
}

func TestCreateCampaignType(t *testing.T) {
	createRandomCampaignType(t)
}

func TestGetAllCampaignType(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCampaignType(t)
	}

	campaignTypes, err := testQueries.GetAllCampaignType(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, campaignTypes)
}
