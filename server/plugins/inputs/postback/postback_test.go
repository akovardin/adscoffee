package postback

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestNew(t *testing.T) {
	postback := New()

	assert.NotNil(t, postback)
}

func TestPostback_Name(t *testing.T) {
	postback := New()

	name := postback.Name()

	assert.Equal(t, "inputs.postback", name)
}

func TestPostback_Copy(t *testing.T) {
	postback := New()

	cfgMap := map[string]any{"key": "value"}
	copied := postback.Copy(cfgMap)

	assert.NotNil(t, copied)
	assert.IsType(t, &Postback{}, copied)
}

func TestPostback_Do(t *testing.T) {
	postback := New()

	ctx := context.Background()
	state := &plugins.State{
		Request:    &http.Request{},
		Response:   nil,
		User:       nil,
		Device:     nil,
		Candidates: []ads.Banner{},
		Winners:    []ads.Banner{},
	}

	result := postback.Do(ctx, state)

	assert.True(t, result)
	assert.NotNil(t, state.User)
	assert.NotNil(t, state.Device)
}
