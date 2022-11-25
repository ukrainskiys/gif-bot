package main

import (
	"fmt"
	"github.com/ukrainskiys/gif-bot/internal/client/giphy"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"log"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	gif := giphy.NewClient(conf.Giphy)
	list, err := gif.GetGifList(giphy.GIF, "бухать")
	if err != nil {
		return
	}

	for _, g := range list {
		fmt.Println(g)
	}

}
