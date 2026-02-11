package web

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.ads.coffee/platform/server/internal/analytics"
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

func TestNew(t *testing.T) {
	// Вызываем тестируемую функцию
	web := New(&analytics.Analytics{})

	// Проверяем результат
	assert.NotNil(t, web)
}

func TestWeb_Name(t *testing.T) {
	// Создаем экземпляр Web
	web := &Web{}

	// Вызываем тестируемую функцию
	name := web.Name()

	// Проверяем результат
	assert.Equal(t, "inputs.web", name)
}

func TestWeb_Copy(t *testing.T) {
	// Создаем mock для analytics
	mockAnalytics := new(MockAnalytics)

	// Создаем экземпляр Web
	web := &Web{
		analytics: mockAnalytics,
	}

	// Подготавливаем конфигурацию
	cfgMap := map[string]any{"key": "value"}

	// Вызываем тестируемую функцию
	copied := web.Copy(cfgMap)

	// Проверяем результат
	assert.NotNil(t, copied)
	assert.IsType(t, &Web{}, copied)

	// Проверяем, что analytics скопирован корректно
	copiedWeb := copied.(*Web)
	assert.Equal(t, mockAnalytics, copiedWeb.analytics)
}

func TestWeb_Do(t *testing.T) {
	// Создаем mock для analytics
	mockAnalytics := new(MockAnalytics)
	mockAnalytics.On("LogRequest", mock.Anything, mock.Anything).Return(nil)

	// Создаем экземпляр Web
	web := &Web{
		analytics: mockAnalytics,
	}

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
	assert.Equal(t, "test-placement", state.Placement.ID)

	// Проверяем, что placement содержит единицу рекламы
	assert.Len(t, state.Placement.Units, 1)
	assert.Equal(t, "yandex-1", state.Placement.Units[0].ID)
	assert.Equal(t, "yandex", state.Placement.Units[0].Network)
	assert.Equal(t, 10, state.Placement.Units[0].Price)
	assert.Equal(t, "banner", state.Placement.Units[0].Format)

	// Проверяем, что analytics.LogRequest был вызван
	mockAnalytics.AssertCalled(t, "LogRequest", ctx, state)
}
