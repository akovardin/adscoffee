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

type testObj struct {
	Budget string `json:"budget"`
}

func TestBudgetComponent(t *testing.T) {
	logger := zap.NewNop()
	budgetComponent := NewBudget(logger)

	t.Run("Component with valid budget data", func(t *testing.T) {
		obj := &testObj{
			Budget: `{"impressions":{"daily":1000,"total":10000,"uniform":true},"clicks":{"daily":100,"total":1000,"uniform":false},"money":{"daily":5000,"total":50000,"uniform":true},"conversions":{"daily":10,"total":100,"uniform":false}}`,
		}

		field := &presets.FieldContext{
			Name:  "Budget",
			Label: "Budget",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := budgetComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Check that all labels are present
		assert.Contains(t, string(html), "Показы")
		assert.Contains(t, string(html), "Клики")
		assert.Contains(t, string(html), "Конверсии")
		assert.Contains(t, string(html), "Деньги")

		// Check that all input fields are present
		assert.Contains(t, string(html), "Budget.Impressions.Daily")
		assert.Contains(t, string(html), "Budget.Impressions.Total")
		assert.Contains(t, string(html), "Budget.Impressions.Uniform")
		assert.Contains(t, string(html), "Budget.Clicks.Daily")
		assert.Contains(t, string(html), "Budget.Clicks.Total")
		assert.Contains(t, string(html), "Budget.Clicks.Uniform")
		assert.Contains(t, string(html), "Budget.Conversions.Daily")
		assert.Contains(t, string(html), "Budget.Conversions.Total")
		assert.Contains(t, string(html), "Budget.Conversions.Uniform")
		assert.Contains(t, string(html), "Budget.Money.Daily")
		assert.Contains(t, string(html), "Budget.Money.Total")
		assert.Contains(t, string(html), "Budget.Money.Uniform")
	})

	t.Run("Component with invalid budget data", func(t *testing.T) {
		obj := &testObj{
			Budget: `invalid json`,
		}

		field := &presets.FieldContext{
			Name:  "Budget",
			Label: "Budget",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := budgetComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should still render the component even with invalid data
		assert.Contains(t, string(html), "budget-field")
	})

	t.Run("Component with empty budget data", func(t *testing.T) {
		obj := &testObj{
			Budget: "",
		}

		field := &presets.FieldContext{
			Name:  "Budget",
			Label: "Budget",
		}

		req := httptest.NewRequest("GET", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		component := budgetComponent.Component(obj, field, ctx)
		require.NotNil(t, component)

		// Render component to HTML
		html, err := component.MarshalHTML(ctx.R.Context())
		require.NoError(t, err)

		// Should render the component with default empty values
		assert.Contains(t, string(html), "budget-field")
	})
}

func TestBudgetSetter(t *testing.T) {
	logger := zap.NewNop()
	budgetComponent := NewBudget(logger)

	t.Run("Setter with valid form data", func(t *testing.T) {
		obj := &testObj{
			Budget: `{"impressions":{"daily":1000,"total":10000,"uniform":true},"clicks":{"daily":100,"total":1000,"uniform":false},"money":{"daily":5000,"total":50000,"uniform":true},"conversions":{"daily":10,"total":100,"uniform":false}}`,
		}

		field := &presets.FieldContext{
			Name:  "Budget",
			Label: "Budget",
		}

		form := make(map[string][]string)
		form["Budget.Impressions.Daily"] = []string{"2000"}
		form["Budget.Impressions.Total"] = []string{"20000"}
		form["Budget.Impressions.Uniform"] = []string{"true"}
		form["Budget.Clicks.Daily"] = []string{"200"}
		form["Budget.Clicks.Total"] = []string{"2000"}
		form["Budget.Clicks.Uniform"] = []string{"false"}
		form["Budget.Conversions.Daily"] = []string{"20"}
		form["Budget.Conversions.Total"] = []string{"200"}
		form["Budget.Conversions.Uniform"] = []string{"true"}
		form["Budget.Money.Daily"] = []string{"10000"}
		form["Budget.Money.Total"] = []string{"100000"}
		form["Budget.Money.Uniform"] = []string{"false"}

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

		err := budgetComponent.Setter(obj, field, ctx)
		assert.NoError(t, err)

		// Parse the updated budget to verify values
		updatedBudget, err := models.NewBudget(obj.Budget)
		assert.NoError(t, err)

		assert.Equal(t, 2000, updatedBudget.Impressions.Daily)
		assert.Equal(t, 20000, updatedBudget.Impressions.Total)
		assert.Equal(t, true, updatedBudget.Impressions.Uniform)
		assert.Equal(t, 200, updatedBudget.Clicks.Daily)
		assert.Equal(t, 2000, updatedBudget.Clicks.Total)
		assert.Equal(t, false, updatedBudget.Clicks.Uniform)
		assert.Equal(t, 20, updatedBudget.Conversions.Daily)
		assert.Equal(t, 200, updatedBudget.Conversions.Total)
		assert.Equal(t, true, updatedBudget.Conversions.Uniform)
		assert.Equal(t, 10000, updatedBudget.Money.Daily)
		assert.Equal(t, 100000, updatedBudget.Money.Total)
		assert.Equal(t, false, updatedBudget.Money.Uniform)
	})

	t.Run("Setter with invalid form data", func(t *testing.T) {
		obj := &testObj{
			Budget: "",
		}

		field := &presets.FieldContext{
			Name:  "Budget",
			Label: "Budget",
		}

		form := make(map[string][]string)
		form["Budget.Impressions.Daily"] = []string{"invalid"}

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

		err := budgetComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
	})

	t.Run("Setter with non-string field value", func(t *testing.T) {
		obj := &struct {
			Budget int `json:"budget"`
		}{
			Budget: 123,
		}

		field := &presets.FieldContext{
			Name:  "Budget",
			Label: "Budget",
		}

		req := httptest.NewRequest("POST", "/", nil)
		ctx := &web.EventContext{
			R: req,
		}

		err := budgetComponent.Setter(obj, field, ctx)
		assert.Error(t, err)
		assert.Equal(t, "budget field value is not string", err.Error())
	})
}
