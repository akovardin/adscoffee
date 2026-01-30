package pipeline

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type testMockInput struct {
	name     string
	doCalled bool
}

func (m *testMockInput) Name() string {
	return m.name
}

func (m *testMockInput) Copy(cfg map[string]any) plugins.Input {
	return &testMockInput{name: m.name}
}

func (m *testMockInput) Do(ctx context.Context, state *plugins.State) bool {
	m.doCalled = true
	return true
}

type testMockStage struct {
	name     string
	doCalled bool
}

func (m *testMockStage) Name() string {
	return m.name
}

func (m *testMockStage) Copy(cfg map[string]any) plugins.Stage {
	return &testMockStage{name: m.name}
}

func (m *testMockStage) Do(ctx context.Context, state *plugins.State) error {
	m.doCalled = true

	return nil
}

type testMockOutput struct {
	name     string
	doCalled bool
}

func (m *testMockOutput) Name() string {
	return m.name
}

func (m *testMockOutput) Copy(cfg map[string]any) plugins.Output {
	return &testMockOutput{name: m.name}
}

func (m *testMockOutput) Formats(ff []plugins.Format) {
}

func (m *testMockOutput) Do(ctx context.Context, state *plugins.State) error {
	m.doCalled = true

	return nil
}

func TestPipeline_Do(t *testing.T) {
	mockInput := &testMockInput{name: "test-input"}
	mockStage1 := &testMockStage{name: "test-stage-1"}
	mockStage2 := &testMockStage{name: "test-stage-2"}
	mockOutput := &testMockOutput{name: "test-output"}

	pipeline := NewPipeline(
		"test-pipeline",
		"/test",
		mockInput,
		mockOutput,
		[]plugins.Stage{mockStage1, mockStage2},
	)

	ctx := context.Background()
	state := &plugins.State{
		Request:  &http.Request{},
		Response: nil,
	}

	err := pipeline.Do(ctx, state)

	assert.NoError(t, err)
	assert.True(t, mockInput.doCalled, "Input.Do should be called")
	assert.True(t, mockStage1.doCalled, "Stage 1 Do should be called")
	assert.True(t, mockStage2.doCalled, "Stage 2 Do should be called")
	assert.True(t, mockOutput.doCalled, "Output.Do should be called")
}

func TestPipeline_DoWithNoStages(t *testing.T) {
	mockInput := &testMockInput{name: "test-input"}
	mockOutput := &testMockOutput{name: "test-output"}

	pipeline := NewPipeline(
		"test-pipeline",
		"/test",
		mockInput,
		mockOutput,
		[]plugins.Stage{},
	)

	ctx := context.Background()
	state := &plugins.State{
		Request:  &http.Request{},
		Response: nil,
	}

	err := pipeline.Do(ctx, state)

	assert.NoError(t, err)
	assert.True(t, mockInput.doCalled, "Input.Do should be called")
	assert.True(t, mockOutput.doCalled, "Output.Do should be called")
}

func TestPipeline_NameAndRoute(t *testing.T) {
	pipeline := NewPipeline(
		"test-pipeline",
		"/test-route",
		&testMockInput{name: "test-input"},
		&testMockOutput{name: "test-output"},
		[]plugins.Stage{},
	)

	assert.Equal(t, "test-pipeline", pipeline.Name(), "Pipeline name should match")
	assert.Equal(t, "/test-route", pipeline.Route(), "Pipeline route should match")
}
