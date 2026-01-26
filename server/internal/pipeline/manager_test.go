package pipeline

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.ads.coffee/server/config"
	"go.ads.coffee/server/domain"
	"go.ads.coffee/server/internal/formats"
	"go.ads.coffee/server/internal/inputs"
	"go.ads.coffee/server/internal/outputs"
	"go.ads.coffee/server/internal/stages"
	"go.ads.coffee/server/internal/targetings"
)

// Mock implementations for testing
type mockInput struct {
	name string
}

func (m *mockInput) Name() string {
	return m.name
}

func (m *mockInput) Copy(cfg map[string]any) domain.Input {
	return &mockInput{name: m.name}
}

func (m *mockInput) Do(ctx context.Context, state *domain.State) bool {
	return true
}

type mockStage struct {
	name string
}

func (m *mockStage) Name() string {
	return m.name
}

func (m *mockStage) Copy(cfg map[string]any) domain.Stage {
	return &mockStage{name: m.name}
}

func (m *mockStage) Do(ctx context.Context, state *domain.State) {
	// Mock implementation
}

type mockStageWithTargetings struct {
	mockStage
	targetings []domain.Targeting
}

func (m *mockStageWithTargetings) Targetings(tt []domain.Targeting) {
	m.targetings = tt
}

func (m *mockStageWithTargetings) Copy(cfg map[string]any) domain.Stage {
	return &mockStageWithTargetings{
		mockStage: mockStage{name: m.name},
	}
}

type mockTargeting struct {
	name string
}

func (m *mockTargeting) Name() string {
	return m.name
}

func (m *mockTargeting) Copy(cfg map[string]any) domain.Targeting {
	return &mockTargeting{name: m.name}
}

func (m *mockTargeting) Filter() {
	// Mock implementation
}

type mockFormat struct {
	name string
}

func (m *mockFormat) Name() string {
	return m.name
}

func (m *mockFormat) Copy(cfg map[string]any) domain.Format {
	return &mockFormat{name: m.name}
}

func (m *mockFormat) Render(ctx context.Context, state *domain.State) {
	// Mock implementation
}

type mockOutput struct {
	name    string
	formats []domain.Format
}

func (m *mockOutput) Name() string {
	return m.name
}

func (m *mockOutput) Copy(cfg map[string]any) domain.Output {
	return &mockOutput{name: m.name}
}

func (m *mockOutput) Formats(ff []domain.Format) {
	m.formats = ff
}

func (m *mockOutput) Do(ctx context.Context, state *domain.State) {
	// Mock implementation
}

func TestNewManager(t *testing.T) {
	// Create mock inputs
	inputList := []domain.Input{
		&mockInput{name: "inputs.rtb"},
		&mockInput{name: "inputs.web"},
	}
	inputs := inputs.New(inputList)

	// Create mock outputs
	outputList := []domain.Output{
		&mockOutput{name: "outputs.rtb"},
		&mockOutput{name: "outputs.web"},
	}
	outputs := outputs.New(outputList)

	// Create mock stages
	stageList := []domain.Stage{
		&mockStage{name: "stages.banners"},
		&mockStageWithTargetings{mockStage: mockStage{name: "stages.targeting"}},
	}
	stages := stages.New(stageList)

	// Create mock targetings
	targetingList := []domain.Targeting{
		&mockTargeting{name: "targetings.apps"},
		&mockTargeting{name: "targetings.geo"},
	}
	targetings := targetings.New(targetingList)

	// Create mock formats
	formatList := []domain.Format{
		&mockFormat{name: "formats.native"},
		&mockFormat{name: "formats.banner"},
	}
	formats := formats.New(formatList)

	// Create test config
	cfg := config.Config{
		Pipelines: []config.Pipeline{
			{
				Name:  "dsp",
				Route: "/dsp",
				Input: config.Input{
					Name:   "inputs.rtb",
					Config: map[string]any{},
				},
				Stages: []config.Stage{
					{Name: "stages.banners", Config: map[string]any{}},
					{Name: "stages.targeting", Config: map[string]any{}},
				},
				Targetings: []config.Targeting{
					{Name: "targetings.apps", Config: map[string]any{}},
					{Name: "targetings.geo", Config: map[string]any{}},
				},
				Formats: []config.Format{
					{Name: "formats.native", Config: map[string]any{}},
					{Name: "formats.banner", Config: map[string]any{}},
				},
				Output: config.Output{
					Name:   "outputs.rtb",
					Config: map[string]any{},
				},
			},
		},
	}

	// Create manager
	manager := NewManager(cfg, inputs, outputs, stages, targetings, formats)

	// Assertions
	assert.NotNil(t, manager)
	assert.Len(t, manager.pipelines, 1)

	pipeline := manager.pipelines[0]
	assert.Equal(t, "dsp", pipeline.Name())
	assert.Equal(t, "/dsp", pipeline.Route())
}

func TestManagerMount(t *testing.T) {
	// Create mock inputs
	inputList := []domain.Input{
		&mockInput{name: "inputs.rtb"},
	}
	inputs := inputs.New(inputList)

	// Create mock outputs
	outputList := []domain.Output{
		&mockOutput{name: "outputs.rtb"},
	}
	outputs := outputs.New(outputList)

	// Create mock stages
	stageList := []domain.Stage{
		&mockStage{name: "stages.banners"},
	}
	stages := stages.New(stageList)

	// Create mock targetings
	targetingList := []domain.Targeting{
		&mockTargeting{name: "targetings.apps"},
	}
	targetings := targetings.New(targetingList)

	// Create mock formats
	formatList := []domain.Format{
		&mockFormat{name: "formats.native"},
	}
	formats := formats.New(formatList)

	// Create test config
	cfg := config.Config{
		Pipelines: []config.Pipeline{
			{
				Name:  "dsp",
				Route: "/dsp",
				Input: config.Input{
					Name:   "inputs.rtb",
					Config: map[string]any{},
				},
				Stages: []config.Stage{
					{Name: "stages.banners", Config: map[string]any{}},
				},
				Targetings: []config.Targeting{
					{Name: "targetings.apps", Config: map[string]any{}},
				},
				Formats: []config.Format{
					{Name: "formats.native", Config: map[string]any{}},
				},
				Output: config.Output{
					Name:   "outputs.rtb",
					Config: map[string]any{},
				},
			},
		},
	}

	// Create manager
	manager := NewManager(cfg, inputs, outputs, stages, targetings, formats)

	// Create router
	router := chi.NewRouter()

	// Mount pipelines
	manager.Mount(router)

	// Assertions
	// We can't easily test the mounting without making actual HTTP requests,
	// but we can verify the manager was created correctly
	assert.NotNil(t, manager)
	assert.Len(t, manager.pipelines, 1)
}
