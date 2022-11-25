package main

import (
	"fmt"
	"github.com/ukrainskiys/gif-bot/internal/client/translation"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"log"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	translator, err := translation.NewClient(conf.Yandex)
	if err != nil {
		log.Fatal(err)
	}

	phrase, err := translator.AutoTranslate("здарова пацаны")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(phrase)

	translator.Close()

}
