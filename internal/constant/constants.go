package constant

const (
	ConfigName = "config.yml"

	VeryLongTextMessage     = "Фраза слишком большая! Длинна ее не должна превышать 50 символов."
	GifNotFoundMessage      = "Поиск по этой фразе не дал результатов:("
	NeedSelectTypeMessage   = "Нужно указать тип [GIF/STICKER]."
	PhraseForGifMessage     = "Введите фразу для подбора гифки ⬇️"
	PhraseForStickerMessage = "Введите фразу для подбора стикера ⬇️"

	GiphyTokenError         = "giphy client doesn't worked (check auth token)"
	YandexTokenError        = "yandex clint doesn't worked (check auth token)"
	UnexpectedLanguageError = "unexpected language"

	TelegramToken    = "TELEGRAM_TOKEN"
	GiphyToken       = "GIPHY_TOKEN"
	YandexOauthToken = "YANDEX_OAUTH"
)
