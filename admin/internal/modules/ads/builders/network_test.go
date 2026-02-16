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

func TestNetwork_Configure(t *testing.T) {
	// Create a logger for testing
	logger := zap.NewNop()

	// Create a mock database (we won't actually use it in this test)
	var db *gorm.DB

	// Create a new Network instance
	network := &Network{
		logger: logger,
		db:     db,
	}

	// Create a presets builder
	b := presets.New()

	// Call the Configure method
	network.Configure(b)

	// Build the presets to initialize everything
	b.Build()

	// Get the model builder for Network
	modelBuilder := b.Model(&models.Network{})

	// Check that the model builder was configured correctly
	assert.NotNil(t, modelBuilder)

	// Check that the listing fields are set correctly
	listing := modelBuilder.Listing()
	fieldNames := listing.FieldNames()
	assert.Len(t, fieldNames, 2)
	assert.Contains(t, fieldNames, "Title")
	assert.Contains(t, fieldNames, "Name")

	// Check that the listing fields are set correctly
	listing = modelBuilder.Listing()
	fieldNames = listing.FieldNames()
	assert.Len(t, fieldNames, 2)
	assert.Contains(t, fieldNames, "Title")
	assert.Contains(t, fieldNames, "Name")
}

func TestNetwork_Validation(t *testing.T) {
	// Create a logger for testing
	logger := zap.NewNop()

	// Create a mock database (we won't actually use it in this test)
	var db *gorm.DB

	// Create a new Network instance
	network := &Network{
		logger: logger,
		db:     db,
	}

	// Create a presets builder
	b := presets.New()

	// Configure the network
	mb := network.Configure(b)

	// Build the presets to initialize everything
	b.Build()

	// Test with a valid network (title and name are not empty)
	t.Run("Valid network", func(t *testing.T) {
		validNetwork := &models.Network{
			Title: "Test Network",
			Name:  "test-network",
		}

		// Create a mock event context
		ctx := &web.EventContext{}

		// Call the validate function directly
		validateFunc := mb.Editing().Validator
		err := validateFunc(validNetwork, ctx)

		// Check that there are no validation errors
		assert.False(t, err.HaveErrors())
	})

	// Test with an invalid network (title is empty)
	t.Run("Invalid network - empty title", func(t *testing.T) {
		invalidNetwork := &models.Network{
			Title: "",
			Name:  "test-network",
		}

		// Create a mock event context
		ctx := &web.EventContext{}

		// Call the validate function directly
		validateFunc := mb.Editing().Validator
		err := validateFunc(invalidNetwork, ctx)

		// Check that there is a validation error for the Title field
		require.True(t, err.HaveErrors())

		// Check that there are field errors
		titleErrors := err.GetFieldErrors("Title")
		require.Len(t, titleErrors, 1)
		assert.Equal(t, "Title is required", titleErrors[0])
	})

	// Test with an invalid network (name is empty)
	t.Run("Invalid network - empty name", func(t *testing.T) {
		invalidNetwork := &models.Network{
			Title: "Test Network",
			Name:  "",
		}

		// Create a mock event context
		ctx := &web.EventContext{}

		// Call the validate function directly
		validateFunc := mb.Editing().Validator
		err := validateFunc(invalidNetwork, ctx)

		// Check that there is a validation error for the Name field
		require.True(t, err.HaveErrors())

		// Check that there are field errors
		nameErrors := err.GetFieldErrors("Name")
		require.Len(t, nameErrors, 1)
		assert.Equal(t, "Name is required", nameErrors[0])
	})
}
