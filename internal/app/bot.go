package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"github.com/ukrainskiys/gif-bot/internal/handlers"
	"os"
	"time"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *handlers.BotHandler
}

func NewBot(conf *config.AppConfig) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(os.Getenv(constant.TelegramToken))
	if err != nil {
		return nil, err
	}

	handler, err := handlers.NewBotHandler(conf)
	if err != nil {
		return nil, err
	}

	return &Bot{api, handler}, nil
}

func (b *Bot) Run() {
	defer b.handler.Close()
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 0

	for update := range b.api.GetUpdatesChan(updateConfig) {
		now := time.Now()
		if update.Message == nil {
			continue
		}

		log.Printf("GET update Chat.ID=%d Text=%s", update.Message.Chat.ID, update.Message.Text)

		msg, err := b.handler.HandleMessage(update.Message)
		if err != nil {
			log.Warn(err)
			continue
		}

		b.send(msg, time.Since(now))
	}
}

func (b *Bot) send(c tgbotapi.Chattable, duration time.Duration) {
	res, err := b.api.Send(c)
	if err != nil {
		log.Warn(err)
	} else {
		log.Printf(`POST update Chat.ID=%d Text="%s" [%v]`, res.Chat.ID, res.Text, duration)
	}
}
