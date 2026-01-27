package ads

import "encoding/json"

type Limit struct {
	Daily   int64 `json:"daily"`
	Total   int64 `json:"total"`
	Uniform bool  `json:"uniform"`	
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
