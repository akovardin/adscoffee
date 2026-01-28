package models

import "encoding/json"

type Limit struct {
	Daily   int  `json:"daily"`
	Total   int  `json:"total"`
	Uniform bool `json:"uniform"`
}

type Budget struct {
	Impressions Limit `json:"impressions"`
	Clicks      Limit `json:"clicks"`
	Money       Limit `json:"money"`
	Conversions Limit `json:"conversions"`
}

func NewBudget(data string) (Budget, error) {
	b := Budget{}

	if data == "" {
		return b, nil
	}

	if err := json.Unmarshal([]byte(data), &b); err != nil {
		return Budget{}, err
	}

	return b, nil
}

func (b Budget) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return ""
	}

	return string(data)
}
