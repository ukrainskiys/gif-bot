package giphy

import (
	"encoding/json"
	"errors"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"net/http"
	"net/url"
	"os"
)

type Service struct {
	endpoint string
	token    string
}

func NewService(endpoint string) (*Service, error) {
	client := &Service{
		endpoint: endpoint,
		token:    os.Getenv(constant.GiphyToken),
	}

	if err := client.check(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (s *Service) GetGifList(request SearchGifRequest) ([]Gif, error) {
	get, err := http.Get(s.buildUrl(request.GifType, request.Phrase).String())
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

func (s *Service) buildUrl(typ GifType, phrase string) *url.URL {
	builder, _ := url.Parse(s.endpoint)
	builder = builder.JoinPath(typ.toPath(), "search")
	params := builder.Query()
	params.Add("api_key", s.token)
	params.Add("q", phrase)
	builder.RawQuery = params.Encode()
	return builder
}

func (s *Service) check() error {
	get, err := http.Get(s.buildUrl(GIF, "hello").String())
	if err != nil {
		return err
	}

	if get.StatusCode != 200 {
		return errors.New(constant.GiphyTokenError)
	} else {
		return nil
	}
}
