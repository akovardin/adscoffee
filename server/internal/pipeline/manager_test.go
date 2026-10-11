//nolint:errcheck
package pipeline

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/server/internal/domain/ads"
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

type mockStage struct {
	name string
}

func (m *mockStage) Name() string {
	return m.name
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

type mockStageWithLimits struct {
	mockStage
	limits []plugins.Limit
}

func (m *mockStageWithLimits) Limits(ll []plugins.Limit) {
	m.limits = ll
}

type mockTargeting struct {
	name string
}

func (m *mockTargeting) Name() string {
	return m.name
}

func (m *mockTargeting) Filter(candidates []ads.Banner, state *plugins.State) []ads.Banner {
	return candidates
}

type mockLimit struct {
	name string
}

func (m *mockLimit) Name() string {
	return m.name
}

func (m *mockLimit) Apply(candidates []ads.Banner, state *plugins.State) []ads.Banner {
	return candidates
}

type mockOutput struct {
	name string
}

func (m *mockOutput) Name() string {
	return m.name
}

func (m *mockOutput) Do(ctx context.Context, state *plugins.State) error {
	return nil
}

func inputProvider(name string) plugins.Provider[plugins.Input] {
	return plugins.Provider[plugins.Input]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Input, error) {
			return &mockInput{name: name}, nil
		},
	}
}

func outputProvider(name string) plugins.Provider[plugins.Output] {
	return plugins.Provider[plugins.Output]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Output, error) {
			return &mockOutput{name: name}, nil
		},
	}
}

func stageProvider(name string, build func() plugins.Stage) plugins.Provider[plugins.Stage] {
	return plugins.Provider[plugins.Stage]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Stage, error) {
			return build(), nil
		},
	}
}

func targetingProvider(name string) plugins.Provider[plugins.Targeting] {
	return plugins.Provider[plugins.Targeting]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Targeting, error) {
			return &mockTargeting{name: name}, nil
		},
	}
}

func limitProvider(name string) plugins.Provider[plugins.Limit] {
	return plugins.Provider[plugins.Limit]{
		Name: name,
		Factory: func(cfg map[string]any) (plugins.Limit, error) {
			return &mockLimit{name: name}, nil
		},
	}
}

func TestNewManager(t *testing.T) {
	in := newInputs([]plugins.Provider[plugins.Input]{
		inputProvider("inputs.rtb"),
		inputProvider("inputs.web"),
	})

	out := newOutputs([]plugins.Provider[plugins.Output]{
		outputProvider("outputs.rtb"),
		outputProvider("outputs.web"),
	})

	limitStage := &mockStageWithLimits{mockStage: mockStage{name: "stages.limit"}}

	st := newStages([]plugins.Provider[plugins.Stage]{
		stageProvider("stages.banners", func() plugins.Stage { return &mockStage{name: "stages.banners"} }),
		stageProvider("stages.targeting", func() plugins.Stage { return &mockStageWithTargetings{mockStage: mockStage{name: "stages.targeting"}} }),
		{
			Name: "stages.limit",
			Factory: func(cfg map[string]any) (plugins.Stage, error) {
				return limitStage, nil
			},
		},
	})

	tg := newTargetings([]plugins.Provider[plugins.Targeting]{
		targetingProvider("targetings.apps"),
		targetingProvider("targetings.geo"),
	})

	lm := newLimits([]plugins.Provider[plugins.Limit]{
		limitProvider("limits.budget"),
	})

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
				{Name: "stages.limit", Config: map[string]any{}},
			},
			Targetings: []Targeting{
				{Name: "targetings.apps", Config: map[string]any{}},
				{Name: "targetings.geo", Config: map[string]any{}},
			},
			Limits: []Limit{
				{Name: "limits.budget", Config: map[string]any{}},
			},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager, err := NewManager(cfg, in, out, st, tg, lm)
	require.NoError(t, err)

	assert.NotNil(t, manager)
	assert.Len(t, manager.pipelines, 1)

	pipeline := manager.pipelines[0]
	assert.Equal(t, "dsp", pipeline.Name())
	assert.Equal(t, "/dsp", pipeline.Route())

	// лимиты из конфига проброшены в стейдж, реализующий WithLimits
	assert.Len(t, limitStage.limits, 1)
}

func TestNewManager_UnknownPlugin(t *testing.T) {
	cfg := []Config{
		{
			Name:  "dsp",
			Route: "/dsp",
			Input: Input{Name: "inputs.missing", Config: map[string]any{}},
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	_, err := NewManager(
		cfg,
		newInputs(nil),
		newOutputs(nil),
		newStages(nil),
		newTargetings(nil),
		newLimits(nil),
	)
	assert.Error(t, err)
}

func TestManager_Mount(t *testing.T) {
	in := newInputs([]plugins.Provider[plugins.Input]{
		inputProvider("inputs.rtb"),
	})

	out := newOutputs([]plugins.Provider[plugins.Output]{
		outputProvider("outputs.rtb"),
	})

	st := newStages([]plugins.Provider[plugins.Stage]{
		stageProvider("stages.banners", func() plugins.Stage { return &mockStage{name: "stages.banners"} }),
	})

	tg := newTargetings([]plugins.Provider[plugins.Targeting]{
		targetingProvider("targetings.apps"),
	})

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

	manager, err := NewManager(cfg, in, out, st, tg, newLimits(nil))
	require.NoError(t, err)

	router := chi.NewRouter()

	manager.Mount(router)

	assert.Len(t, manager.pipelines, 1)
}

func TestManager_MountHandlers(t *testing.T) {
	in := newInputs([]plugins.Provider[plugins.Input]{
		inputProvider("inputs.rtb"),
	})

	out := newOutputs([]plugins.Provider[plugins.Output]{
		outputProvider("outputs.rtb"),
	})

	st := newStages([]plugins.Provider[plugins.Stage]{
		stageProvider("stages.banners", func() plugins.Stage { return &mockStage{name: "stages.banners"} }),
	})

	tg := newTargetings(nil)

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
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager, err := NewManager(cfg, in, out, st, tg, newLimits(nil))
	require.NoError(t, err)

	router := chi.NewRouter()

	manager.Mount(router)

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/dsp")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(ts.URL + "/web")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestManager_MountWithPipelineDo(t *testing.T) {
	in := newInputs([]plugins.Provider[plugins.Input]{
		inputProvider("inputs.rtb"),
	})

	out := newOutputs([]plugins.Provider[plugins.Output]{
		outputProvider("outputs.rtb"),
	})

	st := newStages([]plugins.Provider[plugins.Stage]{
		stageProvider("stages.banners", func() plugins.Stage { return &mockStage{name: "stages.banners"} }),
	})

	tg := newTargetings(nil)

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
			Output: Output{
				Name:   "outputs.rtb",
				Config: map[string]any{},
			},
		},
	}

	manager, err := NewManager(cfg, in, out, st, tg, newLimits(nil))
	require.NoError(t, err)

	router := chi.NewRouter()

	manager.Mount(router)

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
