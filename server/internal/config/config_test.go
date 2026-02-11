package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	// Create a temporary config file with valid YAML content
	content := `
health:
  hostname: :8080
pipelines:
  - name: test
    route: /test
    input:
      name: inputs.test
    output:
     name: outputs.test

circuit-breaker:
 test:
   timeout: 1s
   max-concurrent-calls: 10
   error-percent-threshold: 50
   request-volume-threshold: 20
   sleep-window: 5s
redis-pool:
  main:
    enabled: true
    key_prefix: "test"
    cluster_addrs: ["127.0.0.1:7001"]
    pool_size: 10

telemetry:
  jaeger:
    enabled: false

database:
  user: test
  password: test
  host: localhost
  port: 5432
  dbname: test

kafka-pool:
  main:
    enabled: true
    seeds:
      - localhost:9092
    producer:
      disable_idempotent_write: true
`
	tmpfile, err := ioutil.TempFile("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	// Test the New function with the temporary file
	cfg, err := New(tmpfile.Name())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check that the config was populated correctly
	assert.Len(t, cfg.Pipelines, 1)
	assert.Equal(t, "test", cfg.Pipelines[0].Name)
	assert.Equal(t, "/test", cfg.Pipelines[0].Route)

	assert.Equal(t, ":8080", cfg.Health.Hostname)
	assert.NotNil(t, cfg.CircuitBreaker["test"])
	assert.Equal(t, "1s", cfg.CircuitBreaker["test"].Timeout.String())

	assert.NotNil(t, cfg.RedisPool["main"])
	assert.Equal(t, "test", cfg.RedisPool["main"].KeyPrefix)

	assert.Equal(t, false, cfg.Telemetry.Jaeger.Enabled)

	assert.Equal(t, "test", cfg.Database.User)
	assert.Equal(t, "test", cfg.Database.Password)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "test", cfg.Database.Dbname)
	assert.NotNil(t, cfg.Kafka["main"])
	assert.Equal(t, true, cfg.Kafka["main"].Enabled)
}

func TestNew_FileNotFound(t *testing.T) {
	// Test with a non-existent file
	_, err := New("/non/existent/file.yaml")

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestNew_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML content
	content := `
pipelines:
  - name: test
	 route: /test  # Incorrect indentation
`
	tmpfile, err := ioutil.TempFile("", "config-invalid-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	// Test the New function with the temporary file
	_, err = New(tmpfile.Name())

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't decode source: yaml")
}

func TestNew_EmptyFile(t *testing.T) {
	// Create an empty temporary config file
	tmpfile, err := ioutil.TempFile("", "config-empty-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	err = tmpfile.Close()
	require.NoError(t, err)

	// Test the New function with the empty file
	cfg, err := New(tmpfile.Name())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	// Should have default/empty values
	assert.Empty(t, cfg.Pipelines)
}
