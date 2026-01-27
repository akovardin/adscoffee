package formats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

// Mock implementation of domain.Format for testing
type mockFormat struct {
	name string
}

func (m *mockFormat) Name() string {
	return m.name
}

func (m *mockFormat) Copy(cfg map[string]any) plugins.Format {
	return &mockFormat{name: m.name}
}

func (m *mockFormat) Render(ctx context.Context, state *plugins.State) {
	// Mock implementation
}

func TestNew(t *testing.T) {
	// Create mock formats
	formatList := []plugins.Format{
		&mockFormat{name: "formats.native"},
		&mockFormat{name: "formats.banner"},
	}

	// Create Formats instance
	formats := New(formatList)

	// Assertions
	assert.NotNil(t, formats)
	assert.Len(t, formats.list, 2)
	assert.Contains(t, formats.list, "formats.native")
	assert.Contains(t, formats.list, "formats.banner")
}

func TestFormats_Get(t *testing.T) {
	// Create mock formats
	formatList := []plugins.Format{
		&mockFormat{name: "formats.native"},
		&mockFormat{name: "formats.banner"},
	}

	// Create Formats instance
	formats := New(formatList)

	// Test getting a format
	cfg := map[string]any{"key": "value"}
	format := formats.Get("formats.native", cfg)

	// Assertions
	assert.NotNil(t, format)
	assert.Equal(t, "formats.native", format.Name())
}
