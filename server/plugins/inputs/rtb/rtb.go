package rtb

import (
	"net/http"

	"go.ads.coffee/server/config"
)

type Rtb struct {
}

func New(config config.Config) *Rtb {
	input := &Rtb{}

	return input
}

func (rtb *Rtb) Build(cfg map[string]any) *Rtb {
	return &Rtb{}
}

func (rtb *Rtb) Name() string {
	return "rtb"
}

func (rtb *Rtb) Process(r *http.Request) bool {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины

	return true
}
