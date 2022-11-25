package main

import (
	"github.com/ukrainskiys/gif-bot/internal/client/bot"
	"github.com/ukrainskiys/gif-bot/internal/client/bot/handler"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"log"
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
