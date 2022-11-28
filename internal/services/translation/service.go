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
	client http.Client
	bearer string
)

type Service struct {
	conf config.YandexConfig

	close chan struct{}
}

func NewService(conf config.YandexConfig) (*Service, error) {
	cl := &Service{
		conf:  conf,
		close: make(chan struct{}),
	}

	err := cl.updateAuthToken()
	if err != nil {
		return nil, err
	}

	cl.startSchedulingUpdateAuthToken()

	return cl, nil
}

func (s *Service) AutoTranslate(text string) (string, error) {
	lang, err := s.DetectLanguage(text)
	if err != nil {
		return "", err
	}

	switch lang {
	case RU:
		return s.Translate(text, EN)
	case EN:
		return s.Translate(text, RU)
	}
	return "", errors.New(constant.UnexpectedLanguageError)
}

func (s *Service) Translate(text string, lang Language) (string, error) {
	translateResp, err := post(newTranslateRequest(s.conf.FolderId, text, lang), s.conf.Api.Translate)
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

func (s *Service) DetectLanguage(text string) (Language, error) {
	detectResp, err := post(newDetectLanguageRequest(s.conf.FolderId, text), s.conf.Api.Detect)
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

func (s *Service) Close() {
	s.close <- struct{}{}
}

func (s *Service) startSchedulingUpdateAuthToken() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.updateAuthToken()
				if err != nil {
					ticker.Stop()
					return
				}
			case <-s.close:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Service) updateAuthToken() error {
	body := bytes.NewBuffer([]byte(fmt.Sprintf("{\"yandexPassportOauthToken\":\"%s\"}", os.Getenv(constant.YandexOauthToken))))

	response, err := http.Post(s.conf.Api.Tokens, "application/json", body)
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

	return client.Do(req)
}
