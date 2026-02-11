package analytics

import (
	"context"
	"fmt"
	"time"

	"go.ads.coffee/platform/pkg/kafkapool"
	"go.ads.coffee/platform/pkg/telemetry"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Analytics struct {
	tel      *telemetry.Telemetry
	producer *kafkapool.Producer
}

func New(pool *kafkapool.Pool, tel *telemetry.Telemetry) (*Analytics, error) {
	if err := tel.Register(actions, money); err != nil {
		return nil, err
	}

	kfk, err := pool.GetPool("main")
	if err != nil {
		return nil, fmt.Errorf("kafka pool error: %w", err)
	}

	producer, err := kafkapool.NewProducer(kfk)
	if err != nil {
		return nil, fmt.Errorf("kafka producer error: %w", err)
	}

	return &Analytics{
		tel:      tel,
		producer: producer,
	}, nil
}

func (r *Analytics) Log(ctx context.Context, name string, event ads.Event) error {
	actions.WithLabelValues(
		name,
		event.Network,
	).Inc()

	data, err := event.JSON()
	if err != nil {
		return err
	}
	// write to kafka
	r.producer.SendAsync(
		context.WithoutCancel(ctx), // нужно проверить на отмену по таймауту
		name,
		data,
	)

	return nil
}

func (r *Analytics) LogRequest(ctx context.Context, state *plugins.State) error {
	return r.Log(
		ctx,
		ads.ActionRequest,
		ads.Event{
			RequestID: state.RequestID,
			Timestamp: time.Now().Unix(),
			Action:    ads.ActionRequest,
			GAID:      "",
			OAID:      "",
			// Country:   rc.Request.Country(),
			// Region:    rc.Request.Region(),
			// City:      rc.Request.City(),
			// Network:   rc.Network,
			// Make: rc.Request.Make(),
		},
	)
}

func (r *Analytics) LogResponse(ctx context.Context, w ads.Banner, state *plugins.State) error {
	return r.Log(
		ctx,
		ads.ActionResponse,
		ads.Event{
			RequestID: state.RequestID,
			ClickID:   state.ClickID,
			Timestamp: time.Now().Unix(),
			Action:    ads.ActionResponse,

			BannerID:     w.ID,
			GroupID:      w.GroupID,
			CampaignID:   w.CampaignID,
			AdvertiserID: w.AdvertiserID,

			GAID: "",
			OAID: "",
			// Country: rc.Request.Country(),
			// Region:  rc.Request.Region(),
			// City:    rc.Request.City(),
			// Network:   rc.Network,
			Price: float64(w.Price),
		},
	)
}

func (r *Analytics) LogWin(ctx context.Context, data ads.TrackerInfo) error {
	data.Action = ads.ActionWin
	return r.Log(ctx, ads.ActionWin, ads.Event(data))
}

func (r *Analytics) LogConversion(ctx context.Context, data ads.TrackerInfo) error {
	data.Action = ads.ActionConversion
	return r.Log(ctx, ads.ActionConversion, ads.Event(data))
}

func (r *Analytics) LogImpression(ctx context.Context, data ads.TrackerInfo) error {
	money.WithLabelValues(
		ads.ActionImpression,
		data.Network,
	).Add(data.Price)

	data.Action = ads.ActionImpression
	return r.Log(ctx, ads.ActionImpression, ads.Event(data))
}

func (r *Analytics) LogClick(ctx context.Context, data ads.TrackerInfo) error {
	data.Action = ads.ActionCLick
	return r.Log(ctx, ads.ActionCLick, ads.Event(data))
}
