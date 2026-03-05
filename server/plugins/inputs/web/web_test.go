package web

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

// MockAnalytics is a mock implementation of the Analytics interface
type MockAnalytics struct {
	mock.Mock
}

func (m *MockAnalytics) LogRequest(ctx context.Context, state *plugins.State) error {
	args := m.Called(ctx, state)
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

func TestWeb_Name(t *testing.T) {
	web := &Web{}

	name := web.Name()

	assert.Equal(t, "inputs.web", name)
}

func TestWeb_Copy(t *testing.T) {
	analytics := &MockAnalytics{}
	placements := &MockPlacements{}
	units := &MockUnits{}

	web := &Web{
		analytics:  analytics,
		placements: placements,
		units:      units,
		logger:     zap.NewNop(),
	}

	cfgMap := map[string]any{"key": "value"}

	copied := web.Copy(cfgMap)

	assert.NotNil(t, copied)
	assert.IsType(t, &Web{}, copied)

	w := copied.(*Web)
	assert.Equal(t, analytics, w.analytics)
	assert.Equal(t, placements, w.placements)
	assert.Equal(t, units, w.units)
}

func TestWeb_Do(t *testing.T) {
	analytics := &MockAnalytics{}
	placements := &MockPlacements{}
	units := &MockUnits{}

	analytics.On("LogRequest", mock.Anything, mock.Anything).Return(nil)

	web := &Web{
		analytics:  analytics,
		placements: placements,
		units:      units,
		logger:     zap.NewNop(),
	}

	placements.On("One", mock.Anything, mock.Anything).Return(ads.Placement{
		ID: 1,
	}, true)
	units.On("FindByPlacement", mock.Anything, mock.Anything).Return([]ads.Unit{
		{
			ID:    1,
			Price: 10,
		},
	}, true)

	// Подготавливаем контекст и состояние
	ctx := context.Background()

	// Создаем mock HTTP запрос с параметром placement
	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"placement"},
			Values: []string{"test-placement"},
		},
	}
	req := &http.Request{}
	req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

	state := &plugins.State{
		Request: req,
	}

	// Вызываем тестируемую функцию
	result := web.Do(ctx, state)

	// Проверяем результат
	assert.True(t, result)
	assert.NotNil(t, state.User)
	assert.NotNil(t, state.Device)
	assert.NotNil(t, state.Placement)
	assert.Equal(t, uint(1), state.Placement.ID)

	// Проверяем, что placement содержит единицу рекламы
	assert.Len(t, state.Units, 1)
	assert.Equal(t, uint(1), state.Units[0].ID)
	assert.Equal(t, 10, state.Units[0].Price)

	// Проверяем, что analytics.LogRequest был вызван
	analytics.AssertCalled(t, "LogRequest", ctx, state)
}
