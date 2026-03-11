package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli/v3"
	"go.ads.coffee/platform/analytics/internal/clickhouse"
)

type Aggregate struct {
	click *clickhouse.Clickhouse
}

func NewAggregate(click *clickhouse.Clickhouse) *Aggregate {
	return &Aggregate{
		click: click,
	}
}

func (a *Aggregate) Run(ctx context.Context, cmd *cli.Command) error {
	table := cmd.String("table")

	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	startTime := endTime.Add(-1 * time.Hour)

	log.Printf("aggregating data from %s to %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	// SQL запрос для агрегации данных
	query := `
	INSERT INTO analytics.` + table + `_hour
	SELECT 
		action,
		toDateTime(formatDateTime(timestamp, '%Y-%m-%d %H:00:00')) as tt,
		banner_id,
		group_id,
		campaign_id,
		advertiser_id,
		city,
		country,
		region,
		sum( multiIf(price > 0, price/1000, 0 )) as price,
		count(*) as count,
		network,
		bundle
	FROM analytics.` + table + `
	WHERE timestamp >= ? AND timestamp < ?
	GROUP BY 
		action,
		toDateTime(formatDateTime(timestamp, '%Y-%m-%d %H:00:00')),
		banner_id,
		group_id,
		campaign_id,
		advertiser_id,
		city,
		country,
		region,
		network,
		bundle
	`

	log.Printf("query: %s", query)

	_, err := a.click.DB.ExecContext(ctx, query, startTime, endTime)
	if err != nil {
		return fmt.Errorf("failed to execute aggregation query: %v", err)
	}

	return nil
}
