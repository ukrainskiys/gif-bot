package translation

type DetectLanguageRequest struct {
	FolderId          string     `json:"folderId"`
	Text              string     `json:"text"`
	LanguageCodeHints []Language `json:"languageCodeHints"`
}

func newDetectLanguageRequest(folderId string, text string) DetectLanguageRequest {
	return DetectLanguageRequest{
		FolderId:          folderId,
		Text:              text,
		LanguageCodeHints: []Language{RU, EN},
	}
}

type DetectLanguageResponse struct {
	LanguageCode string `json:"languageCode"`
}

type TranslateRequest struct {
	FolderId           string   `json:"folderId"`
	Texts              []string `json:"texts"`
	TargetLanguageCode Language `json:"targetLanguageCode"`
}

func newTranslateRequest(folderId string, text string, lang Language) TranslateRequest {
	return TranslateRequest{
		FolderId:           folderId,
		Texts:              []string{text},
		TargetLanguageCode: lang,
	}
}

type TranslateResponse struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Text                 string   `json:"text"`
	DetectedLanguageCode Language `json:"detectedLanguageCode"`
}
