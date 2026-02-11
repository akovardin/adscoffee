package ads

const ConversionKeyTemplate = "conversions:%s"

type TrackerInfo struct {
	Timestamp    int64   `json:"timestamp"`
	Action       string  `json:"action"`
	RequestID    string  `json:"request_id"`
	ClickID      string  `json:"click_id"`
	BannerID     string  `json:"banner_id"`
	GroupID      string  `json:"group_id"`
	CampaignID   string  `json:"campaign_id"`
	AdvertiserID string  `json:"advertiser_id"`
	GAID         string  `json:"gaid"`
	OAID         string  `json:"oaid"`
	Bundle       string  `json:"bundle"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	Region       string  `json:"region"`
	Price        float64 `json:"price"`
	Network      string  `json:"network"`
	Size         string  `json:"size"`
	Make         string  `json:"make"`
}
