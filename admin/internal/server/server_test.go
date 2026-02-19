package server

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/login"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestServer_Serve(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create mock dependencies
	lb := login.New()
	pb := presets.New()

	// Create a test config with a random available port
	config := Config{
		Port: ":0", // Use port 0 to let the OS choose an available port
	}

	// Create the server instance
	s := New(config, logger, lb, pb)

	// Test the Serve method
	err := s.Serve()
	require.NoError(t, err)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Check that the server is running
	assert.NotNil(t, s.srv)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Test the server response
	s.srv.Handler.ServeHTTP(w, req)
	// Since we don't have any routes registered, we expect a 404
	assert.Equal(t, 404, w.Code)

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create mock dependencies
	lb := login.New()
	pb := presets.New()

	// Create a test config with a random available port
	config := Config{
		Port: ":0",
	}

	// Create the server instance
	s := New(config, logger, lb, pb)

	// Start the server
	err := s.Serve()
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
