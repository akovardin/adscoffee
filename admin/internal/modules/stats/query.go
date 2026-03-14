package stats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/go-clickhouse/ch"
	"go.uber.org/zap"

	"go.ads.coffee/platform/admin/internal/clickhouse"
)

const (
	Date           = "2006-01-02"
	DateHour       = "2006-01-02 15"
	DateHourMinute = "2006-01-02 15:04"

	ByHour = "hour"
	ByDay  = "day"
)

// достыпные метрики
const (
	MetricRequests    = "requests"
	MetricResponses   = "responses"
	MetricWins        = "wins"
	MetricImpressions = "impressions"
	MetricClicks      = "clicks"
	MetricConversions = "conversions"
	MetricPrice       = "price"
)

// доступные фильтры
const (
	FilterAdvertiserId = "advertiser_id"
	FilterCampaignId   = "campaign_id"
	FilterGroupId      = "group_id"
	FilterBannerId     = "banner_id"
	FilterNetwork      = "network"
	FilterBundle       = "bundle"
)

// доступные группировки
const (
	GroupAdvertiser = "advertiser_id"
	GroupCampaign   = "campaign_id"
	GroupGroup      = "group_id"
	GroupBanner     = "banner_id"
	GroupNetwork    = "network"
	GroupBundle     = "bundle"
	GroupSlot       = "slot"
)

type Query struct {
	logger     *zap.Logger
	clickhouse *clickhouse.Clickhouse
}

func NewQuery(logger *zap.Logger, clickhouse *clickhouse.Clickhouse) *Query {
	return &Query{
		logger:     logger,
		clickhouse: clickhouse,
	}
}

type Stat struct {
	Labels []string
	//        metric name    hour   value
	Datasets map[string]map[string]float64
}

type Condition struct {
	From    time.Time
	To      time.Time
	Metrics []string // из каких табличек нужно доставать данные
	Filters []Filter
	Groups  []string
	By      string
}

type Filter struct {
	Field string
	Value []string
}

func (q *Query) Select(ctx context.Context, condition Condition) (Stat, error) {
	if len(condition.Groups) == 0 {
		return Stat{}, nil
	}

	if len(condition.Metrics) == 0 {
		return Stat{}, nil
	}

	diff := condition.To.Sub(condition.From)
	labels := []string{}
	if condition.By == ByDay {
		cnt := int(diff.Hours()/24) + 1

		for i := 0; i < cnt; i++ {
			labels = append(labels, condition.From.Add(time.Duration(i)*time.Hour*24).Format(Date))
		}

	} else {
		cnt := int(diff.Hours())

		for i := 0; i < cnt; i++ {
			labels = append(labels, condition.From.Add(time.Duration(i)*time.Hour).Format(DateHour))
		}
	}

	datasets := map[string]map[string]float64{}

	for _, metric := range condition.Metrics {
		st, err := q.query(ctx, metric, condition, labels)
		if err != nil {
			q.logger.Error("query", zap.Error(err))

			return Stat{}, err
		}

		for key, item := range st {
			datasets[key] = item
		}
	}

	return Stat{
		Labels:   labels,
		Datasets: datasets,
	}, nil
}

type Item struct {
	Time   time.Time
	Label0 string
	Label1 string
	Label2 string
	Label3 string
	Label4 string
	Value  float64
}

func (i Item) Key(metric string, groups []string) string {
	if i.Label4 != "" {
		return metric + " - " + groups[0] + ":" + i.Label0 + " - " + groups[1] + ":" + i.Label1 + " - " + groups[2] + ":" + i.Label2 + " - " + groups[3] + ":" + i.Label3 + " - " + groups[4] + ":" + i.Label4
	}

	if i.Label3 != "" {
		return metric + " - " + groups[0] + ":" + i.Label0 + " - " + groups[1] + ":" + i.Label1 + " - " + groups[2] + ":" + i.Label2 + " - " + groups[3] + ":" + i.Label3
	}

	if i.Label2 != "" {
		return metric + " - " + groups[0] + ":" + i.Label0 + " - " + groups[1] + ":" + i.Label1 + " - " + groups[2] + ":" + i.Label2
	}

	if i.Label1 != "" {
		return metric + " - " + groups[0] + ":" + i.Label0 + " - " + groups[1] + ":" + i.Label1
	}

	if i.Label0 != "" {
		return metric + " - " + groups[0] + ":" + i.Label0
	}

	return metric

}

func (q *Query) query(ctx context.Context, metric string, condition Condition, hours []string) (map[string]map[string]float64, error) {
	sel := q.clickhouse.DB.NewSelect()

	if condition.By == ByDay {
		sel.ColumnExpr("toDate(timestamp) as time")
	} else {
		sel.ColumnExpr("timestamp as time")
	}

	for i, group := range condition.Groups {
		sel.ColumnExpr(group+" as label?", i)
	}

	// если метрика деньги, то тут берем price
	if metric == MetricPrice {
		sel.ColumnExpr("toFloat64(sum(price)) as value")
	} else {
		sel.ColumnExpr("sum(count) as value")
	}

	// если метрика деньги, то тут берем impressions
	if metric == MetricPrice {
		sel.ModelTableExpr("impressions_hour")
	} else {
		sel.ModelTableExpr(metric + "_hour")
	}

	for _, filter := range condition.Filters {
		if len(filter.Value) == 0 {
			continue
		}

		sel.Where(filter.Field+" IN (?)", ch.In(filter.Value))
	}

	sel.Where("timestamp >= ?", condition.From)
	sel.Where("timestamp <= ?", condition.To)

	sel.Group("timestamp")

	for _, group := range condition.Groups {
		sel.Group(group)
	}

	sel.Order("time DESC")
	sel.Limit(1000)

	fmt.Println(sel.String())

	items := []Item{}

	err := sel.Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	data := map[string]map[string]float64{}

	// раскладываю по метрикам часы с их значениями
	for _, item := range items {
		key := item.Key(metric, condition.Groups)
		hour := ""
		if condition.By == ByDay {
			hour = item.Time.Format(Date)
		} else {
			hour = item.Time.Format(DateHour)
		}

		if _, exist := data[key]; !exist {
			data[key] = map[string]float64{}
		}
		data[key][hour] = item.Value
	}

	if condition.By == ByHour {
		// заполняем дырки в случае отсутствия данных в определенном часе
		for key, item := range data {
			for _, hour := range hours {
				_, exist := item[hour]
				if !exist {
					data[key][hour] = 0
				}
			}
		}
	}

	return data, nil
}

type BundleModel struct {
	Bundle string
}

func (q *Query) bundles(ctx context.Context) ([]Option, error) {
	sel := q.clickhouse.DB.NewSelect()

	sel.ColumnExpr("bundle")
	sel.ModelTableExpr("impressions_hour")
	sel.Group("bundle")

	items := []BundleModel{}

	err := sel.Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	options := make([]Option, 0, len(items))
	for _, item := range items {
		if !strings.Contains(item.Bundle, ".") {
			continue
		}
		options = append(options, Option{
			ID:   item.Bundle,
			Name: item.Bundle,
		})
	}

	return options, err
}

type NetworkModel struct {
	Network string
}

func (q *Query) networks(ctx context.Context) ([]Option, error) {
	sel := q.clickhouse.DB.NewSelect()

	sel.ColumnExpr("network")
	sel.ModelTableExpr("impressions_hour")
	sel.Group("network")

	items := []NetworkModel{}

	err := sel.Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	options := make([]Option, 0, len(items))
	for _, item := range items {
		options = append(options, Option{
			ID:   item.Network,
			Name: item.Network,
		})
	}

	return options, err
}

type SlotModel struct {
	Slot string
}
