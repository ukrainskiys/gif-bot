package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ukrainskiys/gif-bot/internal/client/bot/handler"
	"log"
	"os"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	handle *handler.BotHandler
}

func NewBot(handle *handler.BotHandler) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return nil, err
	}
	api.Debug = true

	return &Bot{api, handle}, nil
}

func (b *Bot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 0

	for update := range b.api.GetUpdatesChan(updateConfig) {
		if update.Message == nil {
			continue
		}

		msg, err := b.handle.HandleMessage(update.Message)
		if err != nil {
			log.Panic(err)
		}

		if _, err := b.api.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
