package giphy

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	endpoint string
	token    string
}

func NewClient(endpoint string) (*Client, error) {
	client := &Client{
		endpoint: endpoint,
		token:    os.Getenv("GIPHY_TOKEN"),
	}

	if err := client.check(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (c *Client) GetGifList(request SearchGifRequest) ([]Gif, error) {
	get, err := http.Get(c.buildUrl(request.GifType, request.Phrase).String())
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

func (c *Client) check() error {
	get, err := http.Get(c.buildUrl(GIF, "hello").String())
	if err != nil {
		return err
	}

	if get.StatusCode != 200 {
		return errors.New("giphy client doesn't worked (check auth token)")
	} else {
		return nil
	}
}
