package pipeline

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/domain"
)

type testMockInput struct {
	name     string
	doCalled bool
}

func (m *testMockInput) Name() string {
	return m.name
}

func (m *testMockInput) Copy(cfg map[string]any) domain.Input {
	return &testMockInput{name: m.name}
}

func (m *testMockInput) Do(ctx context.Context, state *domain.State) bool {
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

func (m *testMockStage) Copy(cfg map[string]any) domain.Stage {
	return &testMockStage{name: m.name}
}

func (m *testMockStage) Do(ctx context.Context, state *domain.State) {
	m.doCalled = true
}

type testMockOutput struct {
	name     string
	doCalled bool
}

func (m *testMockOutput) Name() string {
	return m.name
}

func (m *testMockOutput) Copy(cfg map[string]any) domain.Output {
	return &testMockOutput{name: m.name}
}

func (m *testMockOutput) Formats(ff []domain.Format) {
}

func (m *testMockOutput) Do(ctx context.Context, state *domain.State) {
	m.doCalled = true
}

func TestPipelineDo(t *testing.T) {
	mockInput := &testMockInput{name: "test-input"}
	mockStage1 := &testMockStage{name: "test-stage-1"}
	mockStage2 := &testMockStage{name: "test-stage-2"}
	mockOutput := &testMockOutput{name: "test-output"}

	pipeline := NewPipeline(
		"test-pipeline",
		"/test",
		mockInput,
		mockOutput,
		[]domain.Stage{mockStage1, mockStage2},
		nil,
	)

	ctx := context.Background()
	state := &domain.State{
		Request:  &http.Request{},
		Response: nil,
	}

	result := pipeline.Do(ctx, state)

	assert.True(t, result, "Pipeline.Do should return true")
	assert.True(t, mockInput.doCalled, "Input.Do should be called")
	assert.True(t, mockStage1.doCalled, "Stage 1 Do should be called")
	assert.True(t, mockStage2.doCalled, "Stage 2 Do should be called")
	assert.True(t, mockOutput.doCalled, "Output.Do should be called")
}

func TestPipelineDoWithNoStages(t *testing.T) {
	mockInput := &testMockInput{name: "test-input"}
	mockOutput := &testMockOutput{name: "test-output"}

	pipeline := NewPipeline(
		"test-pipeline",
		"/test",
		mockInput,
		mockOutput,
		[]domain.Stage{},
		nil,
	)

	ctx := context.Background()
	state := &domain.State{
		Request:  &http.Request{},
		Response: nil,
	}

	result := pipeline.Do(ctx, state)

	assert.True(t, result, "Pipeline.Do should return true")
	assert.True(t, mockInput.doCalled, "Input.Do should be called")
	assert.True(t, mockOutput.doCalled, "Output.Do should be called")
}

func TestPipelineNameAndRoute(t *testing.T) {
	pipeline := NewPipeline(
		"test-pipeline",
		"/test-route",
		&testMockInput{name: "test-input"},
		&testMockOutput{name: "test-output"},
		[]domain.Stage{},
		nil,
	)

	assert.Equal(t, "test-pipeline", pipeline.Name(), "Pipeline name should match")
	assert.Equal(t, "/test-route", pipeline.Route(), "Pipeline route should match")
}
