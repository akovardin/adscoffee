package postback

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/config"
	"go.ads.coffee/platform/server/domain"
)

func TestNew(t *testing.T) {
	cfg := config.Config{}
	postback := New(cfg)

	assert.NotNil(t, postback)
}

func TestPostback_Name(t *testing.T) {
	cfg := config.Config{}
	postback := New(cfg)

	name := postback.Name()

	assert.Equal(t, "inputs.postback", name)
}

func TestPostback_Copy(t *testing.T) {
	cfg := config.Config{}
	postback := New(cfg)

	cfgMap := map[string]any{"key": "value"}
	copied := postback.Copy(cfgMap)

	assert.NotNil(t, copied)
	assert.IsType(t, &Postback{}, copied)
}

func TestPostback_Do(t *testing.T) {
	cfg := config.Config{}
	postback := New(cfg)

	ctx := context.Background()
	state := &domain.State{
		Request:    &http.Request{},
		Response:   nil,
		User:       nil,
		Device:     nil,
		Candidates: []domain.Banner{},
		Winners:    []domain.Banner{},
	}

	result := postback.Do(ctx, state)

	assert.True(t, result)
	assert.NotNil(t, state.User)
	assert.NotNil(t, state.Device)
}
