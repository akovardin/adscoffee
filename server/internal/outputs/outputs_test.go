package outputs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type mockOutput struct {
	name string
}

func (m *mockOutput) Name() string {
	return m.name
}

func (m *mockOutput) Do(ctx context.Context, state *plugins.State) error {
	return nil
}

func mockProvider(name string) plugins.Provider[plugins.Output] {
	return plugins.Provider[plugins.Output]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Output, error) {
			return &mockOutput{name: name}, nil
		},
	}
}

func TestNew(t *testing.T) {
	outputs := New([]plugins.Provider[plugins.Output]{
		mockProvider("outputs.rtb"),
		mockProvider("outputs.web"),
	})

	assert.NotNil(t, outputs)
}

func TestOutputs_Get(t *testing.T) {
	outputs := New([]plugins.Provider[plugins.Output]{
		mockProvider("outputs.rtb"),
		mockProvider("outputs.web"),
	})

	cfg := map[string]any{"key": "value"}
	output, err := outputs.Get("outputs.rtb", cfg)
	require.NoError(t, err)

	assert.NotNil(t, output)
	assert.Equal(t, "outputs.rtb", output.Name())
}

func TestNewWithEmptySlice(t *testing.T) {
	outputs := New([]plugins.Provider[plugins.Output]{})

	assert.NotNil(t, outputs)
}

func TestOutput_GetNonExistentOutput(t *testing.T) {
	outputs := New([]plugins.Provider[plugins.Output]{
		mockProvider("outputs.rtb"),
	})

	_, err := outputs.Get("non-existent", map[string]any{})
	assert.Error(t, err)
}
