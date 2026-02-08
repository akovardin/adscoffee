package ads

import (
	"encoding/json"
	"strings"
)

const base = "//platform.hb.vkcloud-storage.ru"

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

	img := strings.ReplaceAll(i.Url, base, cdn)
	img = strings.ReplaceAll(img, "file", "file.image")

	return img
}
