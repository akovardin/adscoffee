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

type MockBanners struct {
	mock.Mock
}

func (m *MockBanners) One(ctx context.Context, id uint) (ads.Banner, bool) {
	args := m.Called(ctx, id)
	banner, _ := args.Get(0).(ads.Banner)
	return banner, args.Bool(1)
}

type MockSession struct {
	mock.Mock
}

func (m *MockSession) LoadWithExpire(r *http.Request) (sessions.Session, bool) {
	args := m.Called(r)
	session, _ := args.Get(0).(sessions.Session)
	return session, args.Bool(1)
}

type MockAnalytics struct {
	mock.Mock
}

func (m *MockAnalytics) LogClick(ctx context.Context, data ads.TrackerInfo) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type MockPlacements struct {
	mock.Mock
}

func (m *MockPlacements) One(ctx context.Context, id uint) (ads.Placement, bool) {
	args := m.Called(ctx, id)
	placement, _ := args.Get(0).(ads.Placement)
	return placement, args.Bool(1)
}

type MockUnits struct {
	mock.Mock
}

func (m *MockUnits) FindByPlacement(ctx context.Context, id uint) ([]ads.Unit, bool) {
	args := m.Called(ctx, id)
	units, _ := args.Get(0).([]ads.Unit)
	return units, args.Bool(1)
}

func TestNew(t *testing.T) {
	logger := zaptest.NewLogger(t)

	banners := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	static := &Static{
		logger:    logger,
		banners:   banners,
		sessions:  session,
		analytics: analytics,
	}

	assert.NotNil(t, static)
}

func TestStatic_Name(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cache := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	static := &Static{
		logger:    logger,
		banners:   cache,
		sessions:  session,
		analytics: analytics,
	}

	name := static.Name()

	assert.Equal(t, "inputs.static", name)
}

func TestStatic_Copy(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cache := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}

	static := &Static{
		logger:    logger,
		banners:   cache,
		sessions:  session,
		analytics: analytics,
	}

	cfgMap := map[string]any{"key": "value"}

	copied := static.Copy(cfgMap)

	assert.NotNil(t, copied)
	assert.IsType(t, &Static{}, copied)

	copiedStatic := copied.(*Static)
	assert.Equal(t, cache, copiedStatic.banners)
	assert.Equal(t, logger, copiedStatic.logger)
	assert.Equal(t, session, copiedStatic.sessions)
	assert.Equal(t, analytics, copiedStatic.analytics)
}

func TestStatic_Do_ViewAction(t *testing.T) {
	logger := zaptest.NewLogger(t)
	banners := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}
	placements := &MockPlacements{}
	units := &MockUnits{}

	static := &Static{
		logger:     logger,
		banners:    banners,
		sessions:   session,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}

	placements.On("One", mock.Anything, mock.Anything).Return(ads.Placement{
		ID: 1,
	}, true)
	units.On("FindByPlacement", mock.Anything, mock.Anything).Return([]ads.Unit{
		{
			ID: 1,
		},
	}, true)

	ctx := context.Background()

	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"view", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	result := static.Do(ctx, state)

	assert.True(t, result)
	assert.NotNil(t, state.User)
	assert.NotNil(t, state.Device)
	assert.NotNil(t, state.Placement)
	assert.Equal(t, uint(1), state.Placement.ID)

	action := state.Request.Context().Value("action")
	assert.Equal(t, "view", action)
}

func TestStatic_Do_ClickAction_SessionNotFound(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cache := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}
	placements := &MockPlacements{}
	units := &MockUnits{}

	static := &Static{
		logger:     logger,
		banners:    cache,
		sessions:   session,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}

	placements.On("One", mock.Anything, mock.Anything).Return(ads.Placement{
		ID: 1,
	}, true)
	units.On("FindByPlacement", mock.Anything, mock.Anything).Return(nil, false)

	ctx := context.Background()

	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"click", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	reqForMock := req.Clone(ctx)
	session.On("LoadWithExpire", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL == reqForMock.URL
	})).Return(sessions.Session{}, false)

	result := static.Do(ctx, state)

	assert.False(t, result)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	session.AssertExpectations(t)
}

func TestStatic_Do_ClickAction_BannerNotFound(t *testing.T) {
	logger := zaptest.NewLogger(t)
	banners := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}
	placements := &MockPlacements{}
	units := &MockUnits{}

	static := &Static{
		logger:     logger,
		banners:    banners,
		sessions:   session,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}

	placements.On("One", mock.Anything, mock.Anything).Return(ads.Placement{
		ID: 1,
	}, true)
	units.On("FindByPlacement", mock.Anything, mock.Anything).Return(nil, false)

	ctx := context.Background()

	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"click", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	reqForMock := req.Clone(ctx)
	sess := sessions.Session{Value: "1"}
	session.On("LoadWithExpire", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL == reqForMock.URL
	})).Return(sess, true)
	banners.On("One", ctx, uint(1)).Return(ads.Banner{}, false)

	result := static.Do(ctx, state)

	assert.False(t, result)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	session.AssertExpectations(t)
	banners.AssertExpectations(t)
}

func TestStatic_Do_ClickAction_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	banners := &MockBanners{}
	session := &MockSession{}
	analytics := &MockAnalytics{}
	placements := &MockPlacements{}
	units := &MockUnits{}

	static := &Static{
		logger:     logger,
		banners:    banners,
		sessions:   session,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}

	placements.On("One", mock.Anything, mock.Anything).Return(ads.Placement{
		ID: 1,
	}, true)
	units.On("FindByPlacement", mock.Anything, mock.Anything).Return(nil, false)

	ctx := context.Background()

	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"action", "placement"},
			Values: []string{"click", "test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	state := &plugins.State{
		Request:  req,
		Response: rr,
	}

	reqForMock := req.Clone(ctx)
	sess := sessions.Session{Value: "1"}
	banner := ads.Banner{Target: "https://example.com/target"}
	session.On("LoadWithExpire", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL == reqForMock.URL
	})).Return(sess, true)
	banners.On("One", ctx, uint(1)).Return(banner, true)
	analytics.On("LogClick", ctx, ads.TrackerInfo{}).Return(nil)

	result := static.Do(ctx, state)

	assert.False(t, result)
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "https://example.com/target", rr.Header().Get("Location"))

	session.AssertExpectations(t)
	banners.AssertExpectations(t)
	analytics.AssertExpectations(t)
}
