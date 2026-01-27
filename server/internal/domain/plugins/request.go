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
}

type User struct {
	ID string
}

type Device struct {
	UA string
	IP string
}
