//nolint:errcheck
package pipeline

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/inputs"
	"go.ads.coffee/platform/server/internal/outputs"
	"go.ads.coffee/platform/server/internal/stages"
	"go.ads.coffee/platform/server/internal/targetings"
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

type mockStageWithTargetings struct {
	mockStage
	targetings []plugins.Targeting
}

func (m *mockStageWithTargetings) Targetings(tt []plugins.Targeting) {
	m.targetings = tt
}

func (m *mockStageWithTargetings) Copy(cfg map[string]any) plugins.Stage {
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

func (m *mockTargeting) Copy(cfg map[string]any) plugins.Targeting {
	return &mockTargeting{name: m.name}
}

func (m *mockTargeting) Filter() {
}

type mockOutput struct {
	name    string
	formats []plugins.Format
}

func (m *mockOutput) Name() string {
	return m.name
}

func (m *mockOutput) Copy(cfg map[string]any) plugins.Output {
	return &mockOutput{name: m.name}
}

func (m *mockOutput) Formats(ff []plugins.Format) {
	m.formats = ff
}

func (m *mockOutput) Do(ctx context.Context, state *plugins.State) error {
	return nil
}

func TestNewManager(t *testing.T) {
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
		&mockInput{name: "inputs.web"},
	}
	inputs := inputs.New(inputList)

	outputList := []plugins.Output{
		&mockOutput{name: "outputs.rtb"},
		&mockOutput{name: "outputs.web"},
	}
	outputs := outputs.New(outputList)

	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
		&mockStageWithTargetings{mockStage: mockStage{name: "stages.targeting"}},
	}
	stages := stages.New(stageList)

	targetingList := []plugins.Targeting{
		&mockTargeting{name: "targetings.apps"},
		&mockTargeting{name: "targetings.geo"},
	}
	targetings := targetings.New(targetingList)

	cfg := []Config{
		{
			Name:  "dsp",
			Route: "/dsp",
			Input: Input{
				Name:   "inputs.rtb",
				Config: map[string]any{},
			},
			Stages: []Stage{
				{Name: "stages.banners", Config: map[string]any{}},
				{Name: "stages.targeting", Config: map[string]any{}},
			},
			Targetings: []Targeting{
				{Name: "targetings.apps", Config: map[string]any{}},
				{Name: "targetings.geo", Config: map[string]any{}},
			},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager := NewManager(cfg, inputs, outputs, stages, targetings)

	assert.NotNil(t, manager)
	assert.Len(t, manager.pipelines, 1)

	pipeline := manager.pipelines[0]
	assert.Equal(t, "dsp", pipeline.Name())
	assert.Equal(t, "/dsp", pipeline.Route())
}

func TestManager_Mount(t *testing.T) {
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
	}
	inputs := inputs.New(inputList)

	outputList := []plugins.Output{
		&mockOutput{name: "outputs.rtb"},
	}
	outputs := outputs.New(outputList)

	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
	}
	stages := stages.New(stageList)

	targetingList := []plugins.Targeting{
		&mockTargeting{name: "targetings.apps"},
	}
	targetings := targetings.New(targetingList)

	cfg := []Config{
		{
			Name:  "dsp",
			Route: "/dsp",
			Input: Input{
				Name:   "inputs.rtb",
				Config: map[string]any{},
			},
			Stages: []Stage{
				{Name: "stages.banners", Config: map[string]any{}},
			},
			Targetings: []Targeting{
				{Name: "targetings.apps", Config: map[string]any{}},
			},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager := NewManager(cfg, inputs, outputs, stages, targetings)

	router := chi.NewRouter()

	manager.Mount(router)

	assert.NotNil(t, manager)
	assert.Len(t, manager.pipelines, 1)
}

func TestManager_MountHandlers(t *testing.T) {
	// Create mock components
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
	}
	inputs := inputs.New(inputList)

	outputList := []plugins.Output{
		&mockOutput{name: "outputs.rtb"},
	}
	outputs := outputs.New(outputList)

	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
	}
	stages := stages.New(stageList)

	targetingList := []plugins.Targeting{
		&mockTargeting{name: "targetings.apps"},
	}
	targetings := targetings.New(targetingList)

	// Create config with multiple pipelines
	cfg := []Config{
		{
			Name:  "dsp",
			Route: "/dsp",
			Input: Input{
				Name:   "inputs.rtb",
				Config: map[string]any{},
			},
			Stages: []Stage{
				{Name: "stages.banners", Config: map[string]any{}},
			},
			Targetings: []Targeting{
				{Name: "targetings.apps", Config: map[string]any{}},
			},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
		{
			Name:  "web",
			Route: "/web",
			Input: Input{
				Name:   "inputs.rtb",
				Config: map[string]any{},
			},
			Stages: []Stage{
				{Name: "stages.banners", Config: map[string]any{}},
			},
			Targetings: []Targeting{
				{Name: "targetings.apps", Config: map[string]any{}},
			},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager := NewManager(cfg, inputs, outputs, stages, targetings)

	router := chi.NewRouter()

	// Mount the pipelines
	manager.Mount(router)

	// Test that routes are registered by making requests
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Test first route
	resp, err := http.Get(ts.URL + "/dsp")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Test second route
	resp, err = http.Get(ts.URL + "/web")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestManager_MountWithPipelineDo(t *testing.T) {
	// Create mock components
	inputList := []plugins.Input{
		&mockInput{name: "inputs.rtb"},
	}
	inputs := inputs.New(inputList)

	outputList := []plugins.Output{
		&mockOutput{name: "outputs.rtb"},
	}
	outputs := outputs.New(outputList)

	stageList := []plugins.Stage{
		&mockStage{name: "stages.banners"},
	}
	stages := stages.New(stageList)

	targetingList := []plugins.Targeting{
		&mockTargeting{name: "targetings.apps"},
	}
	targetings := targetings.New(targetingList)

	// Create config with a single pipeline
	cfg := []Config{
		{
			Name:  "test",
			Route: "/test",
			Input: Input{
				Name:   "inputs.rtb",
				Config: map[string]any{},
			},
			Stages: []Stage{
				{Name: "stages.banners", Config: map[string]any{}},
			},
			Targetings: []Targeting{
				{Name: "targetings.apps", Config: map[string]any{}},
			},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager := NewManager(cfg, inputs, outputs, stages, targetings)

	router := chi.NewRouter()

	// Mount the pipelines
	manager.Mount(router)

	// Test that the handler is properly mounted
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Make a request to trigger the handler
	resp, err := http.Get(ts.URL + "/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
