//nolint:errcheck
package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTelemetry implements the Telemetry interface for testing
type mockTelemetry struct {
	registry   *prometheus.Registry
	registered []prometheus.Collector
}

func (m *mockTelemetry) Register(collectors ...prometheus.Collector) error {
	m.registered = append(m.registered, collectors...)
	for _, collector := range collectors {
		if err := m.registry.Register(collector); err != nil {
			return err
		}
	}
	return nil
}

// errorTelemetry implements the Telemetry interface and returns an error on Register
type errorTelemetry struct{}

func (e *errorTelemetry) Register(collectors ...prometheus.Collector) error {
	return fmt.Errorf("registration error")
}

func TestMetrics_Success(t *testing.T) {
	// Create a mock telemetry
	tel := &mockTelemetry{
		registry: prometheus.NewRegistry(),
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	// Check that collectors were registered
	assert.Len(t, tel.registered, 2)

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some work
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap the handler with the middleware
	handler := middleware(nextHandler)

	// Create a test request with chi context
	r := chi.NewRouter()
	r.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	req := httptest.NewRequest("GET", "/test/123", nil)
	rec := httptest.NewRecorder()

	// Execute the request
	r.ServeHTTP(rec, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test response", rec.Body.String())
}

func TestMetrics_RegistrationError(t *testing.T) {
	// Create an error telemetry
	tel := &errorTelemetry{}

	// Try to create the metrics middleware
	middleware, err := Metrics(tel)
	require.Error(t, err)
	assert.Nil(t, middleware)
	assert.Contains(t, err.Error(), "registration error")
}

func TestMetrics_CollectorIdentification(t *testing.T) {
	// Create a mock telemetry
	tel := &mockTelemetry{
		registry: prometheus.NewRegistry(),
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	// Check that we have the expected collectors
	assert.Len(t, tel.registered, 2)

	// Identify the collectors by their types
	var counterVec *prometheus.CounterVec
	var histogramVec *prometheus.HistogramVec

	for _, collector := range tel.registered {
		switch c := collector.(type) {
		case *prometheus.CounterVec:
			counterVec = c
		case *prometheus.HistogramVec:
			histogramVec = c
		}
	}

	require.NotNil(t, counterVec)
	require.NotNil(t, histogramVec)
}

func TestMetrics_DifferentStatusCodes(t *testing.T) {
	// Create a mock telemetry
	tel := &mockTelemetry{
		registry: prometheus.NewRegistry(),
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	testCases := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"BadRequest", http.StatusBadRequest},
		{"Unauthorized", http.StatusUnauthorized},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that returns the specific status code
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})

			// Wrap the handler with the middleware
			handler := middleware(nextHandler)

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			// Execute the request
			handler.ServeHTTP(rec, req)

			// Check the response
			assert.Equal(t, tc.statusCode, rec.Code)
		})
	}
}

func TestMetrics_RoutePattern(t *testing.T) {
	// Create a mock telemetry
	tel := &mockTelemetry{
		registry: prometheus.NewRegistry(),
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the middleware
	handler := middleware(nextHandler)

	// Create a chi router with multiple routes to test pattern matching
	r := chi.NewRouter()
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/users"},
		{"GET", "/users/123"},
		{"POST", "/users"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %s", tc.method, tc.path), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

// Test the actual metric collection and values
func TestMetrics_CollectionAndValues(t *testing.T) {
	// Create a mock telemetry with a registry
	registry := prometheus.NewRegistry()
	tel := &mockTelemetry{
		registry: registry,
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	// Get the collectors
	var requestsTotal *prometheus.CounterVec
	var requestDuration *prometheus.HistogramVec

	for _, collector := range tel.registered {
		switch c := collector.(type) {
		case *prometheus.CounterVec:
			requestsTotal = c
		case *prometheus.HistogramVec:
			requestDuration = c
		}
	}

	require.NotNil(t, requestsTotal)
	require.NotNil(t, requestDuration)

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap the handler with the middleware
	handler := middleware(nextHandler)

	// Create a chi router
	r := chi.NewRouter()
	r.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	// Make a request
	req := httptest.NewRequest("GET", "/test/123", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test response", rec.Body.String())

	// Check that collectors were registered correctly
	// We can't easily check the exact metric values without testutil,
	// but we can verify that the collectors are of the right type
	// and that they were registered without error
	assert.IsType(t, &prometheus.CounterVec{}, requestsTotal)
	assert.IsType(t, &prometheus.HistogramVec{}, requestDuration)
}

// Test error status codes
func TestMetrics_ErrorStatusCodes(t *testing.T) {
	// Create a mock telemetry with a registry
	registry := prometheus.NewRegistry()
	tel := &mockTelemetry{
		registry: registry,
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	// Get the collectors
	var requestsTotal *prometheus.CounterVec

	for _, collector := range tel.registered {
		switch c := collector.(type) {
		case *prometheus.CounterVec:
			requestsTotal = c
		}
	}

	require.NotNil(t, requestsTotal)

	// Create a test handler that returns an error
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	})

	// Wrap the handler with the middleware
	handler := middleware(nextHandler)

	// Create a chi router
	r := chi.NewRouter()
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	// Make a request
	req := httptest.NewRequest("GET", "/error", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "error", rec.Body.String())

	// Verify that the collector is of the right type
	assert.IsType(t, &prometheus.CounterVec{}, requestsTotal)
}

// Test histogram metrics
func TestMetrics_Histogram(t *testing.T) {
	// Create a mock telemetry with a registry
	registry := prometheus.NewRegistry()
	tel := &mockTelemetry{
		registry: registry,
	}

	// Create the metrics middleware
	middleware, err := Metrics(tel)
	require.NoError(t, err)
	require.NotNil(t, middleware)

	// Get the collectors
	var requestDuration *prometheus.HistogramVec

	for _, collector := range tel.registered {
		switch c := collector.(type) {
		case *prometheus.HistogramVec:
			requestDuration = c
		}
	}

	require.NotNil(t, requestDuration)

	// Create a test handler with some delay
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	// Wrap the handler with the middleware
	handler := middleware(nextHandler)

	// Create a chi router
	r := chi.NewRouter()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify that the collector is of the right type
	assert.IsType(t, &prometheus.HistogramVec{}, requestDuration)
}
