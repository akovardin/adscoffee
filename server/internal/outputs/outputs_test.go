package outputs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/domain"
)

type mockOutput struct {
	name string
}

func (m *mockOutput) Name() string {
	return m.name
}

func (m *mockOutput) Copy(cfg map[string]any) domain.Output {
	return &mockOutput{name: m.name}
}

func (m *mockOutput) Formats(ff []domain.Format) {
}

func (m *mockOutput) Do(ctx context.Context, state *domain.State) {
}

func TestNew(t *testing.T) {
	outputList := []domain.Output{
		&mockOutput{name: "outputs.rtb"},
		&mockOutput{name: "outputs.web"},
	}

	outputs := New(outputList)

	assert.NotNil(t, outputs)
	assert.Len(t, outputs.list, 2)
	assert.Contains(t, outputs.list, "outputs.rtb")
	assert.Contains(t, outputs.list, "outputs.web")
}

func TestOutputs_Get(t *testing.T) {
	outputList := []domain.Output{
		&mockOutput{name: "outputs.rtb"},
		&mockOutput{name: "outputs.web"},
	}

	outputs := New(outputList)

	cfg := map[string]any{"key": "value"}
	output := outputs.Get("outputs.rtb", cfg)

	assert.NotNil(t, output)
	assert.Equal(t, "outputs.rtb", output.Name())
}

func TestNewWithEmptySlice(t *testing.T) {
	outputs := New([]domain.Output{})

	assert.NotNil(t, outputs)
	assert.Len(t, outputs.list, 0)
}

func TestOutput_GetNonExistentOutput(t *testing.T) {
	outputList := []domain.Output{
		&mockOutput{name: "outputs.rtb"},
	}

	outputs := New(outputList)

	assert.Panics(t, func() {
		cfg := map[string]any{}
		outputs.Get("non-existent", cfg)
	})
}
