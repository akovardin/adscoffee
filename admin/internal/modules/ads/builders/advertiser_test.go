package builders

import (
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

// TestNewAdvertiser tests the NewAdvertiser function
func TestNewAdvertiser(t *testing.T) {
	logger := zaptest.NewLogger(t)
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

	advertiser := NewAdvertiser(logger, gdb)

	assert.NotNil(t, advertiser)
	assert.Equal(t, logger, advertiser.logger)
	assert.Equal(t, gdb, advertiser.db)
}

// TestAdvertiserConfigure tests the Configure method of Advertiser
func TestAdvertiserConfigure(t *testing.T) {
	logger := zaptest.NewLogger(t)
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

	advertiser := NewAdvertiser(logger, gdb)
	b := presets.New()

	// This test is limited because Configure method heavily depends on presets.Builder
	// and other external dependencies that are hard to mock completely
	// In a real scenario, you would use integration tests with a test database
	assert.NotPanics(t, func() {
		advertiser.Configure(b)
	})

	// Verify that the model was registered
	assert.NotNil(t, b.Model(&models.Advertiser{}))
}

// TestCopyAdvertiser tests the copyAdvertiser method
func TestCopyAdvertiser(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a GORM database instance
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	advertiser := &Advertiser{
		logger: logger,
		db:     gdb,
	}

	// Create a mock context
	ctx := &web.EventContext{
		R: &http.Request{
			Form: map[string][]string{
				"id": {"1"}, // This will be converted to string when FormValue is called
			},
		},
	}

	// Create a test advertiser
	now := time.Now()

	// Mock the database calls for finding the original advertiser
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "info", "active", "start", "end", "targeting", "budget", "capping", "timetable", "ord_contract", "archived_at"}).
		AddRow(1, now, now, nil, "Test Advertiser", "Test Info", true, now, now.Add(24*time.Hour), "", "", "", "", "", nil)
	mock.ExpectQuery(`SELECT \* FROM "advertisers" WHERE "advertisers"\."id" = \$1 AND "advertisers"\."deleted_at" IS NULL ORDER BY "advertisers"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for creating the copy
	mock.ExpectBegin()

	// Используем ExpectQuery вместо ExpectExec для INSERT с RETURNING
	rowsInsert := sqlmock.NewRows([]string{"id"}).AddRow(2)
	mock.ExpectQuery(`INSERT INTO "advertisers" \("created_at","updated_at","deleted_at","title","info","active","start","end","targeting","budget","capping","timetable","ord_contract","archived_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13,\$14\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),          // created_at
			sqlmock.AnyArg(),          // updated_at
			sqlmock.AnyArg(),          // deleted_at
			"Test Advertiser (Копия)", // title
			"Test Info",               // info
			false,                     // active
			sqlmock.AnyArg(),          // start
			sqlmock.AnyArg(),          // end
			sqlmock.AnyArg(),          // targeting
			sqlmock.AnyArg(),          // budget
			sqlmock.AnyArg(),          // capping
			sqlmock.AnyArg(),          // timetable
			sqlmock.AnyArg(),          // ord_contract
			sqlmock.AnyArg(),          // archived_at
		).
		WillReturnRows(rowsInsert)

	mock.ExpectCommit()

	// Call the method
	response, err := advertiser.copyAdvertiser(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestArchiveAdvertiser(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a GORM database instance
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	advertiser := &Advertiser{
		logger: logger,
		db:     gdb,
	}

	// Create a mock context
	ctx := &web.EventContext{
		R: &http.Request{
			Form: map[string][]string{
				"id": {"1"},
			},
		},
	}

	// Create a test advertiser
	now := time.Now()

	// Mock the database calls
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "info", "active", "start", "end", "targeting", "budget", "capping", "timetable", "ord_contract", "archived_at"}).
		AddRow(1, now, now, nil, "Test Advertiser", "", true, time.Time{}, time.Time{}, "", "", "", "", "", nil)
	mock.ExpectQuery(`SELECT \* FROM "advertisers" WHERE "advertisers"\."id" = \$1 AND "advertisers"\."deleted_at" IS NULL ORDER BY "advertisers"\."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "advertisers" SET .* WHERE "advertisers"\."deleted_at" IS NULL AND "id" = \$16`).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // info
			sqlmock.AnyArg(), // active
			sqlmock.AnyArg(), // start
			sqlmock.AnyArg(), // end
			sqlmock.AnyArg(), // targeting
			sqlmock.AnyArg(), // budget
			sqlmock.AnyArg(), // capping
			sqlmock.AnyArg(), // timetable
			sqlmock.AnyArg(), // ord_contract
			sqlmock.AnyArg(), // archived_at
			int64(1),         // WHERE id = 1
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Этот запрос должен быть ДО коммита, а не после
	mock.ExpectQuery(`SELECT \* FROM "campaigns" WHERE advertiser_id = \$1 AND "campaigns"\."deleted_at" IS NULL`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "active", "bundle", "start", "end", "targeting", "budget", "capping", "timetable", "advertiser_id", "archived_at"}))

	// Call the method
	response, err := advertiser.archiveAdvertiser(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUnarchiveAdvertiser tests the unarchiveAdvertiser method
func TestUnarchiveAdvertiser(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a GORM database instance
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	advertiser := &Advertiser{
		logger: logger,
		db:     gdb,
	}

	// Create a mock context
	ctx := &web.EventContext{
		R: &http.Request{
			Form: map[string][]string{
				"id": {"1"}, // This will be converted to string when FormValue is called
			},
		},
	}

	// Create a test advertiser
	now := time.Now()

	// Mock the database calls
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "info", "active", "start", "end", "targeting", "budget", "capping", "timetable", "ord_contract", "archived_at"}).
		AddRow(1, now, now, nil, "Test Advertiser", "", true, time.Time{}, time.Time{}, "", "", "", "", "", nil)
	mock.ExpectQuery(`SELECT \* FROM "advertisers" WHERE "advertisers"\."id" = \$1 AND "advertisers"\."deleted_at" IS NULL ORDER BY "advertisers"\."id" LIMIT \$2`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "advertisers" SET .* WHERE "advertisers"\."deleted_at" IS NULL AND "id" = \$16`).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // info
			sqlmock.AnyArg(), // active
			sqlmock.AnyArg(), // start
			sqlmock.AnyArg(), // end
			sqlmock.AnyArg(), // targeting
			sqlmock.AnyArg(), // budget
			sqlmock.AnyArg(), // capping
			sqlmock.AnyArg(), // timetable
			sqlmock.AnyArg(), // ord_contract
			sqlmock.AnyArg(), // archived_at
			int64(1),         // WHERE id = 1
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Этот запрос должен быть ДО коммита, а не после
	mock.ExpectQuery(`SELECT \* FROM "campaigns" WHERE advertiser_id = \$1 AND "campaigns"\."deleted_at" IS NULL`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "active", "bundle", "start", "end", "targeting", "budget", "capping", "timetable", "advertiser_id", "archived_at"}))

	// Call the method
	response, err := advertiser.unarchiveAdvertiser(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
