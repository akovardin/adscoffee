package pipeline

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/server/domain"
)

// Mock implementations for testing with unique names to avoid conflicts
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
	// Mock implementation
}

func (m *testMockOutput) Do(ctx context.Context, state *domain.State) {
	m.doCalled = true
}

func TestPipelineDo(t *testing.T) {
	// Create mock components
	mockInput := &testMockInput{name: "test-input"}
	mockStage1 := &testMockStage{name: "test-stage-1"}
	mockStage2 := &testMockStage{name: "test-stage-2"}
	mockOutput := &testMockOutput{name: "test-output"}

	// Create pipeline
	pipeline := NewPipeline(
		"test-pipeline",
		"/test",
		mockInput,
		mockOutput,
		[]domain.Stage{mockStage1, mockStage2},
		nil,
	)

	// Create test context and state
	ctx := context.Background()
	state := &domain.State{
		Request:  &http.Request{},
		Response: nil,
	}

	// Execute pipeline
	result := pipeline.Do(ctx, state)

	// Assertions
	assert.True(t, result, "Pipeline.Do should return true")

	// Verify input was called
	assert.True(t, mockInput.doCalled, "Input.Do should be called")

	// Verify stages were called
	assert.True(t, mockStage1.doCalled, "Stage 1 Do should be called")
	assert.True(t, mockStage2.doCalled, "Stage 2 Do should be called")

	// Verify output was called
	assert.True(t, mockOutput.doCalled, "Output.Do should be called")
}

func TestPipelineDoWithNoStages(t *testing.T) {
	// Create mock components
	mockInput := &testMockInput{name: "test-input"}
	mockOutput := &testMockOutput{name: "test-output"}

	// Create pipeline with no stages
	pipeline := NewPipeline(
		"test-pipeline",
		"/test",
		mockInput,
		mockOutput,
		[]domain.Stage{}, // No stages
		nil,
	)

	// Create test context and state
	ctx := context.Background()
	state := &domain.State{
		Request:  &http.Request{},
		Response: nil,
	}

	// Execute pipeline
	result := pipeline.Do(ctx, state)

	// Assertions
	assert.True(t, result, "Pipeline.Do should return true")

	// Verify input was called
	assert.True(t, mockInput.doCalled, "Input.Do should be called")

	// Verify output was called
	assert.True(t, mockOutput.doCalled, "Output.Do should be called")
}

func TestPipelineNameAndRoute(t *testing.T) {
	// Create pipeline
	pipeline := NewPipeline(
		"test-pipeline",
		"/test-route",
		&testMockInput{name: "test-input"},
		&testMockOutput{name: "test-output"},
		[]domain.Stage{},
		nil,
	)

	// Assertions
	assert.Equal(t, "test-pipeline", pipeline.Name(), "Pipeline name should match")
	assert.Equal(t, "/test-route", pipeline.Route(), "Pipeline route should match")
}
