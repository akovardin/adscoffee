//nolint:errcheck
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

type testTargetingObj struct {
	Targeting string `json:"targeting"`
}

func TestTargetingComponent(t *testing.T) {
	logger := zap.NewNop()
	targetingComponent := NewTargeting(logger)

	t.Run("Component with valid targeting data", func(t *testing.T) {
		obj := &testTargetingObj{
			Targeting: `{"bundle":{"include_or":["com.example","ru.rustore"],"exclude_or":["com.test"]},"country":{"include_or":["RU","US"],"exclude_or":["CN"]},"region":{"include_or":["SPE","MOW"],"exclude_or":["LEN"]},"city":{"include_or":["KUF","OMS"],"exclude_or":["MOW"]},"ip":{"include":["188.170.172.0/22"],"exclude":["188.170.192.0/22"]}}`,
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := targetingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Check that all section headers are present
		assert.Contains(t, string(html), "Бандлы")
		assert.Contains(t, string(html), "Страны")
		assert.Contains(t, string(html), "Регионы")
		assert.Contains(t, string(html), "Города")
		assert.Contains(t, string(html), "IP")

		// Check that all textarea fields are present with correct values
		assert.Contains(t, string(html), "Targeting.Bundle.IncludeOr")
		assert.Contains(t, string(html), "Targeting.Bundle.ExcludeOr")
		assert.Contains(t, string(html), "Targeting.Country.IncludeOr")
		assert.Contains(t, string(html), "Targeting.Country.ExcludeOr")
		assert.Contains(t, string(html), "Targeting.Region.IncludeOr")
		assert.Contains(t, string(html), "Targeting.Region.ExcludeOr")
		assert.Contains(t, string(html), "Targeting.City.IncludeOr")
		assert.Contains(t, string(html), "Targeting.City.ExcludeOr")
		assert.Contains(t, string(html), "Targeting.IP.Include")
		assert.Contains(t, string(html), "Targeting.IP.Exclude")

		// Check that values are correctly set in textareas
		assert.Contains(t, string(html), "com.example ru.rustore")
		assert.Contains(t, string(html), "com.test")
		assert.Contains(t, string(html), "RU US")
		assert.Contains(t, string(html), "CN")
		assert.Contains(t, string(html), "SPE MOW")
		assert.Contains(t, string(html), "LEN")
		assert.Contains(t, string(html), "KUF OMS")
		assert.Contains(t, string(html), "MOW")
		assert.Contains(t, string(html), "188.170.172.0/22")
		assert.Contains(t, string(html), "188.170.192.0/22")
	})

	t.Run("Component with invalid targeting data", func(t *testing.T) {
		obj := &testTargetingObj{
			Targeting: `invalid json`,
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := targetingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with invalid data
		assert.Contains(t, string(html), "targeting-field")
	})

	t.Run("Component with empty targeting data", func(t *testing.T) {
		obj := &testTargetingObj{
			Targeting: "",
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := targetingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should render the component with empty values
		assert.Contains(t, string(html), "targeting-field")
	})

	t.Run("Component with non-string field value", func(t *testing.T) {
		obj := &struct {
			Targeting int `json:"targeting"`
		}{
			Targeting: 123,
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := targetingComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with wrong type
		assert.Contains(t, string(html), "targeting-field")
	})
}

func TestTargetingSetter(t *testing.T) {
	logger := zap.NewNop()
	targetingComponent := NewTargeting(logger)

	t.Run("Setter with valid form data", func(t *testing.T) {
		obj := &testTargetingObj{
			Targeting: `{"bundle":{"include_or":["com.example","ru.rustore"],"exclude_or":["com.test"]},"country":{"include_or":["RU","US"],"exclude_or":["CN"]},"region":{"include_or":["SPE","MOW"],"exclude_or":["LEN"]},"city":{"include_or":["KUF","OMS"],"exclude_or":["MOW"]},"ip":{"include":["188.170.172.0/22"],"exclude":["188.170.192.0/22"]}}`,
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		form := make(map[string][]string)
		form["Targeting.Bundle.IncludeOr"] = []string{"com.newapp com.another"}
		form["Targeting.Bundle.ExcludeOr"] = []string{"com.exclude"}
		form["Targeting.Country.IncludeOr"] = []string{"DE FR"}
		form["Targeting.Country.ExcludeOr"] = []string{"UK"}
		form["Targeting.Region.IncludeOr"] = []string{"SPE LEN"}
		form["Targeting.Region.ExcludeOr"] = []string{"MOW"}
		form["Targeting.City.IncludeOr"] = []string{"KUF OMS"}
		form["Targeting.City.ExcludeOr"] = []string{"MOW"}
		form["Targeting.IP.Include"] = []string{"192.168.1.0/24 10.0.0.0/8"}
		form["Targeting.IP.Exclude"] = []string{"172.16.0.0/12"}

		formData := make([]string, 0)
		for key, values := range form {
			for _, value := range values {
				formData = append(formData, key+"="+value)
			}
		}

		formString := strings.Join(formData, "&")
		req := httptest.NewRequest("POST", "/", strings.NewReader(formString))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// Parse the form to populate req.Form
		req.ParseForm()

		ctx := &web.EventContext{
			R: req,
		}

		err := targetingComponent.Setter(obj, field, ctx)
		assert.NoError(t, err)

		// Parse the updated targeting to verify values
		updatedTargeting, err := models.NewTargeting(obj.Targeting)
		assert.NoError(t, err)

		assert.Equal(t, []string{"com.newapp", "com.another"}, updatedTargeting.Bundle.IncludeOr)
		assert.Equal(t, []string{"com.exclude"}, updatedTargeting.Bundle.ExcludeOr)
		assert.Equal(t, []string{"DE", "FR"}, updatedTargeting.Country.IncludeOr)
		assert.Equal(t, []string{"UK"}, updatedTargeting.Country.ExcludeOr)
		assert.Equal(t, []string{"SPE", "LEN"}, updatedTargeting.Region.IncludeOr)
		assert.Equal(t, []string{"MOW"}, updatedTargeting.Region.ExcludeOr)
		assert.Equal(t, []string{"KUF", "OMS"}, updatedTargeting.City.IncludeOr)
		assert.Equal(t, []string{"MOW"}, updatedTargeting.City.ExcludeOr)
		assert.Equal(t, []string{"192.168.1.0/24", "10.0.0.0/8"}, updatedTargeting.IP.Include)
		assert.Equal(t, []string{"172.16.0.0/12"}, updatedTargeting.IP.Exclude)
	})

	t.Run("Setter with non-string field value", func(t *testing.T) {
		obj := &struct {
			Targeting int `json:"targeting"`
		}{
			Targeting: 123,
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		req := httptest.NewRequest("POST", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		err := targetingComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
		assert.Equal(t, "budget field value is not string", err.Error())
	})

	t.Run("Setter with invalid targeting data", func(t *testing.T) {
		obj := &testTargetingObj{
			Targeting: "invalid json",
		}

		field := &presets.FieldContext{
			Name:  "Targeting",
			Label: "Targeting",
		}

		req := httptest.NewRequest("POST", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		err := targetingComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
	})
}
