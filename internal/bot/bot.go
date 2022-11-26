package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/bot/handler"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"os"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	handle *handler.BotHandler
}

func NewBot(handle *handler.BotHandler) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(os.Getenv(constant.TelegramToken))
	if err != nil {
		return nil, err
	}

	return &Bot{api, handle}, nil
}

func (b *Bot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 0

	for update := range b.api.GetUpdatesChan(updateConfig) {
		if update.Message == nil {
			continue
		}

		log.Printf("GET update Chat.ID=%d Text=%s", update.Message.Chat.ID, update.Message.Text)

		msg, err := b.handle.HandleMessage(update.Message)
		if err != nil {
			log.Warn(err)
			continue
		}

		b.send(msg)
	}
}

func (b *Bot) send(c tgbotapi.Chattable) {
	res, err := b.api.Send(c)
	if err != nil {
		log.Warn(err)
	}
	log.Printf(`POST update Chat.ID=%d Text="%s"`, res.Chat.ID, res.Text)
}