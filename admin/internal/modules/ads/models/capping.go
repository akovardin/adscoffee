package models

import "encoding/json"

const CappingKeyTemplate = "%s:%s:%s" //  uid, item, id

type Capping struct {
	Count  int `json:"count"`
	Period int `json:"period"`
}

func NewCapping(data string) (Capping, error) {
	c := Capping{}

	if data == "" {
		return c, nil
	}

	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return Capping{}, err
	}

	return c, nil
}

func (c Capping) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(data)
}
