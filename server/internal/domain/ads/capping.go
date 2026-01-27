package ads

import "encoding/json"

const CappingKeyTemplate = "%s:%s:%s" //  uid, item, id

type Capping struct {
	Count  int64 `json:"count"`
	Period int64 `json:"period"`
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

type CappingInfo struct {
	LastSeen int64 `json:"last_seen"`
	Count    int64 `json:"count"`
}
