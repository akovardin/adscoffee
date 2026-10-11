package postback

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestNew(t *testing.T) {
	v, err := New().Factory(nil)
	require.NoError(t, err)

	assert.IsType(t, &Postback{}, v)
}

func TestPostback_Name(t *testing.T) {
	postback := &Postback{}

	name := postback.Name()

	assert.Equal(t, "inputs.postback", name)
}

func TestPostback_Do(t *testing.T) {
	postback := &Postback{}

	ctx := context.Background()
	state := &plugins.State{
		Request:    &http.Request{},
		Response:   nil,
		User:       nil,
		Device:     nil,
		Candidates: []*ads.Banner{},
		Winners:    []*ads.Banner{},
	}

	result := postback.Do(ctx, state)

	assert.True(t, result)
	assert.NotNil(t, state.User)
	assert.NotNil(t, state.Device)
}
