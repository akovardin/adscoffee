package banners

import (
	"context"
	"encoding/json"
	"net"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type Repo struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewRepo(logger *zap.Logger, db *gorm.DB) *Repo {
	return &Repo{
		logger: logger.Named("banners"),
		db:     db,
	}
}

func (b *Repo) All(ctx context.Context) ([]ads.Banner, error) {
	rows := []Row{}

	err := b.db.Model(Row{}).Raw(`select 
    banners.id as id, 
    banners.title as title, 
    CASE 
        WHEN banners.price IS NOT NULL AND banners.price > 0  THEN banners.price
        ELSE bgroups.price
    END AS price,
    banners.active as active, 
    
    banners.targeting as banner_targeting,
    bgroups.targeting as bgroup_targeting,
    campaigns.targeting as campaign_targeting,
    advertisers.targeting as advertiser_targeting,

	banners.timetable as banner_timetable,
    bgroups.timetable as bgroup_timetable,
    campaigns.timetable as campaign_timetable,
    advertisers.timetable as advertiser_timetable,

    banners.budget as banner_budget,
    bgroups.budget as bgroup_budget,
    campaigns.budget as campaign_budget,
    advertisers.budget as advertiser_budget,

    banners.capping as banner_capping,
    bgroups.capping as bgroup_capping,
    campaigns.capping as campaign_capping,
    advertisers.capping as advertiser_capping,

    banners.image as image,
    banners.icon as icon,
	banners.clicktracker as clicktracker,
	banners.imptracker as imptracker,
	banners.target as target,

	banners.label as label,
	banners.description as description,
	campaigns.bundle as bundle,

	banners.erid as erid,
 	
    bgroups.id as bgroup_id,
    campaigns.id as campaign_id,
	advertisers.id as advertiser_id,
    
    banners.start as banner_start,
    banners.end as banner_end,
    bgroups.start as bgroups_start,
    bgroups.end as bgroups_end,
    campaigns.start as campaign_start,
    campaigns.end as campaign_end,
    advertisers.start as advertiser_start,
    advertisers.end as advertiser_end
from banners 
join bgroups ON (banners.bgroup_id = bgroups.id)
join campaigns ON (bgroups.campaign_id = campaigns.id)
join advertisers ON (campaigns.advertiser_id = advertisers.id)
where
    banners.active = true 
    and banners.deleted_at is NULL
    and banners.archived_at is NULL
    and bgroups.active = true 
	and bgroups.deleted_at is NULL
	and bgroups.archived_at is NULL
    and campaigns.active = true 
	and campaigns.deleted_at is NULL
	and campaigns.archived_at is NULL
    and advertisers.active = true
	and advertisers.deleted_at is NULL
	and advertisers.archived_at is NULL
    and (banners."end" is null or banners."end" >  NOW() or banners."end" < '2001-01-02 00:00:00')
    and (campaigns."end" is null or campaigns."end" >  NOW() or campaigns."end" < '2001-01-02 00:00:00')
    and (bgroups."end" is null or bgroups."end" >  NOW() or bgroups."end" < '2001-01-02 00:00:00')
    and (advertisers."end" is null or advertisers."end" >  NOW() or advertisers."end" < '2001-01-02 00:00:00')`).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	banners := make([]ads.Banner, 0, len(rows))

	for _, row := range rows {
		banner, err := toModel(row)
		if err != nil {
			b.logger.Warn("error on convert row to model", zap.Error(err), zap.String("id", row.ID))

			continue
		}
		banners = append(banners, banner)
	}

	return banners, nil
}

func toModel(row Row) (ads.Banner, error) {
	banner := ads.Banner{
		ID:     row.ID,
		Title:  row.Title,
		Price:  row.Price,
		Active: row.Active,

		Clicktracker: row.Clicktracker,
		Imptracker:   row.Imptracker,
		Target:       row.Target,

		Label:       row.Label,
		Description: row.Description,
		Bundle:      row.Bundle,

		Erid: row.Erid,

		GroupID:      row.GroupID,
		CampaignID:   row.CampaignID,
		AdvertiserID: row.AdvertiserID,

		BannerStart: row.BannerStart,
		BannerEnd:   row.BannerEnd,

		GroupStart: row.GroupStart,
		GroupEnd:   row.GroupEnd,

		CampaignStart: row.CampaignStart,
		CampaignEnd:   row.CampaignEnd,

		AdvertiserStart: row.AdvertiserStart,
		AdvertiserEnd:   row.AdvertiserEnd,
	}

	var err error

	// targetings
	atg, err := newTargeting(row.AdvertiserTargeting)
	if err != nil {
		return ads.Banner{}, err
	}

	ctg, err := newTargeting(row.CampaignTargeting)
	if err != nil {
		return ads.Banner{}, err
	}

	gtg, err := newTargeting(row.GroupTargeting)
	if err != nil {
		return ads.Banner{}, err
	}

	btg, err := newTargeting(row.BannerTargeting)
	if err != nil {
		return ads.Banner{}, err
	}

	banner.Targeting = atg.merge(ctg).merge(gtg).merge(btg).toDomain()

	// timetables
	att, err := ads.NewTimetable(row.AdvertiserTimetable)
	if err != nil {
		return ads.Banner{}, err
	}

	ctt, err := ads.NewTimetable(row.CampaignTimetable)
	if err != nil {
		return ads.Banner{}, err
	}

	gtt, err := ads.NewTimetable(row.GroupTimetable)
	if err != nil {
		return ads.Banner{}, err
	}

	btt, err := ads.NewTimetable(row.BannerTimetable)
	if err != nil {
		return ads.Banner{}, err
	}

	banner.Timetable = att.Merge(ctt).Merge(gtt).Merge(btt)

	// budget
	if banner.BannerBudget, err = ads.NewBudget(row.BannerBudget); err != nil {
		return ads.Banner{}, err
	}

	if banner.GroupBudget, err = ads.NewBudget(row.GroupBudget); err != nil {
		return ads.Banner{}, err
	}

	if banner.CampaignBudget, err = ads.NewBudget(row.CampaignBudget); err != nil {
		return ads.Banner{}, err
	}

	if banner.AdvertiserBudget, err = ads.NewBudget(row.AdvertiserBudget); err != nil {
		return ads.Banner{}, err
	}

	// capping
	if banner.BannerCapping, err = ads.NewCapping(row.BannerCapping); err != nil {
		return ads.Banner{}, err
	}

	if banner.GroupCapping, err = ads.NewCapping(row.GroupCapping); err != nil {
		return ads.Banner{}, err
	}

	if banner.CampaignCapping, err = ads.NewCapping(row.CampaignCapping); err != nil {
		return ads.Banner{}, err
	}

	if banner.AdvertiserCapping, err = ads.NewCapping(row.AdvertiserCapping); err != nil {
		return ads.Banner{}, err
	}

	// images
	if banner.Image, err = ads.NewImage(row.Image); err != nil {
		return ads.Banner{}, err
	}

	if banner.Icon, err = ads.NewImage(row.Icon); err != nil {
		return ads.Banner{}, err
	}

	return banner, nil
}

type excludeInclude struct {
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	IncludeOr  []string `json:"include_or"`
	ExcludeOr  []string `json:"exclude_or"`
	IncludeAnd []string `json:"include_and"`
	ExcludeAnd []string `json:"exclude_and"`
}

func (e excludeInclude) merge(source excludeInclude) excludeInclude {
	result := e

	if len(source.IncludeOr) > 0 {
		result.IncludeOr = source.IncludeOr
	}
	if len(source.ExcludeOr) > 0 {
		result.ExcludeOr = source.ExcludeOr
	}
	if len(source.IncludeAnd) > 0 {
		result.IncludeAnd = source.IncludeAnd
	}
	if len(source.ExcludeAnd) > 0 {
		result.ExcludeAnd = source.ExcludeAnd
	}

	return result
}

func (e excludeInclude) toExcludeIncludeString() ads.ExcludeInclude {
	return ads.ExcludeInclude{
		IncludeOr:  e.IncludeOr,
		ExcludeOr:  e.ExcludeOr,
		IncludeAnd: e.IncludeAnd,
		ExcludeAnd: e.ExcludeAnd,
	}
}

func (e excludeInclude) toExcludeIncludeIP() ads.ExcludeIncludeIP {
	t := ads.ExcludeIncludeIP{
		Include: []*net.IPNet{},
		Exclude: []*net.IPNet{},
	}

	for _, v := range e.Include {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			t.Include = append(t.Include, network)
		}
	}

	for _, v := range e.Exclude {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			t.Exclude = append(t.Exclude, network)
		}
	}

	return t
}

type targeting struct {
	Bundle   excludeInclude `json:"bundle"`
	Audience excludeInclude `json:"audience"`
	Bapp     excludeInclude `json:"bapp"`
	IP       excludeInclude `json:"ip"`
	Country  excludeInclude `json:"country"`
	City     excludeInclude `json:"city"`
	Region   excludeInclude `json:"region"`
	Network  excludeInclude `json:"network"`
}

func newTargeting(data string) (targeting, error) {
	t := targeting{}

	if data == "" {
		return t, nil
	}

	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return targeting{}, err
	}

	return t, nil
}

func (t targeting) merge(source targeting) targeting {
	return targeting{
		Bundle:   t.Bundle.merge(source.Bundle),
		Audience: t.Audience.merge(source.Audience),
		Bapp:     t.Bapp.merge(source.Bapp),
		IP:       t.IP.merge(source.IP),
		Country:  t.Country.merge(source.Country),
		City:     t.City.merge(source.City),
		Region:   t.Region.merge(source.Region),
		Network:  t.Network.merge(source.Network),
	}
}

func (t targeting) toDomain() ads.Targeting {
	return ads.Targeting{
		Bundle:   t.Bundle.toExcludeIncludeString(),
		Audience: t.Audience.toExcludeIncludeString(),
		Bapp:     t.Bapp.toExcludeIncludeString(),
		IP:       t.Bapp.toExcludeIncludeIP(),
		Country:  t.Country.toExcludeIncludeString(),
		City:     t.City.toExcludeIncludeString(),
		Region:   t.Region.toExcludeIncludeString(),
		Network:  t.Network.toExcludeIncludeString(),
	}
}
