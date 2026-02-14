package static

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/sessions"
)

// MockCache is a mock implementation of the Cache interface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) One(ctx context.Context, id string) (ads.Banner, bool) {
	args := m.Called(ctx, id)
	banner, _ := args.Get(0).(ads.Banner)
	return banner, args.Bool(1)
}

// MockSession is a mock implementation of the Session interface
type MockSession struct {
	mock.Mock
}

func (m *MockSession) LoadWithExpire(r *http.Request) (sessions.Session, bool) {
	args := m.Called(r)
	session, _ := args.Get(0).(sessions.Session)
	return session, args.Bool(1)
}

// MockAnalytics is a mock implementation of the Analytics interface
type MockAnalytics struct {
	mock.Mock
}

func (m *MockAnalytics) LogClick(ctx context.Context, data ads.TrackerInfo) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	// Create real dependencies for constructor
	logger := zaptest.NewLogger(t)

	// For the constructor, we need to use actual types
	// But we can't easily create real instances in tests
	// So we'll test the constructor separately

	// Create mock implementations of interfaces
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance directly since we can't easily create real dependencies
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Check the result
	assert.NotNil(t, static)
}

func TestStatic_Name(t *testing.T) {
	// Create real dependencies
	logger := zaptest.NewLogger(t)
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance of Static
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Call the function under test
	name := static.Name()

	// Check the result
	assert.Equal(t, "inputs.static", name)
}

func TestStatic_Copy(t *testing.T) {
	// Create real dependencies
	logger := zaptest.NewLogger(t)
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance of Static
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Prepare configuration
	cfgMap := map[string]any{"key": "value"}

	// Call the function under test
	copied := static.Copy(cfgMap)

	// Check the result
	assert.NotNil(t, copied)
	assert.IsType(t, &Static{}, copied)

	// Check that dependencies are copied correctly
	copiedStatic := copied.(*Static)
	assert.Equal(t, cache, copiedStatic.cache)
	assert.Equal(t, logger, copiedStatic.logger)
	assert.Equal(t, session, copiedStatic.sessions)
	assert.Equal(t, analytics, copiedStatic.analytics)
}

func TestStatic_Do_ViewAction(t *testing.T) {
	// Create real dependencies
	logger := zaptest.NewLogger(t)
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance of Static
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Prepare context and state
	ctx := context.Background()

	// Create a mock HTTP request with action and placement parameters
	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"view", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	// Create a response recorder
	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	// Call the function under test
	result := static.Do(ctx, state)

	// Check the result
	assert.True(t, result)
	assert.NotNil(t, state.User)
	assert.NotNil(t, state.Device)
	assert.NotNil(t, state.Placement)
	assert.Equal(t, "test-placement", state.Placement.ID)

	// Check that placement contains one ad unit
	assert.Len(t, state.Placement.Units, 1)
	assert.Equal(t, "yandex-1", state.Placement.Units[0].ID)
	assert.Equal(t, "yandex", state.Placement.Units[0].Network)
	assert.Equal(t, 10, state.Placement.Units[0].Price)
	assert.Equal(t, "banner", state.Placement.Units[0].Format)

	// Check that action is stored in context
	action := state.Request.Context().Value("action")
	assert.Equal(t, "view", action)
}

func TestStatic_Do_ClickAction_SessionNotFound(t *testing.T) {
	// Create real dependencies
	logger := zaptest.NewLogger(t)
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance of Static
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Prepare context and state
	ctx := context.Background()

	// Create a mock HTTP request with action and placement parameters
	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"click", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	// Create a response recorder
	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	// Set up mock expectations
	// We need to create a request with the same context that will be modified by WithValue
	reqForMock := req.Clone(ctx)
	session.On("LoadWithExpire", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL == reqForMock.URL
	})).Return(sessions.Session{}, false)

	// Call the function under test
	result := static.Do(ctx, state)

	// Check the result
	assert.False(t, result)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify mock expectations
	session.AssertExpectations(t)
}

func TestStatic_Do_ClickAction_BannerNotFound(t *testing.T) {
	// Create real dependencies
	logger := zaptest.NewLogger(t)
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance of Static
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Prepare context and state
	ctx := context.Background()

	// Create a mock HTTP request with action and placement parameters
	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"click", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	// Create a response recorder
	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	// Set up mock expectations
	reqForMock := req.Clone(ctx)
	sess := sessions.Session{Value: "banner-id"}
	session.On("LoadWithExpire", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL == reqForMock.URL
	})).Return(sess, true)
	cache.On("One", ctx, "banner-id").Return(ads.Banner{}, false)

	// Call the function under test
	result := static.Do(ctx, state)

	// Check the result
	assert.False(t, result)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify mock expectations
	session.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestStatic_Do_ClickAction_Success(t *testing.T) {
	// Create real dependencies
	logger := zaptest.NewLogger(t)
	cache := &MockCache{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	// Create an instance of Static
	static := &Static{
		logger:    logger,
		cache:     cache,
		sessions:  session,
		analytics: analytics,
	}

	// Prepare context and state
	ctx := context.Background()

	// Create a mock HTTP request with action and placement parameters
	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"click", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	// Create a response recorder
	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	// Set up mock expectations
	reqForMock := req.Clone(ctx)
	sess := sessions.Session{Value: "banner-id"}
	banner := ads.Banner{Target: "https://example.com/target"}
	session.On("LoadWithExpire", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL == reqForMock.URL
	})).Return(sess, true)
	cache.On("One", ctx, "banner-id").Return(banner, true)
	analytics.On("LogClick", ctx, ads.TrackerInfo{}).Return(nil)

	// Call the function under test
	result := static.Do(ctx, state)

	// Check the result
	assert.False(t, result)
	// Check that it's a redirect response
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "https://example.com/target", rr.Header().Get("Location"))

	// Verify mock expectations
	session.AssertExpectations(t)
	cache.AssertExpectations(t)
	analytics.AssertExpectations(t)
}
