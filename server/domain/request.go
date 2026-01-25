package domain

import "net/http"

type State struct {
	Request  *http.Request
	Response http.ResponseWriter
}
