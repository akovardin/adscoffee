package static

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

// Mock implementations for testing
type MockSession struct {
	mock.Mock
}

func (m *MockSession) Start(r *http.Request, value string) error {
	args := m.Called(r, value)
	return args.Error(0)
}

type MockAnalytics struct {
	mock.Mock
}

func (m *MockAnalytics) LogImpression(ctx context.Context, data ads.TrackerInfo) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type MockBanner struct {
	mock.Mock
}

func (m *MockBanner) Banner(ctx context.Context, base string, banner ads.Banner, w http.ResponseWriter) error {
	args := m.Called(ctx, base, banner, w)
	return args.Error(0)
}

func TestStatic_Name(t *testing.T) {
	static := &Static{}
	name := static.Name()
	assert.Equal(t, "outputs.static", name)
}

func TestStatic_Copy(t *testing.T) {
	original := &Static{
		base: "http://example.com",
	}

	cfg := map[string]any{
		"base": "http://newbase.com",
	}

	copied := original.Copy(cfg)
	staticCopied, ok := copied.(*Static)
	assert.True(t, ok)
	assert.Equal(t, "http://newbase.com", staticCopied.base)
}

func TestStatic_Do_NonImgAction(t *testing.T) {
	static := &Static{}

	// Create a request with a non-img action context
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "action", "other")
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	err := static.Do(context.Background(), state)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestStatic_Do_EmptyWinners(t *testing.T) {
	static := &Static{}

	// Create a request with img action context
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "action", "img")
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
		Winners:  []ads.Banner{}, // Empty winners
	}

	err := static.Do(context.Background(), state)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestStatic_Do_SessionStartError(t *testing.T) {
	// Create mocks
	mockSession := new(MockSession)
	mockAnalytics := new(MockAnalytics)
	mockBanner := new(MockBanner)

	static := &Static{
		sessions:  mockSession,
		analytics: mockAnalytics,
		format:    mockBanner,
		base:      "http://example.com",
	}

	// Create a request with img action context
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "action", "img")
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a test banner
	banner := ads.Banner{
		ID: "test-banner-id",
	}

	state := &plugins.State{
		Request:  req,
		Response: rr,
		Winners:  []ads.Banner{banner},
	}

	// Set up mock expectations
	mockSession.On("Start", req, "test-banner-id").Return(fmt.Errorf("session error"))

	err := static.Do(context.Background(), state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error on start session")

	// Assert that the mock expectations were met
	mockSession.AssertExpectations(t)
}

func TestStatic_Do_Success(t *testing.T) {
	// Create mocks
	mockSession := new(MockSession)
	mockAnalytics := new(MockAnalytics)
	mockBanner := new(MockBanner)

	static := &Static{
		base:      "http://example.com",
		sessions:  mockSession,
		analytics: mockAnalytics,
		format:    mockBanner,
	}

	// Create a request with img action context
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "action", "img")
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a test banner
	banner := ads.Banner{
		ID: "test-banner-id",
	}

	state := &plugins.State{
		Request:  req,
		Response: rr,
		Winners:  []ads.Banner{banner},
	}

	// Set up mock expectations
	mockSession.On("Start", req, "test-banner-id").Return(nil)
	mockAnalytics.On("LogImpression", mock.Anything, ads.TrackerInfo{}).Return(nil)
	mockBanner.On("Banner", mock.Anything, "http://example.com", banner, rr).Return(nil)

	err := static.Do(context.Background(), state)
	assert.NoError(t, err)

	// Assert that the mock expectations were met
	mockSession.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
	mockBanner.AssertExpectations(t)
}
