package inputs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type mockInput struct {
	name string
}

func (m *mockInput) Name() string {
	return m.name
}

func (m *mockInput) Do(ctx context.Context, state *plugins.State) bool {
	return true
}

func mockProvider(name string) plugins.Provider[plugins.Input] {
	return plugins.Provider[plugins.Input]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Input, error) {
			return &mockInput{name: name}, nil
		},
	}
}

func TestNew(t *testing.T) {
	inputs := New([]plugins.Provider[plugins.Input]{
		mockProvider("inputs.rtb"),
		mockProvider("inputs.web"),
	})

	assert.NotNil(t, inputs)
}

func TestInputs_Get(t *testing.T) {
	inputs := New([]plugins.Provider[plugins.Input]{
		mockProvider("inputs.rtb"),
		mockProvider("inputs.web"),
	})

	cfg := map[string]any{"key": "value"}
	input, err := inputs.Get("inputs.rtb", cfg)
	require.NoError(t, err)

	assert.NotNil(t, input)
	assert.Equal(t, "inputs.rtb", input.Name())
}

func TestNewWithEmptySlice(t *testing.T) {
	inputs := New([]plugins.Provider[plugins.Input]{})

	assert.NotNil(t, inputs)
}

func TestInput_GetNonExistentInput(t *testing.T) {
	inputs := New([]plugins.Provider[plugins.Input]{
		mockProvider("inputs.rtb"),
	})

	_, err := inputs.Get("non-existent", map[string]any{})
	assert.Error(t, err)
}
