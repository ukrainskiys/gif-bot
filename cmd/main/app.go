package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/bot"
	"github.com/ukrainskiys/gif-bot/internal/bot/handler"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"time"
)

func main() {
	//test()

	now := time.Now()

	conf, err := config.NewConfig()
	handle, err := handler.NewBotHandler(conf)
	telegram, err := bot.NewBot(handle)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started [%v]", time.Since(now))

	defer handle.Close()
	telegram.Run()
}

//func test() {
//	conf, _ := config.NewConfig()
//	cacheClient := cache.NewClient(conf.RedisConfig)
//}
