package builders

import (
	"testing"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

func TestUnit_Configure(t *testing.T) {
	// Create a logger for testing
	logger := zap.NewNop()

	// Create a mock database (we won't actually use it in this test)
	var db *gorm.DB

	// Create a new Unit instance
	unit := &Unit{
		logger: logger,
		db:     db,
	}

	// Create a presets builder
	b := presets.New()

	// Call the Configure method
	unit.Configure(b)

	// Build the presets to initialize everything
	b.Build()

	// Get the model builder for Unit
	modelBuilder := b.Model(&models.Unit{})

	// Check that the model builder was configured correctly
	assert.NotNil(t, modelBuilder)

	// Check that the listing fields are set correctly
	listing := modelBuilder.Listing()
	assert.Equal(t, []interface{}{"Name", "Price", "NetworkID", "PlacementID", "Data", "Active"}, listing.FieldNames())
}

func TestUnit_Validation(t *testing.T) {
	// Create a logger for testing
	logger := zap.NewNop()

	// Create a mock database (we won't actually use it in this test)
	var db *gorm.DB

	// Create a new Unit instance
	unit := &Unit{
		logger: logger,
		db:     db,
	}

	// Create a presets builder
	b := presets.New()

	// Configure the unit
	mb := unit.Configure(b)

	// Build the presets to initialize everything
	b.Build()

	// Test with a valid unit (name is not empty)
	t.Run("Valid unit", func(t *testing.T) {
		validUnit := &models.Unit{
			Name: "Test Unit",
		}

		// Create a mock event context
		ctx := &web.EventContext{}

		// Call the validate function directly
		validateFunc := mb.Editing().Validator
		err := validateFunc(validUnit, ctx)

		// Check that there are no validation errors
		assert.False(t, err.HaveErrors())
	})

	// Test with an invalid unit (name is empty)
	t.Run("Invalid unit - empty name", func(t *testing.T) {
		invalidUnit := &models.Unit{
			Name: "",
		}

		// Create a mock event context
		ctx := &web.EventContext{}

		// Call the validate function directly
		validateFunc := mb.Editing().Validator
		err := validateFunc(invalidUnit, ctx)

		// Check that there is a validation error for the Name field
		require.True(t, err.HaveErrors())

		// Check that there are field errors
		fieldErrors := err.FieldErrors()
		require.NotEmpty(t, fieldErrors)

		// Check that the error is for the Name field
		nameErrors := err.GetFieldErrors("Name")
		require.Len(t, nameErrors, 1)
		assert.Equal(t, "Name is required", nameErrors[0])
	})
}
