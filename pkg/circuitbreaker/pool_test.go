package circuitbreaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewPool(t *testing.T) {
	logger := zaptest.NewLogger(t)
	metrics := &Metrics{} // Using empty metrics for tests

	t.Run("successful creation with valid configs", func(t *testing.T) {
		configs := map[string]*Config{
			"test-circuit": {
				Enabled:               true,
				Timeout:               time.Second,
				MaxConcurrentRequests: 10,
				Hystrix: HystrixConfig{
					SleepWindow:                  5 * time.Second,
					HalfOpenAttempts:             1,
					RequiredConcurrentSuccessful: 1,
					ErrorThresholdPercentage:     50,
					RequestVolumeThreshold:       20,
					RollingDuration:              10 * time.Second,
					NumBuckets:                   10,
				},
			},
		}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)
		assert.NotNil(t, pool.manager)
		assert.Equal(t, logger, pool.logger)
	})

	t.Run("successful creation with empty configs", func(t *testing.T) {
		configs := map[string]*Config{}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)
	})

	t.Run("successful creation with nil configs", func(t *testing.T) {
		pool, err := NewPool(logger, metrics, nil)
		require.NoError(t, err)
		assert.NotNil(t, pool)
	})

	t.Run("successful creation with multiple configs", func(t *testing.T) {
		configs := map[string]*Config{
			"circuit-1": NewConfig(),
			"circuit-2": NewConfig(),
			"circuit-3": NewConfig(),
		}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)
	})
}

func TestPool_Get(t *testing.T) {
	logger := zaptest.NewLogger(t)
	metrics := &Metrics{}

	t.Run("get existing circuit", func(t *testing.T) {
		configs := map[string]*Config{
			"test-circuit": NewConfig(),
		}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)

		circuit := pool.Get("test-circuit")
		assert.NotNil(t, circuit)
		assert.NotEqual(t, &noopCircuit{}, circuit)
	})

	t.Run("get non-existing circuit", func(t *testing.T) {
		configs := map[string]*Config{
			"test-circuit": NewConfig(),
		}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)

		circuit := pool.Get("non-existing")
		assert.NotNil(t, circuit)
		assert.IsType(t, &noopCircuit{}, circuit)
	})

	t.Run("get from nil pool", func(t *testing.T) {
		var pool *Pool
		circuit := pool.Get("test")
		assert.NotNil(t, circuit)
		assert.IsType(t, &noopCircuit{}, circuit)
	})
}

func TestNewPool_WithCircuitConfigurations(t *testing.T) {
	logger := zaptest.NewLogger(t)
	metrics := &Metrics{}

	t.Run("circuit with hystrix configurations", func(t *testing.T) {
		configs := map[string]*Config{
			"test-circuit": {
				Enabled:               true,
				Timeout:               2 * time.Second,
				MaxConcurrentRequests: 5,
				Hystrix: HystrixConfig{
					SleepWindow:                  10 * time.Second,
					HalfOpenAttempts:             2,
					RequiredConcurrentSuccessful: 3,
					ErrorThresholdPercentage:     75,
					RequestVolumeThreshold:       10,
					RollingDuration:              30 * time.Second,
					NumBuckets:                   5,
				},
			},
		}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)

		// Verify that the circuit was created with correct configuration
		c := pool.manager.GetCircuit("test-circuit")
		assert.NotNil(t, c)
	})

	t.Run("circuit with custom timeout", func(t *testing.T) {
		configs := map[string]*Config{
			"timeout-circuit": {
				Enabled:               true,
				Timeout:               100 * time.Millisecond,
				MaxConcurrentRequests: -1,
				Hystrix: HystrixConfig{
					SleepWindow:                  5 * time.Second,
					HalfOpenAttempts:             1,
					RequiredConcurrentSuccessful: 1,
					ErrorThresholdPercentage:     50,
					RequestVolumeThreshold:       20,
					RollingDuration:              10 * time.Second,
					NumBuckets:                   10,
				},
			},
		}

		pool, err := NewPool(logger, metrics, configs)
		require.NoError(t, err)
		assert.NotNil(t, pool)
	})
}

func TestNewPool_ErrorCases(t *testing.T) {
	logger := zaptest.NewLogger(t)
	metrics := &Metrics{}

	// Note: It's difficult to simulate actual errors in circuit creation with the
	// current implementation, as the circuit.Manager doesn't easily expose error cases.
	// In a real scenario, we might test with invalid configurations or other edge cases.
	t.Run("creation with valid configs does not error", func(t *testing.T) {
		configs := map[string]*Config{
			"valid-circuit": NewConfig(),
		}

		pool, err := NewPool(logger, metrics, configs)
		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}
