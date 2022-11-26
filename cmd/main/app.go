package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/bot"
	"github.com/ukrainskiys/gif-bot/internal/bot/handler"
	"github.com/ukrainskiys/gif-bot/internal/config"
)

func main() {
	conf, err := config.NewConfig()
	handle, err := handler.NewBotHandler(conf)
	telegram, err := bot.NewBot(handle)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()
	telegram.Run()
}
