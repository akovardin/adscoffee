package domain

import "net/http"

type State struct {
	Request    *http.Request
	Response   http.ResponseWriter
	User       *User
	Device     *Device
	Candidates []Banner
	Winners    []Banner
}

type User struct {
	ID string
}

type Device struct {
	UA string
	IP string
}
