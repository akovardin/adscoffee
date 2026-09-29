package stats

import (
	"testing"

	"github.com/qor5/admin/v3/presets"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	stats := New()
	assert.NotNil(t, stats)
}

func TestConfigure(t *testing.T) {
	stats := New()
	b := presets.New()

	assert.NotPanics(t, func() {
		stats.Configure(b)
	})
	assert.NotNil(t, b.Model(&Dashboard{}))
}