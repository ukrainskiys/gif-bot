package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/app"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"os"
	"time"
)

func init() {
	fileName := fmt.Sprintf("log/%s.log", time.Now().UTC().Format(constant.LayoutCustom))
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
}

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
