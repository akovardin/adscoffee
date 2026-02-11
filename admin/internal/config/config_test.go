//nolint:errcheck,staticcheck
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
server:
  port: ":9090"

database:
  debug: true
  user: testuser
  password: testpass
  host: localhost
  port: 5432
  dbname: testdb

s3storage:
  bucket: "test-bucket"
  s3Endpoint: "http://localhost:9000"
  accessId: "test-access-id"
  accessKey: "test-access-key"
  region: "us-east-1"
  s3ForcePathStyle: true
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
	assert.Equal(t, ":9090", cfg.Server.Port)

	assert.Equal(t, true, cfg.Database.Debug)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "testdb", cfg.Database.Dbname)

	assert.Equal(t, "test-bucket", cfg.S3Storage.Bucket)
	assert.Equal(t, "http://localhost:9000", cfg.S3Storage.S3Endpoint)
	assert.Equal(t, "test-access-id", cfg.S3Storage.AccessID)
	assert.Equal(t, "test-access-key", cfg.S3Storage.AccessKey)
	assert.Equal(t, "us-east-1", cfg.S3Storage.Region)
	assert.Equal(t, true, cfg.S3Storage.S3ForcePathStyle)
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
server:
  port: ":8080"
 database:  # Incorrect indentation
   debug: true
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
	assert.Equal(t, "", cfg.Server.Port)
	assert.Equal(t, "", cfg.Database.User)
	assert.Equal(t, "", cfg.S3Storage.Bucket)
}
