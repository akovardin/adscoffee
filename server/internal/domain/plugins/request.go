package plugins

import (
	"context"
	"net/http"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type State struct {
	RequestID  string
	ClickID    string
	Request    *http.Request
	Response   http.ResponseWriter
	User       *User
	Device     *Device
	Candidates []ads.Banner
	Winners    []ads.Banner
	Placement  *Placement
}

func (s *State) Value(key any) any {
	return s.Request.Context().Value("action")
}

func (s *State) WithValue(key, value any) {
	ctx := s.Request.Context()
	ctx = context.WithValue(ctx, key, value)
	s.Request = s.Request.WithContext(ctx)
}

type User struct {
	ID string
}

type Device struct {
	UA string
	IP string
}

type Placement struct {
	ID    string
	Units []ads.Unit
}
