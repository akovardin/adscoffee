package server

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	plogin "github.com/qor5/admin/v3/login"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestServer_Serve(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create mock dependencies
	pb := presets.New()
	lb := plogin.New(pb).
		DB(db).
		Secret("test-secret").
		UserModel(&TestUser{})

	// Create a test config with a random available port
	config := Config{
		Port: ":0", // Use port 0 to let the OS choose an available port
	}

	// Create the server instance
	s := New(config, logger, lb, pb)

	// Test the Serve method
	err = s.Serve()
	require.NoError(t, err)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Check that the server is running
	assert.NotNil(t, s.srv)
	assert.NotNil(t, s.srv.Handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	// Test the server response
	s.srv.Handler.ServeHTTP(w, req)
	// We expect a redirect to login page
	assert.Equal(t, 302, w.Code)

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create mock dependencies
	pb := presets.New()
	lb := plogin.New(pb).
		DB(db).
		Secret("test-secret").
		UserModel(&TestUser{})

	// Create a test config with a random available port
	config := Config{
		Port: ":0",
	}

	// Create the server instance
	s := New(config, logger, lb, pb)

	// Start the server
	err = s.Serve()
	require.NoError(t, err)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test the Shutdown method
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	assert.NoError(t, err)
}

// TestDB is a simple data operator for testing
type TestDB struct{}

func NewTestDB() *TestDB {
	return &TestDB{}
}

func (db *TestDB) Fetch(obj interface{}, id string, ctx *web.EventContext) (interface{}, error) {
	return nil, nil
}

func (db *TestDB) Save(obj interface{}, id string, ctx *web.EventContext) error {
	return nil
}

func (db *TestDB) Delete(obj interface{}, id string, ctx *web.EventContext) error {
	return nil
}

func (db *TestDB) Search(ctx *web.EventContext, params *presets.SearchParams) (*presets.SearchResult, error) {
	return &presets.SearchResult{}, nil
}

// TestUser is a simple user model for testing
type TestUser struct {
	ID       uint
	Account  string
	Password string
}

func (u *TestUser) GetID() string {
	return "1"
}

func (u *TestUser) GetAccountName() string {
	return u.Account
}

func (u *TestUser) GetPassword() string {
	return u.Password
}

func (u *TestUser) GetName() string {
	return u.Account
}

func (u *TestUser) GetAvatar() string {
	return ""
}
