package formats

import (
	"context"
	"testing"

	"github.com/qor5/admin/v3/media/media_library"
	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestNative_Name(t *testing.T) {
	native := &Native{}

	name := native.Name()

	assert.Equal(t, "native", name)
}

func TestNative_Copy(t *testing.T) {
	native := &Native{}
	cfg := map[string]any{
		"base": "https://example.com",
	}

	copied := native.Copy(cfg)

	assert.NotNil(t, copied)
	assert.IsType(t, &Native{}, copied)

	copiedNative := copied.(*Native)
	assert.Equal(t, "https://example.com", copiedNative.base)
}

func TestNative_Copy_WithEmptyConfig(t *testing.T) {
	native := &Native{}
	cfg := map[string]any{}

	copied := native.Copy(cfg)

	assert.NotNil(t, copied)
	assert.IsType(t, &Native{}, copied)

	copiedNative := copied.(*Native)
	assert.Equal(t, "", copiedNative.base)
}

func TestNative_Render_WithEmptyWinners(t *testing.T) {
	native := &Native{}
	ctx := context.Background()
	state := &plugins.State{
		Winners: []ads.Banner{},
	}

	result, err := native.Render(ctx, state)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 0)
}

func TestNative_Render_WithOneWinner(t *testing.T) {
	native := &Native{}
	ctx := context.Background()

	image := media_library.MediaBox{Url: "/test/image.jpg"}

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

	result, err := native.Render(ctx, state)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "Test Native Ad", item.Title)
	assert.Equal(t, "This is a test native advertisement", item.Description)
	assert.Equal(t, "https://example.com/click", item.Target)
	assert.Equal(t, "https:/test/image.250x250.jpg", item.Image) // base URL is prepended

	assert.NotEmpty(t, item.Impressions)
	assert.NotEmpty(t, item.Clicks)
}

func TestNative_Render_WithMultipleWinners(t *testing.T) {
	native := &Native{}
	ctx := context.Background()

	image1 := media_library.MediaBox{
		Url: "/test/image1.jpg",
	}

	image2 := media_library.MediaBox{
		Url: "/test/image2.jpg",
	}

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

	result, err := native.Render(ctx, state)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 2)

	item1 := items[0]
	assert.Equal(t, "First Native Ad", item1.Title)
	assert.Equal(t, "This is the first test native advertisement", item1.Description)
	assert.Equal(t, "https://example.com/click1", item1.Target)
	assert.Equal(t, "https:/test/image1.250x250.jpg", item1.Image)

	assert.NotEmpty(t, item1.Impressions)
	assert.NotEmpty(t, item1.Clicks)

	item2 := items[1]
	assert.Equal(t, "Second Native Ad", item2.Title)
	assert.Equal(t, "This is the second test native advertisement", item2.Description)
	assert.Equal(t, "https://example.com/click2", item2.Target)
	assert.Equal(t, "https:/test/image2.250x250.jpg", item2.Image)

	assert.NotEmpty(t, item2.Impressions)
	assert.NotEmpty(t, item2.Clicks)
}

func TestNative_Render_WithEmptyFields(t *testing.T) {
	native := &Native{}
	ctx := context.Background()

	image := media_library.MediaBox{}

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

	result, err := native.Render(ctx, state)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	items, ok := result.([]NativeResponse)
	assert.True(t, ok)
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "", item.Title)
	assert.Equal(t, "", item.Description)
	assert.Equal(t, "", item.Target)
	assert.Equal(t, "", item.Image) // Just the protocol prefix

	assert.NotEmpty(t, item.Impressions)
	assert.NotEmpty(t, item.Clicks)
}
