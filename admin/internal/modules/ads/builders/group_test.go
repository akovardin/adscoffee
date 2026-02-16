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

// TestNewGroup tests the NewGroup function
func TestNewGroup(t *testing.T) {
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

	group := NewGroup(logger, gdb)

	assert.NotNil(t, group)
	assert.Equal(t, logger, group.logger)
	assert.Equal(t, gdb, group.db)
}

// TestGroupConfigure tests the Configure method of Group
func TestGroupConfigure(t *testing.T) {
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

	group := NewGroup(logger, gdb)
	b := presets.New()

	// This test is limited because Configure method heavily depends on presets.Builder
	// and other external dependencies that are hard to mock completely
	// In a real scenario, you would use integration tests with a test database
	assert.NotPanics(t, func() {
		group.Configure(b)
	})

	// Verify that the model was registered
	assert.NotNil(t, b.Model(&models.Bgroup{}))
}

// TestCopyGroup tests the copyGroup method
func TestCopyGroup(t *testing.T) {
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

	groupBuilder := &Group{
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

	// Create a test group
	now := time.Now()

	// Mock the database calls for finding the original group
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "active", "price",
		"start", "end", "targeting", "budget", "capping", "timetable", "campaign_id", "archived_at",
	}).AddRow(
		1, now, now, nil, "Test Group", true, 1000,
		now, now.Add(24*time.Hour), "targeting", "budget", "capping", "timetable", 1, nil,
	)

	mock.ExpectQuery(`SELECT \* FROM "bgroups" WHERE "bgroups"\."id" = \$1 AND "bgroups"\."deleted_at" IS NULL ORDER BY "bgroups"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for creating the copy
	mock.ExpectBegin()

	// Use ExpectQuery instead of ExpectExec for INSERT with RETURNING
	rowsInsert := sqlmock.NewRows([]string{"id"}).AddRow(2)
	mock.ExpectQuery(`INSERT INTO "bgroups" \("created_at","updated_at","deleted_at","title","active","price","start","end","targeting","budget","capping","timetable","campaign_id","archived_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13,\$14\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),     // created_at
			sqlmock.AnyArg(),     // updated_at
			sqlmock.AnyArg(),     // deleted_at
			"Test Group (Копия)", // title
			false,                // active
			1000,                 // price
			sqlmock.AnyArg(),     // start
			sqlmock.AnyArg(),     // end
			"targeting",          // targeting
			"budget",             // budget
			"capping",            // capping
			"timetable",          // timetable
			1,                    // campaign_id
			sqlmock.AnyArg(),     // archived_at
		).
		WillReturnRows(rowsInsert)

	mock.ExpectCommit()

	// Mock the database calls for finding banners
	bannerRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "label", "description", "active",
		"erid", "ord_category", "ord_targeting", "ord_format", "ord_kktu", "price", "image", "icon",
		"start", "end", "clicktracker", "imptracker", "target", "targeting", "budget", "capping",
		"bgroup_id", "timetable", "archived_at",
	})
	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE bgroup_id = \$1`).
		WithArgs(1).
		WillReturnRows(bannerRows)

	// Call the method
	response, err := groupBuilder.copyGroup(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestArchiveGroup tests the archiveGroup method
func TestArchiveGroup(t *testing.T) {
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

	groupBuilder := &Group{
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

	// Create a test group
	now := time.Now()

	// Mock the database calls for finding the original group
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "active", "price",
		"start", "end", "targeting", "budget", "capping", "timetable", "campaign_id", "archived_at",
	}).AddRow(
		1, now, now, nil, "Test Group", true, 1000,
		now, now.Add(24*time.Hour), "targeting", "budget", "capping", "timetable", 1, nil,
	)

	mock.ExpectQuery(`SELECT \* FROM "bgroups" WHERE "bgroups"\."id" = \$1 AND "bgroups"\."deleted_at" IS NULL ORDER BY "bgroups"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for updating the group
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "bgroups" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"title"=\$4,"active"=\$5,"price"=\$6,"start"=\$7,"end"=\$8,"targeting"=\$9,"budget"=\$10,"capping"=\$11,"timetable"=\$12,"campaign_id"=\$13,"archived_at"=\$14 WHERE "bgroups"\."deleted_at" IS NULL AND "id" = \$15`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // active
			sqlmock.AnyArg(), // price
			sqlmock.AnyArg(), // start
			sqlmock.AnyArg(), // end
			sqlmock.AnyArg(), // targeting
			sqlmock.AnyArg(), // budget
			sqlmock.AnyArg(), // capping
			sqlmock.AnyArg(), // timetable
			sqlmock.AnyArg(), // campaign_id
			sqlmock.AnyArg(), // archived_at (this should be the new time value)
			int64(1),         // WHERE id = 1
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Mock the database calls for finding banners
	bannerRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "label", "description", "active",
		"erid", "ord_category", "ord_targeting", "ord_format", "ord_kktu", "price", "image", "icon",
		"start", "end", "clicktracker", "imptracker", "target", "targeting", "budget", "capping",
		"bgroup_id", "timetable", "archived_at",
	})
	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE bgroup_id = \$1`).
		WithArgs(1).
		WillReturnRows(bannerRows)

	// Call the method
	response, err := groupBuilder.archiveGroup(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUnarchiveGroup tests the unarchiveGroup method
func TestUnarchiveGroup(t *testing.T) {
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

	groupBuilder := &Group{
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

	// Create a test group (that is already archived)
	now := time.Now()
	archivedTime := now.Add(-24 * time.Hour)

	// Mock the database calls for finding the original group
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "active", "price",
		"start", "end", "targeting", "budget", "capping", "timetable", "campaign_id", "archived_at",
	}).AddRow(
		1, now, now, nil, "Test Group", true, 1000,
		now, now.Add(24*time.Hour), "targeting", "budget", "capping", "timetable", 1, &archivedTime,
	)

	mock.ExpectQuery(`SELECT \* FROM "bgroups" WHERE "bgroups"\."id" = \$1 AND "bgroups"\."deleted_at" IS NULL ORDER BY "bgroups"\."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	// Mock the database calls for updating the group
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "bgroups" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"title"=\$4,"active"=\$5,"price"=\$6,"start"=\$7,"end"=\$8,"targeting"=\$9,"budget"=\$10,"capping"=\$11,"timetable"=\$12,"campaign_id"=\$13,"archived_at"=\$14 WHERE "bgroups"\."deleted_at" IS NULL AND "id" = \$15`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // active
			sqlmock.AnyArg(), // price
			sqlmock.AnyArg(), // start
			sqlmock.AnyArg(), // end
			sqlmock.AnyArg(), // targeting
			sqlmock.AnyArg(), // budget
			sqlmock.AnyArg(), // capping
			sqlmock.AnyArg(), // timetable
			sqlmock.AnyArg(), // campaign_id
			sqlmock.AnyArg(), // archived_at (this should be NULL now)
			int64(1),         // WHERE id = 1
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Mock the database calls for finding banners
	bannerRows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "title", "label", "description", "active",
		"erid", "ord_category", "ord_targeting", "ord_format", "ord_kktu", "price", "image", "icon",
		"start", "end", "clicktracker", "imptracker", "target", "targeting", "budget", "capping",
		"bgroup_id", "timetable", "archived_at",
	})
	mock.ExpectQuery(`SELECT \* FROM "banners" WHERE bgroup_id = \$1`).
		WithArgs(1).
		WillReturnRows(bannerRows)

	// Call the method
	response, err := groupBuilder.unarchiveGroup(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
