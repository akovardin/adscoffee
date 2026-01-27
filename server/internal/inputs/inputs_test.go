package inputs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type mockInput struct {
	name string
}

func (m *mockInput) Name() string {
	return m.name
}

func (m *mockInput) Copy(cfg map[string]any) plugins.Input {
	return &mockInput{name: m.name}
}

func (m *mockInput) Do(ctx context.Context, state *plugins.State) bool {
	return true
}

func TestNew(t *testing.T) {
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
		&mockInput{name: "inputs.web"},
	}

	inputs := New(inputList)

	assert.NotNil(t, inputs)
	assert.Len(t, inputs.plugins, 2)
	assert.Contains(t, inputs.plugins, "inputs.rtb")
	assert.Contains(t, inputs.plugins, "inputs.web")
}

func TestInputs_Get(t *testing.T) {
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
		&mockInput{name: "inputs.web"},
	}

	inputs := New(inputList)

	cfg := map[string]any{"key": "value"}
	input := inputs.Get("inputs.rtb", cfg)

	assert.NotNil(t, input)
	assert.Equal(t, "inputs.rtb", input.Name())
}

func TestNewWithEmptySlice(t *testing.T) {
	inputs := New([]plugins.Input{})

	assert.NotNil(t, inputs)
	assert.Len(t, inputs.plugins, 0)
}

func TestInput_GetNonExistentInput(t *testing.T) {
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
	}

	inputs := New(inputList)

	assert.Panics(t, func() {
		cfg := map[string]any{}
		inputs.Get("non-existent", cfg)
	})
}
