package ads

import (
	"encoding/json"
	"strings"
)

const base = "//adexproadmin.hb.vkcloud-storage.ru"

type Image struct {
	Url string `json:"url"`
}

func NewImage(data string) (Image, error) {
	i := Image{}

	if data == "" {
		return i, nil
	}

	if err := json.Unmarshal([]byte(data), &i); err != nil {
		return Image{}, err
	}

	return i, nil
}

func (i Image) Full(cdn string) string {
	if cdn == "" {
		return "https:" + i.Url
	}

	img := strings.Replace(i.Url, base, cdn, -1)
	img = strings.Replace(img, "file", "file.image", -1)

	return img
}
