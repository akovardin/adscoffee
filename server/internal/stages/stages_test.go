package stages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type mockStage struct {
	name string
}

func (m *mockStage) Name() string {
	return m.name
}

func (m *mockStage) Copy(cfg map[string]any) plugins.Stage {
	return &mockStage{name: m.name}
}

func (m *mockStage) Do(ctx context.Context, state *plugins.State) error {
	return nil
}

func TestNew(t *testing.T) {
	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
		&mockStage{name: "stages.rotation"},
	}

	stages := New(stageList)

	assert.NotNil(t, stages)
	assert.Len(t, stages.list, 2)
	assert.Contains(t, stages.list, "stages.banners")
	assert.Contains(t, stages.list, "stages.rotation")
}

func TestStages_Get(t *testing.T) {
	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
		&mockStage{name: "stages.rotation"},
	}

	stages := New(stageList)

	cfg := map[string]any{"key": "value"}
	stage := stages.Get("stages.banners", cfg)

	assert.NotNil(t, stage)
	assert.Equal(t, "stages.banners", stage.Name())
}

func TestNewWithEmptySlice(t *testing.T) {
	stages := New([]plugins.Stage{})

	assert.NotNil(t, stages)
	assert.Len(t, stages.list, 0)
}

func TestStages_GetNonExistentStage(t *testing.T) {
	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
	}

	stages := New(stageList)

	assert.Panics(t, func() {
		cfg := map[string]any{}
		stages.Get("non-existent", cfg)
	})
}
