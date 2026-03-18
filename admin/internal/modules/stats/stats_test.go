//nolint:errcheck
package stats

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	query := &Query{}
	stats := New(gdb, query)

	assert.NotNil(t, stats)
	assert.Equal(t, gdb, stats.db)
	assert.Equal(t, query, stats.query)
	assert.NotNil(t, stats.template)
}

// TestConfigure tests the Configure method
func TestConfigure(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	query := &Query{}
	stats := New(gdb, query)
	b := presets.New()

	// Ensure Configure does not panic
	assert.NotPanics(t, func() {
		stats.Configure(b)
	})

	// Verify that the model was registered
	assert.NotNil(t, b.Model(&Dashboard{}))
}

// TestParse tests the parse function
func TestParse(t *testing.T) {
	ctx := &web.EventContext{
		R: &http.Request{
			Form: map[string][]string{
				"key": {"value1,value2,value3"},
			},
		},
	}

	result := parse(ctx, "key")
	assert.Equal(t, []string{"value1", "value2", "value3"}, result)

	// Test empty value
	ctx.R.Form = map[string][]string{}
	result = parse(ctx, "key")
	assert.Nil(t, result)

	// Test single value
	ctx.R.Form = map[string][]string{"key": {"single"}}
	result = parse(ctx, "key")
	assert.Equal(t, []string{"single"}, result)

	// Test with spaces
	ctx.R.Form = map[string][]string{"key": {"a, b , c"}}
	result = parse(ctx, "key")
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

// TestBanners tests the banners method
func TestBanners(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	// Mock the database query
	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "Banner 1").
		AddRow(2, "Banner 2")
	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE "banners"\."deleted_at" IS NULL`).
		WillReturnRows(rows)

	stats := &Stats{db: gdb}
	options := stats.banners()

	assert.Len(t, options, 2)
	assert.Equal(t, "1", options[0].ID)
	assert.Equal(t, "Banner 1", options[0].Name)
	assert.Equal(t, "2", options[1].ID)
	assert.Equal(t, "Banner 2", options[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGroups tests the groups method
func TestGroups(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(10, "Group A").
		AddRow(20, "Group B")
	mock.ExpectQuery(`SELECT \* FROM "bgroups" WHERE "bgroups"\."deleted_at" IS NULL`).
		WillReturnRows(rows)

	stats := &Stats{db: gdb}
	options := stats.groups()

	assert.Len(t, options, 2)
	assert.Equal(t, "10", options[0].ID)
	assert.Equal(t, "Group A", options[0].Name)
	assert.Equal(t, "20", options[1].ID)
	assert.Equal(t, "Group B", options[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCampaigns tests the campaigns method
func TestCampaigns(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(100, "Campaign X").
		AddRow(200, "Campaign Y")
	mock.ExpectQuery(`SELECT \* FROM "campaigns" WHERE "campaigns"\."deleted_at" IS NULL`).
		WillReturnRows(rows)

	stats := &Stats{db: gdb}
	options := stats.campaigns()

	assert.Len(t, options, 2)
	assert.Equal(t, "100", options[0].ID)
	assert.Equal(t, "Campaign X", options[0].Name)
	assert.Equal(t, "200", options[1].ID)
	assert.Equal(t, "Campaign Y", options[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestAdvertisers tests the advertisers method
func TestAdvertisers(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(500, "Advertiser One").
		AddRow(600, "Advertiser Two")
	mock.ExpectQuery(`SELECT \* FROM "advertisers" WHERE "advertisers"\."deleted_at" IS NULL`).
		WillReturnRows(rows)

	stats := &Stats{db: gdb}
	options := stats.advertisers()

	assert.Len(t, options, 2)
	assert.Equal(t, "500", options[0].ID)
	assert.Equal(t, "Advertiser One", options[0].Name)
	assert.Equal(t, "600", options[1].ID)
	assert.Equal(t, "Advertiser Two", options[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestMetrics tests the metrics method
func TestMetrics(t *testing.T) {
	stats := &Stats{}
	options := stats.metrics()

	expected := []Option{
		{ID: MetricRequests, Name: "Запросы"},
		{ID: MetricResponses, Name: "Респонсы"},
		{ID: MetricWins, Name: "Победы"},
		{ID: MetricImpressions, Name: "Показы"},
		{ID: MetricClicks, Name: "Клики"},
		{ID: MetricConversions, Name: "Конверсии"},
		{ID: MetricPrice, Name: "Деньги"},
	}

	assert.Equal(t, expected, options)
}

// TestGrouped tests the grouped method
func TestGrouped(t *testing.T) {
	stats := &Stats{}
	options := stats.grouped()

	expected := []Option{
		{ID: GroupBanner, Name: "Баннеры"},
		{ID: GroupGroup, Name: "Группы"},
		{ID: GroupCampaign, Name: "Кампании"},
		{ID: GroupAdvertiser, Name: "Рекламодатели"},
		{ID: GroupNetwork, Name: "Сети"},
		{ID: GroupBundle, Name: "Бандлы"},
	}

	assert.Equal(t, expected, options)
}

// TestQueryBundles tests the bundles method of Query (via stats.query)
func TestQueryBundles(t *testing.T) {
	// This would require mocking clickhouse, which is complex.
	// We'll skip for now, but could be added later.
	t.Skip("Skipping test that requires clickhouse mock")
}

// TestQueryNetworks tests the networks method of Query
func TestQueryNetworks(t *testing.T) {
	t.Skip("Skipping test that requires clickhouse mock")
}
