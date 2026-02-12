package formats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestNative_Name(t *testing.T) {
	// Arrange
	native := &Native{}

	// Act
	name := native.Name()

	// Assert
	assert.Equal(t, "native", name)
}

func TestNative_Copy(t *testing.T) {
	// Arrange
	native := &Native{}
	cfg := map[string]any{
		"base": "https://example.com",
	}

	// Act
	copied := native.Copy(cfg)

	// Assert
	assert.NotNil(t, copied)
	assert.IsType(t, &Native{}, copied)

	// Check that the copied native has the correct base
	copiedNative := copied.(*Native)
	assert.Equal(t, "https://example.com", copiedNative.base)
}

func TestNative_Copy_WithEmptyConfig(t *testing.T) {
	// Arrange
	native := &Native{}
	cfg := map[string]any{}

	// Act
	copied := native.Copy(cfg)

	// Assert
	assert.NotNil(t, copied)
	assert.IsType(t, &Native{}, copied)

	// Check that the copied native has an empty base
	copiedNative := copied.(*Native)
	assert.Equal(t, "", copiedNative.base)
}

func TestNative_Render_WithEmptyWinners(t *testing.T) {
	// Arrange
	native := &Native{}
	ctx := context.Background()
	state := &plugins.State{
		Winners: []ads.Banner{},
	}

	// Act
	result, err := native.Render(ctx, state)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that result is an empty slice
	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 0)
}

func TestNative_Render_WithOneWinner(t *testing.T) {
	// Arrange
	native := &Native{}
	ctx := context.Background()

	// Create a test image
	image, _ := ads.NewImage(`{"url":"/test/image.jpg"}`)

	state := &plugins.State{
		Winners: []ads.Banner{
			{
				Title:       "Test Native Ad",
				Description: "This is a test native advertisement",
				Target:      "https://example.com/click",
				Image:       image,
			},
		},
	}

	// Act
	result, err := native.Render(ctx, state)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that result is a slice with one item
	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 1)

	// Check the item properties
	item := items[0]
	assert.Equal(t, "Test Native Ad", item.Title)
	assert.Equal(t, "This is a test native advertisement", item.Description)
	assert.Equal(t, "https://example.com/click", item.Target)
	assert.Equal(t, "/test/image.jpg", item.Image) // base URL is prepended

	// Check that trackers are empty arrays
	assert.Empty(t, item.Impressions)
	assert.Empty(t, item.Clicks)
}

func TestNative_Render_WithMultipleWinners(t *testing.T) {
	// Arrange
	native := &Native{}
	ctx := context.Background()

	// Create test images
	image1, _ := ads.NewImage(`{"url":"/test/image1.jpg"}`)
	image2, _ := ads.NewImage(`{"url":"/test/image2.jpg"}`)

	state := &plugins.State{
		Winners: []ads.Banner{
			{
				Title:       "First Native Ad",
				Description: "This is the first test native advertisement",
				Target:      "https://example.com/click1",
				Image:       image1,
			},
			{
				Title:       "Second Native Ad",
				Description: "This is the second test native advertisement",
				Target:      "https://example.com/click2",
				Image:       image2,
			},
		},
	}

	// Act
	result, err := native.Render(ctx, state)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that result is a slice with two items
	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 2)

	// Check the first item properties
	item1 := items[0]
	assert.Equal(t, "First Native Ad", item1.Title)
	assert.Equal(t, "This is the first test native advertisement", item1.Description)
	assert.Equal(t, "https://example.com/click1", item1.Target)
	assert.Equal(t, "/test/image1.jpg", item1.Image)

	// Check that trackers are empty arrays
	assert.Empty(t, item1.Impressions)
	assert.Empty(t, item1.Clicks)

	// Check the second item properties
	item2 := items[1]
	assert.Equal(t, "Second Native Ad", item2.Title)
	assert.Equal(t, "This is the second test native advertisement", item2.Description)
	assert.Equal(t, "https://example.com/click2", item2.Target)
	assert.Equal(t, "/test/image2.jpg", item2.Image)

	// Check that trackers are empty arrays
	assert.Empty(t, item2.Impressions)
	assert.Empty(t, item2.Clicks)
}

func TestNative_Render_WithEmptyFields(t *testing.T) {
	// Arrange
	native := &Native{}
	ctx := context.Background()

	// Create an empty image
	image, _ := ads.NewImage(`{"url":""}`)

	state := &plugins.State{
		Winners: []ads.Banner{
			{
				Title:       "",
				Description: "",
				Target:      "",
				Image:       image,
			},
		},
	}

	// Act
	result, err := native.Render(ctx, state)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that result is a slice with one item
	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 1)

	// Check the item properties are empty
	item := items[0]
	assert.Equal(t, "", item.Title)
	assert.Equal(t, "", item.Description)
	assert.Equal(t, "", item.Target)
	assert.Equal(t, "", item.Image) // Just the protocol prefix

	// Check that trackers are empty arrays
	assert.Empty(t, item.Impressions)
	assert.Empty(t, item.Clicks)
}
