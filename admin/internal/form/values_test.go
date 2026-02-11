//nolint:errcheck,staticcheck
package form

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	// Test case 1: Form is nil, key exists in form data
	t.Run("Form is nil, key exists in form data", func(t *testing.T) {
		// Create a request with form data
		form := url.Values{}
		form.Add("testKey", "testValue1")
		form.Add("testKey", "testValue2")

		req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Call Values function
		result := Values(req, "testKey")

		// Assertions
		assert.NotNil(t, result)
		assert.Equal(t, []string{"testValue1", "testValue2"}, result)
	})

	// Test case 2: Form is nil, key does not exist in form data
	t.Run("Form is nil, key does not exist in form data", func(t *testing.T) {
		// Create a request with form data
		form := url.Values{}
		form.Add("otherKey", "otherValue")

		req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Call Values function
		result := Values(req, "testKey")

		// Assertions
		assert.Nil(t, result)
	})

	// Test case 3: Form is already parsed, key exists
	t.Run("Form is already parsed, key exists", func(t *testing.T) {
		// Create a request with form data
		form := url.Values{}
		form.Add("testKey", "testValue1")
		form.Add("testKey", "testValue2")

		req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Parse form manually to simulate already parsed form
		req.ParseForm()

		// Call Values function
		result := Values(req, "testKey")

		// Assertions
		assert.NotNil(t, result)
		assert.Equal(t, []string{"testValue1", "testValue2"}, result)
	})

	// Test case 4: Form is already parsed, key does not exist
	t.Run("Form is already parsed, key does not exist", func(t *testing.T) {
		// Create a request with form data
		form := url.Values{}
		form.Add("otherKey", "otherValue")

		req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Parse form manually to simulate already parsed form
		req.ParseForm()

		// Call Values function
		result := Values(req, "testKey")

		// Assertions
		assert.Nil(t, result)
	})

	// Test case 5: Multipart form data
	t.Run("Multipart form data", func(t *testing.T) {
		// Create a request with multipart form data
		body := `--boundary
Content-Disposition: form-data; name="testKey"

testValue1
--boundary
Content-Disposition: form-data; name="testKey"

testValue2
--boundary--`

		req := httptest.NewRequest("POST", "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

		// Call Values function
		result := Values(req, "testKey")

		// Assertions
		assert.NotNil(t, result)
		assert.Equal(t, []string{"testValue1", "testValue2"}, result)
	})

	// Test case 6: Empty form values
	t.Run("Empty form values", func(t *testing.T) {
		// Create a request with empty form data for a key
		form := url.Values{}
		form.Add("testKey", "") // Empty value

		req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Call Values function
		result := Values(req, "testKey")

		// Assertions
		assert.NotNil(t, result)
		assert.Equal(t, []string{""}, result)
	})
}
