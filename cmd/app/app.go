package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/app"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"time"
)

func main() {
	now := time.Now()

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	telegram, err := app.NewBot(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started [%v]", time.Since(now))

	telegram.Run()
}
