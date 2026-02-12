package components

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type testCappingObj struct {
	Capping string `json:"capping"`
}

func TestCappingComponent(t *testing.T) {
	logger := zap.NewNop()
	cappingComponent := NewCapping(logger)

	t.Run("Component with valid capping data", func(t *testing.T) {
		obj := &testCappingObj{
			Capping: `{"count":1000,"period":24}`,
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := cappingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Check that all labels are present
		assert.Contains(t, string(html), "Показы")
		assert.Contains(t, string(html), "Период (часы)")

		// Check that all input fields are present
		assert.Contains(t, string(html), "Capping.Count")
		assert.Contains(t, string(html), "Capping.Period")

		// Check that values are correctly set
		assert.Contains(t, string(html), "{\"Capping.Count\":1000}")
		assert.Contains(t, string(html), "{\"Capping.Period\":24}")
	})

	t.Run("Component with invalid capping data", func(t *testing.T) {
		obj := &testCappingObj{
			Capping: `invalid json`,
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := cappingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with invalid data
		assert.Contains(t, string(html), "capping-field")
	})

	t.Run("Component with empty capping data", func(t *testing.T) {
		obj := &testCappingObj{
			Capping: "",
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := cappingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should render the component with default empty values
		assert.Contains(t, string(html), "capping-field")
	})

	t.Run("Component with non-string field value", func(t *testing.T) {
		obj := &struct {
			Capping int `json:"capping"`
		}{
			Capping: 123,
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := cappingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with wrong type
		assert.Contains(t, string(html), "capping-field")
	})
}

func TestCappingSetter(t *testing.T) {
	logger := zap.NewNop()
	cappingComponent := NewCapping(logger)

	t.Run("Setter with valid form data", func(t *testing.T) {
		obj := &testCappingObj{
			Capping: `{"count":1000,"period":24}`,
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		form := make(map[string][]string)
		form["Capping.Count"] = []string{"2000"}
		form["Capping.Period"] = []string{"48"}

		formData := make([]string, 0)
		for key, values := range form {
			for _, value := range values {
				formData = append(formData, key+"="+value)
			}
		}

		formString := strings.Join(formData, "&")
		req := httptest.NewRequest("POST", "/", strings.NewReader(formString))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := &web.EventContext{
			R: req,
		}

		err := cappingComponent.Setter(obj, field, ctx)
		assert.NoError(t, err)

		// Parse the updated capping to verify values
		updatedCapping, err := models.NewCapping(obj.Capping)
		assert.NoError(t, err)

		assert.Equal(t, 2000, updatedCapping.Count)
		assert.Equal(t, 48, updatedCapping.Period)
	})

	t.Run("Setter with invalid count value", func(t *testing.T) {
		obj := &testCappingObj{
			Capping: "",
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		form := make(map[string][]string)
		form["Capping.Count"] = []string{"invalid"}

		formData := make([]string, 0)
		for key, values := range form {
			for _, value := range values {
				formData = append(formData, key+"="+value)
			}
		}

		formString := strings.Join(formData, "&")
		req := httptest.NewRequest("POST", "/", strings.NewReader(formString))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := &web.EventContext{
			R: req,
		}

		err := cappingComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
	})

	t.Run("Setter with invalid period value", func(t *testing.T) {
		obj := &testCappingObj{
			Capping: "",
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		form := make(map[string][]string)
		form["Capping.Period"] = []string{"invalid"}

		formData := make([]string, 0)
		for key, values := range form {
			for _, value := range values {
				formData = append(formData, key+"="+value)
			}
		}

		formString := strings.Join(formData, "&")
		req := httptest.NewRequest("POST", "/", strings.NewReader(formString))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := &web.EventContext{
			R: req,
		}

		err := cappingComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
	})

	t.Run("Setter with non-string field value", func(t *testing.T) {
		obj := &struct {
			Capping int `json:"capping"`
		}{
			Capping: 123,
		}

		field := &presets.FieldContext{
			Name:  "Capping",
			Label: "Capping",
		}

		req := httptest.NewRequest("POST", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		err := cappingComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
		assert.Equal(t, "capping field value is not string", err.Error())
	})
}
