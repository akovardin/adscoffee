package stages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type mockStage struct {
	name string
}

func (m *mockStage) Name() string {
	return m.name
}

func (m *mockStage) Do(ctx context.Context, state *plugins.State) error {
	return nil
}

func mockProvider(name string) plugins.Provider[plugins.Stage] {
	return plugins.Provider[plugins.Stage]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Stage, error) {
			return &mockStage{name: name}, nil
		},
	}
}

func TestNew(t *testing.T) {
	stages := New([]plugins.Provider[plugins.Stage]{
		mockProvider("stages.banners"),
		mockProvider("stages.rotation"),
	})

	assert.NotNil(t, stages)
}

func TestStages_Get(t *testing.T) {
	stages := New([]plugins.Provider[plugins.Stage]{
		mockProvider("stages.banners"),
		mockProvider("stages.rotation"),
	})

	cfg := map[string]any{"key": "value"}
	stage, err := stages.Get("stages.banners", cfg)
	require.NoError(t, err)

	assert.NotNil(t, stage)
	assert.Equal(t, "stages.banners", stage.Name())
}

func TestNewWithEmptySlice(t *testing.T) {
	stages := New([]plugins.Provider[plugins.Stage]{})

	assert.NotNil(t, stages)
}

func TestStages_GetNonExistentStage(t *testing.T) {
	stages := New([]plugins.Provider[plugins.Stage]{
		mockProvider("stages.banners"),
	})

	_, err := stages.Get("non-existent", map[string]any{})
	assert.Error(t, err)
}
