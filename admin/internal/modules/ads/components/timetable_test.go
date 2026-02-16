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

type testTimetableObj struct {
	Timetable string `json:"timetable"`
}

func TestTimetableComponent(t *testing.T) {
	logger := zap.NewNop()
	timetableComponent := NewTimetable(logger)

	t.Run("Component with valid timetable data", func(t *testing.T) {
		obj := &testTimetableObj{
			Timetable: `{"0":{"0":true,"1":false},"1":{"0":false,"1":true}}`,
		}

		field := &presets.FieldContext{
			Name:  "Timetable",
			Label: "Timetable",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := timetableComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Check that all labels are present
		assert.Contains(t, string(html), "Расписание активности")
		assert.Contains(t, string(html), "Пн")
		assert.Contains(t, string(html), "Вт")

		// Check that all input fields are present with correct values
		assert.Contains(t, string(html), "Timetable[0][0]")
		assert.Contains(t, string(html), "Timetable[0][1]")
		assert.Contains(t, string(html), "Timetable[1][0]")
		assert.Contains(t, string(html), "Timetable[1][1]")

		// Check that values are correctly set
		assert.Contains(t, string(html), "{\"Timetable[0][0]\":true}")
		assert.Contains(t, string(html), "{\"Timetable[0][1]\":false}")
		assert.Contains(t, string(html), "{\"Timetable[1][0]\":false}")
		assert.Contains(t, string(html), "{\"Timetable[1][1]\":true}")
	})

	t.Run("Component with invalid timetable data", func(t *testing.T) {
		obj := &testTimetableObj{
			Timetable: `invalid json`,
		}

		field := &presets.FieldContext{
			Name:  "Timetable",
			Label: "Timetable",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := timetableComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with invalid data
		assert.Contains(t, string(html), "timetable-field")
	})

	t.Run("Component with empty timetable data", func(t *testing.T) {
		obj := &testTimetableObj{
			Timetable: "",
		}

		field := &presets.FieldContext{
			Name:  "Timetable",
			Label: "Timetable",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := timetableComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should render the component with default values (all true)
		assert.Contains(t, string(html), "timetable-field")
	})

	t.Run("Component with non-string field value", func(t *testing.T) {
		obj := &struct {
			Timetable int `json:"timetable"`
		}{
			Timetable: 123,
		}

		field := &presets.FieldContext{
			Name:  "Timetable",
			Label: "Timetable",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := timetableComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with wrong type
		assert.Contains(t, string(html), "timetable-field")
	})
}

func TestTimetableSetter(t *testing.T) {
	logger := zap.NewNop()
	timetableComponent := NewTimetable(logger)

	t.Run("Setter with valid form data", func(t *testing.T) {
		obj := &testTimetableObj{
			Timetable: `{"0":{"0":true,"1":false},"1":{"0":false,"1":true}}`,
		}

		field := &presets.FieldContext{
			Name:  "Timetable",
			Label: "Timetable",
		}

		form := make(map[string][]string)
		form["Timetable[0][0]"] = []string{"true"}
		form["Timetable[0][1]"] = []string{"false"}
		form["Timetable[1][0]"] = []string{"false"}
		form["Timetable[1][1]"] = []string{"true"}

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

		err := timetableComponent.Setter(obj, field, ctx)
		assert.NoError(t, err)

		// Parse the updated timetable to verify values
		updatedTimetable, err := models.NewTimetable(obj.Timetable)
		assert.NoError(t, err)

		assert.Equal(t, true, updatedTimetable[0][0])
		assert.Equal(t, false, updatedTimetable[0][1])
		assert.Equal(t, false, updatedTimetable[1][0])
		assert.Equal(t, true, updatedTimetable[1][1])
	})

	t.Run("Setter with non-string field value", func(t *testing.T) {
		obj := &struct {
			Timetable int `json:"timetable"`
		}{
			Timetable: 123,
		}

		field := &presets.FieldContext{
			Name:  "Timetable",
			Label: "Timetable",
		}

		req := httptest.NewRequest("POST", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		err := timetableComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "strconv.ParseInt: parsing")
	})
}
