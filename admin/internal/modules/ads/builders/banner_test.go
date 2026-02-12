//nolint:errcheck
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

// TestNewBanner tests the NewBanner function
func TestNewBanner(t *testing.T) {
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

	banner := NewBanner(logger, gdb)

	assert.NotNil(t, banner)
	assert.Equal(t, logger, banner.logger)
	assert.Equal(t, gdb, banner.db)
}

// TestBannerConfigure tests the Configure method of Banner
func TestBannerConfigure(t *testing.T) {
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

	banner := NewBanner(logger, gdb)
	b := presets.New()

	// This test is limited because Configure method heavily depends on presets.Builder
	// and other external dependencies that are hard to mock completely
	// In a real scenario, you would use integration tests with a test database
	assert.NotPanics(t, func() {
		banner.Configure(b)
	})

	// Verify that the model was registered
	assert.NotNil(t, b.Model(&models.Banner{}))
}

// TestCopyBanner tests the copyBanner method
func TestCopyBanner(t *testing.T) {
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

	bannerBuilder := &Banner{
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

	// Create a test banner
	now := time.Now()

	// Mock the database calls for finding the original banner
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "label", "description", "active",
		"erid", "ord_category", "ord_targeting", "ord_format", "ord_kktu", "price", "image", "icon",
		"start", "end", "clicktracker", "imptracker", "target", "targeting", "budget", "capping",
		"bgroup_id", "timetable", "archived_at",
	}).AddRow(
		1, now, now, nil, "Test Banner", "Test Label", "Test Description", true,
		"test-erid", "category", "targeting", "format", "kktu", 1000,
		"{}", "{}", now, now.Add(24*time.Hour), "clicktracker", "imptracker", "target",
		"targeting", "budget", "capping", 1, "timetable", nil,
	)

	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE "banners"\."id" = \$1 AND "banners"\."deleted_at" IS NULL ORDER BY "banners"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for creating the copy
	mock.ExpectBegin()

	// Use ExpectQuery instead of ExpectExec for INSERT with RETURNING
	rowsInsert := sqlmock.NewRows([]string{"id"}).AddRow(2)
	mock.ExpectQuery(`INSERT INTO "banners" \("created_at","updated_at","deleted_at","title","label","description","active","erid","ord_category","ord_targeting","ord_format","ord_kktu","price","image","icon","start","end","clicktracker","imptracker","target","targeting","budget","capping","bgroup_id","timetable","archived_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13,\$14,\$15,\$16,\$17,\$18,\$19,\$20,\$21,\$22,\$23,\$24,\$25,\$26\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),      // created_at
			sqlmock.AnyArg(),      // updated_at
			sqlmock.AnyArg(),      // deleted_at
			"Test Banner (Копия)", // title
			"Test Label",          // label
			"Test Description",    // description
			false,                 // active
			"test-erid",           // erid
			"category",            // ord_category
			"targeting",           // ord_targeting
			"format",              // ord_format
			"kktu",                // ord_kktu
			1000,                  // price
			sqlmock.AnyArg(),      // image
			sqlmock.AnyArg(),      // icon
			sqlmock.AnyArg(),      // start
			sqlmock.AnyArg(),      // end
			"clicktracker",        // clicktracker
			"imptracker",          // imptracker
			"target",              // target
			"targeting",           // targeting
			"budget",              // budget
			"capping",             // capping
			1,                     // bgroup_id
			"timetable",           // timetable
			sqlmock.AnyArg(),      // archived_at
		).
		WillReturnRows(rowsInsert)

	mock.ExpectCommit()

	// Call the method
	response, err := bannerBuilder.copyBanner(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestArchiveBanner tests the archiveBanner method
func TestArchiveBanner(t *testing.T) {
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

	bannerBuilder := &Banner{
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

	// Create a test banner
	now := time.Now()

	// Mock the database calls for finding the original banner
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "label", "description", "active",
		"erid", "ord_category", "ord_targeting", "ord_format", "ord_kktu", "price", "image", "icon",
		"start", "end", "clicktracker", "imptracker", "target", "targeting", "budget", "capping",
		"bgroup_id", "timetable", "archived_at",
	}).AddRow(
		1, now, now, nil, "Test Banner", "Test Label", "Test Description", true,
		"test-erid", "category", "targeting", "format", "kktu", 1000,
		"{}", "{}", now, now.Add(24*time.Hour), "clicktracker", "imptracker", "target",
		"targeting", "budget", "capping", 1, "timetable", nil,
	)

	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE "banners"\."id" = \$1 AND "banners"\."deleted_at" IS NULL ORDER BY "banners"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for updating the banner
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "banners" SET "id"=\$1,"created_at"=\$2,"updated_at"=\$3,"deleted_at"=\$4,"title"=\$5,"label"=\$6,"description"=\$7,"active"=\$8,"erid"=\$9,"ord_category"=\$10,"ord_targeting"=\$11,"ord_format"=\$12,"ord_kktu"=\$13,"price"=\$14,"image"=\$15,"icon"=\$16,"start"=\$17,"end"=\$18,"clicktracker"=\$19,"imptracker"=\$20,"target"=\$21,"targeting"=\$22,"budget"=\$23,"capping"=\$24,"bgroup_id"=\$25,"timetable"=\$26,"archived_at"=\$27 WHERE "banners"\."deleted_at" IS NULL AND "id" = \$28`).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // label
			sqlmock.AnyArg(), // description
			sqlmock.AnyArg(), // active
			sqlmock.AnyArg(), // erid
			sqlmock.AnyArg(), // ord_category
			sqlmock.AnyArg(), // ord_targeting
			sqlmock.AnyArg(), // ord_format
			sqlmock.AnyArg(), // ord_kktu
			sqlmock.AnyArg(), // price
			sqlmock.AnyArg(), // image
			sqlmock.AnyArg(), // icon
			sqlmock.AnyArg(), // start
			sqlmock.AnyArg(), // end
			sqlmock.AnyArg(), // clicktracker
			sqlmock.AnyArg(), // imptracker
			sqlmock.AnyArg(), // target
			sqlmock.AnyArg(), // targeting
			sqlmock.AnyArg(), // budget
			sqlmock.AnyArg(), // capping
			sqlmock.AnyArg(), // bgroup_id
			sqlmock.AnyArg(), // timetable
			sqlmock.AnyArg(), // archived_at (this should be the new time value)
			int64(1),         // WHERE id = 1
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	response, err := bannerBuilder.archiveBanner(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUnarchiveBanner tests the unarchiveBanner method
func TestUnarchiveBanner(t *testing.T) {
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

	bannerBuilder := &Banner{
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

	// Create a test banner (that is already archived)
	now := time.Now()
	archivedTime := now.Add(-24 * time.Hour)

	// Mock the database calls for finding the original banner
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "label", "description", "active",
		"erid", "ord_category", "ord_targeting", "ord_format", "ord_kktu", "price", "image", "icon",
		"start", "end", "clicktracker", "imptracker", "target", "targeting", "budget", "capping",
		"bgroup_id", "timetable", "archived_at",
	}).AddRow(
		1, now, now, nil, "Test Banner", "Test Label", "Test Description", true,
		"test-erid", "category", "targeting", "format", "kktu", 1000,
		"{}", "{}", now, now.Add(24*time.Hour), "clicktracker", "imptracker", "target",
		"targeting", "budget", "capping", 1, "timetable", &archivedTime,
	)

	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE "banners"\."id" = \$1 AND "banners"\."deleted_at" IS NULL ORDER BY "banners"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for updating the banner
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "banners" SET "id"=\$1,"created_at"=\$2,"updated_at"=\$3,"deleted_at"=\$4,"title"=\$5,"label"=\$6,"description"=\$7,"active"=\$8,"erid"=\$9,"ord_category"=\$10,"ord_targeting"=\$11,"ord_format"=\$12,"ord_kktu"=\$13,"price"=\$14,"image"=\$15,"icon"=\$16,"start"=\$17,"end"=\$18,"clicktracker"=\$19,"imptracker"=\$20,"target"=\$21,"targeting"=\$22,"budget"=\$23,"capping"=\$24,"bgroup_id"=\$25,"timetable"=\$26,"archived_at"=\$27 WHERE "banners"\."deleted_at" IS NULL AND "id" = \$28`).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // label
			sqlmock.AnyArg(), // description
			sqlmock.AnyArg(), // active
			sqlmock.AnyArg(), // erid
			sqlmock.AnyArg(), // ord_category
			sqlmock.AnyArg(), // ord_targeting
			sqlmock.AnyArg(), // ord_format
			sqlmock.AnyArg(), // ord_kktu
			sqlmock.AnyArg(), // price
			sqlmock.AnyArg(), // image
			sqlmock.AnyArg(), // icon
			sqlmock.AnyArg(), // start
			sqlmock.AnyArg(), // end
			sqlmock.AnyArg(), // clicktracker
			sqlmock.AnyArg(), // imptracker
			sqlmock.AnyArg(), // target
			sqlmock.AnyArg(), // targeting
			sqlmock.AnyArg(), // budget
			sqlmock.AnyArg(), // capping
			sqlmock.AnyArg(), // bgroup_id
			sqlmock.AnyArg(), // timetable
			sqlmock.AnyArg(), // archived_at (this should be NULL now)
			int64(1),         // WHERE id = 1
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	response, err := bannerBuilder.unarchiveBanner(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
