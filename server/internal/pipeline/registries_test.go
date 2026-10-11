package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestRegistries_Get(t *testing.T) {
	in := newInputs([]plugins.Provider[plugins.Input]{inputProvider("inputs.rtb")})
	input, err := in.Get("inputs.rtb", map[string]any{"k": "v"})
	require.NoError(t, err)
	assert.Equal(t, "inputs.rtb", input.Name())

	out := newOutputs([]plugins.Provider[plugins.Output]{outputProvider("outputs.rtb")})
	output, err := out.Get("outputs.rtb", nil)
	require.NoError(t, err)
	assert.Equal(t, "outputs.rtb", output.Name())

	st := newStages([]plugins.Provider[plugins.Stage]{
		stageProvider("stages.banners", func() plugins.Stage { return &mockStage{name: "stages.banners"} }),
	})
	stage, err := st.Get("stages.banners", nil)
	require.NoError(t, err)
	assert.Equal(t, "stages.banners", stage.Name())

	tg := newTargetings([]plugins.Provider[plugins.Targeting]{targetingProvider("targetings.apps")})
	targeting, err := tg.Get("targetings.apps", nil)
	require.NoError(t, err)
	assert.Equal(t, "targetings.apps", targeting.Name())

	lm := newLimits([]plugins.Provider[plugins.Limit]{limitProvider("limits.budget")})
	limit, err := lm.Get("limits.budget", nil)
	require.NoError(t, err)
	assert.Equal(t, "limits.budget", limit.Name())
}

func TestRegistries_GetUnknown(t *testing.T) {
	in := newInputs(nil)
	_, err := in.Get("missing", nil)
	assert.Error(t, err)
}
