package plugins

import (
	"net/http"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type State struct {
	Request    *http.Request
	Response   http.ResponseWriter
	User       *User
	Device     *Device
	Candidates []ads.Banner
	Winners    []ads.Banner
	Placement  *Placement
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
