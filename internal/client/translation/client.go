package translation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"net/http"
	"os"
	"time"
)

var (
	cl     http.Client
	bearer string
)

type Client struct {
	conf config.Yandex

	close chan struct{}
}

func NewClient(conf config.Yandex) (*Client, error) {
	client := &Client{
		conf:  conf,
		close: make(chan struct{}),
	}

	err := client.updateAuthToken()
	if err != nil {
		return nil, err
	}

	client.startSchedulingUpdateAuthToken()

	return client, nil
}

func (c *Client) AutoTranslate(text string) (string, error) {
	lang, err := c.DetectLanguage(text)
	if err != nil {
		return "", err
	}

	switch lang {
	case RU:
		return c.Translate(text, EN)
	case EN:
		return c.Translate(text, RU)
	}
	return "", errors.New(constant.UnexpectedLanguageError)
}

func (c *Client) Translate(text string, lang Language) (string, error) {
	translateResp, err := post(newTranslateRequest(c.conf.FolderId, text, lang), c.conf.Api.Translate)
	if err != nil {
		return "", err
	}

	var translate TranslateResponse
	err = json.NewDecoder(translateResp.Body).Decode(&translate)
	if err != nil {
		return "", err
	}

	return translate.Translations[0].Text, nil
}

func (c *Client) DetectLanguage(text string) (Language, error) {
	detectResp, err := post(newDetectLanguageRequest(c.conf.FolderId, text), c.conf.Api.Detect)
	if err != nil {
		return "", err
	}

	var detect DetectLanguageResponse
	err = json.NewDecoder(detectResp.Body).Decode(&detect)
	if err != nil {
		return "", err
	}

	return Language(detect.LanguageCode), nil
}

func (c *Client) Close() {
	c.close <- struct{}{}
}

func (c *Client) startSchedulingUpdateAuthToken() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.updateAuthToken()
				if err != nil {
					ticker.Stop()
					return
				}
			case <-c.close:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *Client) updateAuthToken() error {
	body := bytes.NewBuffer([]byte(fmt.Sprintf("{\"yandexPassportOauthToken\":\"%s\"}", os.Getenv(constant.YandexOauthToken))))

	response, err := http.Post(c.conf.Api.Tokens, "application/json", body)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(constant.YandexTokenError)
	}

	var resp map[string]any
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return err
	}

	bearer = resp["iamToken"].(string)
	return nil
}

func post(request any, endpoint string) (*http.Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	req.Header.Add("Authorization", "Bearer "+bearer)

	return cl.Do(req)
}
