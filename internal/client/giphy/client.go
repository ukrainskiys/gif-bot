package giphy

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	endpoint string
	token    string
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    os.Getenv("GIPHY_TOKEN"),
	}
}

func (c *Client) GetGifList(typ GifType, phrase string) ([]Gif, error) {
	get, err := http.Get(c.buildUrl(typ, phrase).String())
	if err != nil {
		return nil, err
	}

	var response Response
	if err = json.NewDecoder(get.Body).Decode(&response); err != nil {
		return nil, err
	} else {
		return response.Data, nil
	}
}

func (c *Client) buildUrl(typ GifType, phrase string) *url.URL {
	builder, _ := url.Parse(c.endpoint)
	builder = builder.JoinPath(typ.toPath(), "search")
	params := builder.Query()
	params.Add("api_key", c.token)
	params.Add("q", phrase)
	builder.RawQuery = params.Encode()
	return builder
}
