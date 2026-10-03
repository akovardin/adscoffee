package stats

import (
	"testing"

	"github.com/qor5/admin/v3/presets"
	"github.com/stretchr/testify/assert"

	"go.ads.coffee/platform/admin/internal/config"
)

func TestNew(t *testing.T) {
	stats := New(config.Grafana{})
	assert.NotNil(t, stats)
}

func TestConfigure(t *testing.T) {
	stats := New(config.Grafana{})
	b := presets.New()

	assert.NotPanics(t, func() {
		stats.Configure(b)
	})
	assert.NotNil(t, b.Model(&Dashboard{}))
}