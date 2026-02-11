package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

// MockFormat is a mock implementation of the plugins.Format interface
type MockFormat struct {
	mock.Mock
}

func (m *MockFormat) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFormat) Copy(cfg map[string]any) plugins.Format {
	args := m.Called(cfg)
	return args.Get(0).(plugins.Format)
}

func (m *MockFormat) Render(ctx context.Context, state *plugins.State) (any, error) {
	args := m.Called(ctx, state)
	return args.Get(0), args.Error(1)
}

// MockAnalytics is a mock implementation of the Analytics interface
type MockAnalytics struct {
	mock.Mock
}

func (m *MockAnalytics) LogResponse(ctx context.Context, w ads.Banner, state *plugins.State) error {
	args := m.Called(ctx, w, state)
	return args.Error(0)
}

// MockResponseWriter is a mock implementation of http.ResponseWriter
type MockResponseWriter struct {
	mock.Mock
	Buffer bytes.Buffer
}

func (m *MockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	m.Buffer.Write(data)
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

// Helper function to create a Web instance with mock analytics
func newWebWithMockAnalytics(formats []plugins.Format, analytics Analytics) *Web {
	web := &Web{
		analytics: analytics,
		formats:   make(map[string]plugins.Format),
	}

	for _, f := range formats {
		web.formats[f.Name()] = f
	}

	return web
}

func TestWeb_Do_FormatNotFound(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")

	formats := []plugins.Format{mockFormat}
	mockAnalytics := new(MockAnalytics)

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "nonexistent" // Set to a format that doesn't exist

	ctx := context.Background()
	state := &plugins.State{
		Response: nil,
		Winners:  []ads.Banner{},
	}

	// Act
	err := web.Do(ctx, state)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "format nonexistent not found")
	mockFormat.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
}

func TestWeb_Do_RenderError(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")
	mockFormat.On("Render", mock.Anything, mock.Anything).Return(nil, errors.New("render error"))

	formats := []plugins.Format{mockFormat}
	mockAnalytics := new(MockAnalytics)

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "banner"

	ctx := context.Background()
	state := &plugins.State{
		Response: nil,
		Winners:  []ads.Banner{},
	}

	// Act
	err := web.Do(ctx, state)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error on render format")
	mockFormat.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
}

func TestWeb_Do_JSONMarshalError(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")

	// Create a structure that can't be marshaled
	invalidData := make(chan int)
	mockFormat.On("Render", mock.Anything, mock.Anything).Return(invalidData, nil)

	formats := []plugins.Format{mockFormat}
	mockAnalytics := new(MockAnalytics)

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "banner"

	ctx := context.Background()
	state := &plugins.State{
		Response: nil,
		Winners:  []ads.Banner{},
	}

	// Act
	err := web.Do(ctx, state)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error on render format")
	mockFormat.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
}

func TestWeb_Do_Success(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")

	responseData := []map[string]string{
		{"title": "Test Banner", "img": "http://example.com/image.jpg"},
	}
	mockFormat.On("Render", mock.Anything, mock.Anything).Return(responseData, nil)

	formats := []plugins.Format{mockFormat}
	mockAnalytics := new(MockAnalytics)

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "banner"

	mockResponseWriter := new(MockResponseWriter)
	mockResponseWriter.On("Write", mock.Anything).Return(len(responseData), nil)

	ctx := context.Background()
	state := &plugins.State{
		Response: mockResponseWriter,
		Winners:  []ads.Banner{},
	}

	// Act
	err := web.Do(ctx, state)

	// Assert
	assert.NoError(t, err)

	// Check that the response was written correctly
	expectedData, _ := json.Marshal(responseData)
	assert.Equal(t, string(expectedData), mockResponseWriter.Buffer.String())

	mockFormat.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestWeb_Do_WithWinners_AnalyticsCalled(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")

	responseData := []map[string]string{
		{"title": "Test Banner", "img": "http://example.com/image.jpg"},
	}
	mockFormat.On("Render", mock.Anything, mock.Anything).Return(responseData, nil)

	mockAnalytics := new(MockAnalytics)
	mockAnalytics.On("LogResponse", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	formats := []plugins.Format{mockFormat}

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "banner"

	mockResponseWriter := new(MockResponseWriter)
	mockResponseWriter.On("Write", mock.Anything).Return(len(responseData), nil)

	winnerBanner := ads.Banner{
		ID:    "1",
		Title: "Test Banner",
		Price: 100,
	}

	ctx := context.Background()
	state := &plugins.State{
		Response:  mockResponseWriter,
		Winners:   []ads.Banner{winnerBanner},
		RequestID: "test-request-id",
		ClickID:   "test-click-id",
	}

	// Act
	err := web.Do(ctx, state)

	// Assert
	assert.NoError(t, err)

	// Check that analytics was called with the correct parameters
	mockAnalytics.AssertCalled(t, "LogResponse", ctx, winnerBanner, state)

	mockFormat.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestWeb_Do_WithoutWinners_AnalyticsNotCalled(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")

	responseData := []map[string]string{
		{"title": "Test Banner", "img": "http://example.com/image.jpg"},
	}
	mockFormat.On("Render", mock.Anything, mock.Anything).Return(responseData, nil)

	mockAnalytics := new(MockAnalytics)

	formats := []plugins.Format{mockFormat}

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "banner"

	mockResponseWriter := new(MockResponseWriter)
	mockResponseWriter.On("Write", mock.Anything).Return(len(responseData), nil)

	ctx := context.Background()
	state := &plugins.State{
		Response: mockResponseWriter,
		Winners:  []ads.Banner{}, // Empty winners list
	}

	// Act
	err := web.Do(ctx, state)

	// Assert
	assert.NoError(t, err)

	// Check that analytics was not called
	mockAnalytics.AssertNotCalled(t, "LogResponse")

	mockFormat.AssertExpectations(t)
	mockAnalytics.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestWeb_Name(t *testing.T) {
	// Arrange
	web := &Web{}

	// Act
	name := web.Name()

	// Assert
	assert.Equal(t, "outputs.web", name)
}

func TestWeb_Copy(t *testing.T) {
	// Arrange
	mockFormat := new(MockFormat)
	mockFormat.On("Name").Return("banner")
	mockFormat.On("Copy", mock.Anything).Return(mockFormat)

	formats := []plugins.Format{mockFormat}
	mockAnalytics := new(MockAnalytics)

	web := newWebWithMockAnalytics(formats, mockAnalytics)
	web.format = "banner"

	cfg := map[string]any{
		"format": "banner",
	}

	// Act
	copied := web.Copy(cfg)

	// Assert
	assert.NotNil(t, copied)
	assert.IsType(t, &Web{}, copied)

	// Check that the copied web has the correct format
	copiedWeb := copied.(*Web)
	assert.Equal(t, "banner", copiedWeb.format)

	mockFormat.AssertExpectations(t)
}
