package ads

import (
	"strconv"
	"strings"
	"time"

	"github.com/qor5/admin/v3/media/media_library"
)

const (
	CreativeTypeBanner   = "banner"
	CreativeTypeVideo    = "video"
	CreativeTypeNative   = "native"
	CreativeTypeMediator = "mediator"
)

type Banner struct {
	ID uint

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

	Image        media_library.MediaBox
	Icon         media_library.MediaBox
	Clicktracker string
	Imptracker   string
	Target       string

	Label       string
	Description string
	Bundle      string

	Erid string

	GroupID      string
	CampaignID   string
	AdvertiserID string

	BannerStart time.Time
	BannerEnd   time.Time

	GroupStart time.Time
	GroupEnd   time.Time

	CampaignStart time.Time
	CampaignEnd   time.Time

	AdvertiserStart time.Time
	AdvertiserEnd   time.Time
}

func (b Banner) PriceFormated() string {
	return strconv.FormatFloat(float64(b.Price), 'f', -1, 64)
}

func (b Banner) Media(style string) string {
	if b.Image.Url == "" {
		return ""
	}

	u := b.Image.URL(style)

	if strings.Contains(u, "http:") || strings.Contains(u, "https:") {
		return u // TODO: replace to cdn
	}

	return "https:" + u // TODO: replace to cdn
}
