package giphy

import (
	"encoding/json"
)

type Response struct {
	Data []Gif `json:"data"`
}

type Gif struct {
	Id   string  `json:"id"`
	Type GifType `json:"type"`
	Url  Url     `json:"images"`
}

type Url string

func (u *Url) UnmarshalJSON(b []byte) error {
	var data map[string]any
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	*u = Url(data["original"].(map[string]any)["url"].(string))
	return nil
}
