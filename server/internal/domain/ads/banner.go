package ads

import (
	"strconv"
	"time"
)

const (
	CreativeTypeBanner   = "banner"
	CreativeTypeVideo    = "video"
	CreativeTypeNative   = "native"
	CreativeTypeMediator = "mediator"
)

type Banner struct {
	ID     string
	Title  string
	Price  int
	Active bool

	Type    string
	Network string

	Targeting Targeting
	Timetable Timetable

	BannerBudget     Budget
	GroupBudget      Budget
	CampaignBudget   Budget
	AdvertiserBudget Budget

	BannerCapping     Capping
	GroupCapping      Capping
	CampaignCapping   Capping
	AdvertiserCapping Capping

	Image        Image
	Icon         Image
	Clicktracker string
	Imptracker   string
	Target       string

	Label       string
	Description string
	Bundle      string

	Erid string

	GroupID      string `gorm:"bgroup_id"`
	CampaignID   string `gorm:"campaign_id"`
	AdvertiserID string `gorm:"advertiser_id"`

	BannerStart time.Time `gorm:"banner_start"`
	BannerEnd   time.Time `gorm:"banner_end"`

	GroupStart time.Time `gorm:"bgroup_start"`
	GroupEnd   time.Time `gorm:"bgroup_end"`

	CampaignStart time.Time `gorm:"campaign_start"`
	CampaignEnd   time.Time `gorm:"campaign_end"`

	AdvertiserStart time.Time `gorm:"advertiser_start"`
	AdvertiserEnd   time.Time `gorm:"advertiser_end"`
}

func (b Banner) PriceFormated() string {
	return strconv.FormatFloat(float64(b.Price), 'f', -1, 64)
}
