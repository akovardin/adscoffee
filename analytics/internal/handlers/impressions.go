package handlers

import "fmt"

type Impressions struct {
}

func NewImpressions() *Impressions {
	return &Impressions{}
}

func (i *Impressions) Run() {
	fmt.Println("prepare impressions")
}
