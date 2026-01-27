package banners

import "time"

type Row struct {
	ID     string
	Title  string
	Price  int
	Active bool

	BannerTargeting     string `gorm:"column:banner_targeting"`
	GroupTargeting      string `gorm:"column:bgroup_targeting"`
	CampaignTargeting   string `gorm:"column:campaign_targeting"`
	AdvertiserTargeting string `gorm:"column:advertiser_targeting"`

	BannerTimetable     string `gorm:"column:banner_timetable"`
	GroupTimetable      string `gorm:"column:bgroup_timetable"`
	CampaignTimetable   string `gorm:"column:campaign_timetable"`
	AdvertiserTimetable string `gorm:"column:advertiser_timetable"`

	BannerBudget     string `gorm:"column:banner_budget"`
	GroupBudget      string `gorm:"column:bgroup_budget"`
	CampaignBudget   string `gorm:"column:ampaign_budget"`
	AdvertiserBudget string `gorm:"column:advertiser_budget"`

	BannerCapping     string `gorm:"column:banner_capping"`
	GroupCapping      string `gorm:"column:bgroup_capping"`
	CampaignCapping   string `gorm:"column:campaign_capping"`
	AdvertiserCapping string `gorm:"column:advertiser_capping"`

	Image        string
	Icon         string
	Clicktracker string
	Imptracker   string
	Target       string

	Label       string
	Description string
	Bundle      string

	Erid string

	GroupID      string `gorm:"column:bgroup_id"`
	CampaignID   string `gorm:"column:campaign_id"`
	AdvertiserID string `gorm:"column:advertiser_id"`

	BannerStart time.Time `gorm:"column:banner_start"`
	BannerEnd   time.Time `gorm:"column:banner_end"`

	GroupStart time.Time `gorm:"column:bgroup_start"`
	GroupEnd   time.Time `gorm:"column:bgroup_end"`

	CampaignStart time.Time `gorm:"column:campaign_start"`
	CampaignEnd   time.Time `gorm:"column:campaign_end"`

	AdvertiserStart time.Time `gorm:"column:advertiser_start"`
	AdvertiserEnd   time.Time `gorm:"column:advertiser_end"`
}
