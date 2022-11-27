package cache

import (
	"encoding/json"
	"github.com/ukrainskiys/gif-bot/internal/client/giphy"
	"strings"
)

type Phrase struct {
	FirstLang  string
	SecondLang string
}

func is(phrase string, s string) bool {
	words := strings.Split(phrase, "-")
	return words[0] == s || words[1] == s
}

func (p *Phrase) toKey(gifType giphy.GifType) string {
	return gifType.String() + "-" + strings.ToUpper(p.FirstLang) + "-" + strings.ToUpper(p.SecondLang)
}

func parsePhrase(s string) (giphy.GifType, Phrase) {
	words := strings.Split(s, "-")
	return giphy.ParseType(words[0]), Phrase{FirstLang: words[1], SecondLang: words[2]}
}

type AccountInfo struct {
	GifType   giphy.GifType
	GifsCache map[string][]string
}

func (ai *AccountInfo) UpdateGifs(phrase Phrase, urls map[string]struct{}) {
	key := phrase.toKey(ai.GifType)
	ai.GifsCache[key] = make([]string, len(urls))
	i := 0
	for val := range urls {
		ai.GifsCache[key][i] = val
		i++
	}
}

func (ai *AccountInfo) GetGifsByPhrase(s string) ([]string, Phrase) {
	searchable := strings.ToUpper(s)
	for phrase, gifs := range ai.GifsCache {
		if is(phrase, searchable) {
			gifType, phr := parsePhrase(phrase)
			if ai.GifType == gifType {
				return gifs, phr
			}
		}
	}
	return nil, Phrase{}
}

func (ai *AccountInfo) UnmarshalJSON(b []byte) error {
	var data map[string]any
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	ai.GifType = giphy.ParseType(data["GifType"].(string))
	ai.GifsCache = make(map[string][]string)

	if data["GifsCache"] != nil {
		for k, v := range data["GifsCache"].(map[string]any) {
			arr := v.([]any)
			gifType, phrase := parsePhrase(k)
			key := phrase.toKey(gifType)
			ai.GifsCache[k] = make([]string, len(arr))
			for i, val := range arr {
				ai.GifsCache[key][i] = val.(string)
			}
		}
	}
	return nil
}
