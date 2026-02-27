package formats

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

const TypeNative = "native"

type Native struct {
	base string
}

func NewNative() *Native {
	return &Native{}
}

func (b *Native) Name() string {
	return "native"
}

type NativeResponse struct {
	Description string   `json:"description"`
	Title       string   `json:"information"`
	Image       string   `json:"image"`
	Target      string   `json:"target"`
	Impressions []string `json:"impressions"`
	Clicks      []string `json:"clicks"`
	Data        string   `json:"data,omitempty"`
}

func (f *Native) Copy(cfg map[string]any) plugins.Format {
	base, _ := cfg["base"].(string)

	return &Native{
		base: base,
	}
}

func (f *Native) Render(ctx context.Context, state *plugins.State) (any, error) {
	items := []NativeResponse{}

	for _, b := range state.Winners {
		click, err := f.tracker(b, state, ads.ActionClick)
		if err != nil {
			return nil, err
		}

		clicktrackers := []string{
			click,
		}

		if b.Clicktracker != "" {
			clicktrackers = append(clicktrackers, b.Clicktracker)
		}

		impression, err := f.tracker(b, state, ads.ActionImpression)
		if err != nil {
			return nil, err
		}

		impressiontrackers := []string{
			impression,
		}

		if b.Imptracker != "" {
			impressiontrackers = append(impressiontrackers, b.Imptracker)
		}

		items = append(items, NativeResponse{
			Title:       b.Title,
			Description: b.Description,
			Target:      b.Target,
			Image:       b.Media("image"),
			Data:        b.Data,

			Impressions: impressiontrackers,
			Clicks:      clicktrackers,
		})
	}

	return items, nil
}

func (f *Native) tracker(w ads.Banner, state *plugins.State, action string) (string, error) {
	info := ads.TrackerInfo{
		Action:       action,
		BannerID:     w.ID,
		GroupID:      w.GroupID,
		CampaignID:   w.CampaignID,
		AdvertiserID: w.AdvertiserID,
		ClickID:      state.ClickID,
		RequestID:    state.RequestID,
	}

	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	return f.base + "/tracker/" + base64.URLEncoding.EncodeToString(data) + ".gif", nil
}
