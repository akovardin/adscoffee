package domain

import "net/http"

type Input interface {
	Process(r *http.Request) bool
	Name() string
}
