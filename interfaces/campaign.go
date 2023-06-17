package interfaces

type Campaigns struct {
	CampaignType string `json:"campaign_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Goal         int64  `json:"goal"`
	Deadline     int64  `json:"deadline"`
	Image        string `json:"image"`
}

type Donations struct {
	Amount    int64    `json:"amount"`
	Donations []string `json:"donations"`
	Address   []string `json:"address"`
}


type Donation struct {
	Amount float32 `json:"amount"`
	CampaignId int `json:"campaign_id"`

}