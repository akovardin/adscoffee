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
	logger := zap.NewNop()

	var db *gorm.DB

	unit := &Unit{
		logger: logger,
		db:     db,
	}

	b := presets.New()

	unit.Configure(b)

	b.Build()

	modelBuilder := b.Model(&models.Unit{})

	assert.NotNil(t, modelBuilder)

	listing := modelBuilder.Listing()
	assert.Equal(t, []interface{}{"Name", "Price", "NetworkID", "PlacementID", "Data", "Active"}, listing.FieldNames())
}

func TestUnit_Validation(t *testing.T) {
	logger := zap.NewNop()

	var db *gorm.DB

	unit := &Unit{
		logger: logger,
		db:     db,
	}

	b := presets.New()

	mb := unit.Configure(b)

	b.Build()

	t.Run("Valid unit", func(t *testing.T) {
		validUnit := &models.Unit{
			Title: "Test Unit",
		}

		ctx := &web.EventContext{}

		validateFunc := mb.Editing().Validator
		err := validateFunc(validUnit, ctx)

		assert.False(t, err.HaveErrors())
	})

	t.Run("Invalid unit - empty name", func(t *testing.T) {
		invalidUnit := &models.Unit{
			Title: "",
		}

		ctx := &web.EventContext{}

		validateFunc := mb.Editing().Validator
		err := validateFunc(invalidUnit, ctx)

		require.True(t, err.HaveErrors())

		fieldErrors := err.FieldErrors()
		require.NotEmpty(t, fieldErrors)

		nameErrors := err.GetFieldErrors("Title")
		require.Len(t, nameErrors, 1)
		assert.Equal(t, "Name is required", nameErrors[0])
	})
}
