package models

import "encoding/json"

type ExcludeIncludeOrAnd struct {
	IncludeOr []string `json:"include_or"`
	ExcludeOr []string `json:"exclude_or"`

	IncludeAnd []string `json:"include_and"`
	ExcludeAnd []string `json:"exclude_and"`
}

type ExcludeInclude struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type Targeting struct {
	Bundle   ExcludeIncludeOrAnd `json:"bundle"`
	Audience ExcludeIncludeOrAnd `json:"audience"`
	Bapp     ExcludeIncludeOrAnd `json:"bapp"`
	Country  ExcludeIncludeOrAnd `json:"country"`
	Region   ExcludeIncludeOrAnd `json:"region"`
	City     ExcludeIncludeOrAnd `json:"city"`
	IP       ExcludeInclude      `json:"ip"`
	Network  ExcludeIncludeOrAnd `json:"network"`
}

func NewTargeting(data string) (Targeting, error) {
	t := Targeting{}

	if data == "" {
		return t, nil
	}

	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return Targeting{}, err
	}

	return t, nil
}

func (t Targeting) String() string {
	data, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(data)
}
