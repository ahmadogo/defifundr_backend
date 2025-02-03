package interfaces

import (
	"encoding/json"
	"fmt"
	"time"
)


type CacheKey string

const (
	AllCampaigns           CacheKey = "campaigns_all"
	LatestActiveCampaigns CacheKey = "campaigns_active"
	CampaignsByCategory  CacheKey = "campaigns_category_%s" 
	CampaignsByOwner     CacheKey = "campaigns_owner_%s"
)

func FormatCacheKey(key CacheKey, value string) string {
	return fmt.Sprintf(string(key), value)
}


type UserResponseInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Avatar   string `json:"avatar"`
}

type Campaigns struct {
	CampaignType       string             `json:"campaign_id"`
	Title              string             `json:"title"`
	Description        string             `json:"description"`
	Goal               float64            `json:"goal"`
	Deadline           time.Time          `json:"deadline"`
	TotalAmountDonated float64            `json:"total_amount_donated"`
	ID                 int                `json:"id"`
	Image              string             `json:"image"`
	Owner              string             `json:"owner"`
	TotalNumber        int64              `json:"total_number"`
	User               []UserResponseInfo `json:"user"`
	Donations          []DonorDetails     `json:"donations"`
}

type SearchCampaignRequest struct {
	Name string `json:"name"`
}

type DonorDetails struct {
	Amount   float64 `json:"amount"`
	Donor    string  `json:"donor"`
	Image    string  `json:"image"`
	Username string  `json:"username"`
}

type Donation struct {
	Amount     string `json:"amount"`
	CampaignId string     `json:"campaign_id"`
}

type Withdraw struct {
	CampaignId int `json:"campaign_id"`
}

func UnmarshalCurrentPrice(data []byte) (CurrentPrice, error) {
	var r CurrentPrice
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CurrentPrice) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CurrentPrice struct {
	Data Data `json:"data"`
}

type Data struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type CampaignCategory struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Id          string `json:"id"`
}
