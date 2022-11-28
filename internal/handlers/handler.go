package handlers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"github.com/ukrainskiys/gif-bot/internal/services/cache"
	"github.com/ukrainskiys/gif-bot/internal/services/giphy"
	"github.com/ukrainskiys/gif-bot/internal/services/translation"
	"github.com/ukrainskiys/gif-bot/pkg/concurent"
	"math/rand"
	"time"
)

var (
	tgKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(giphy.GIF.String()),
			tgbotapi.NewKeyboardButton(giphy.STICKER.String()),
		),
	)
)

type BotHandler struct {
	giftApi    *giphy.Service
	translator *translation.Service
	cacheCl    *cache.Service
}

func NewBotHandler(conf *config.AppConfig) (*BotHandler, error) {
	translator, err := translation.NewService(conf.YandexConfig)
	if err != nil {
		return nil, err
	}

	giftApi, err := giphy.NewService(conf.Giphy)
	if err != nil {
		return nil, err
	}

	cacheClient, err := cache.NewService(conf.RedisConfig)
	if err != nil {
		return nil, err
	}

	return &BotHandler{giftApi, translator, cacheClient}, nil
}

func (bh *BotHandler) HandleMessage(message *tgbotapi.Message) (tgbotapi.Chattable, error) {
	switch message.Text {
	case giphy.GIF.String(), giphy.STICKER.String():
		gifType := giphy.ParseType(message.Text)
		bh.cacheCl.SetNewTypeForAccount(message.Chat.ID, gifType)
		return getDefaultMessage(message.Chat.ID, gifType), nil

	default:
		if len(message.Text) >= 50 {
			return messageAndErrorNil(message.Chat.ID, constant.VeryLongTextMessage)
		}

		msg, err := bh.getAnimation(message.Chat.ID, message.Text)
		if errors.Is(err, &GifTypeNotSpecifiedError{}) {
			return messageAndErrorNil(message.Chat.ID, constant.NeedSelectTypeMessage)
		} else if errors.Is(err, &GifsNotFoundError{}) {
			return messageAndErrorNil(message.Chat.ID, constant.GifNotFoundMessage)
		} else if err != nil {
			return nil, err
		} else {
			msg.ReplyMarkup = tgKeyboard
			return msg, nil
		}
	}
}

func (bh *BotHandler) Close() {
	bh.translator.Close()
}

func getDefaultMessage(chatId int64, gifType giphy.GifType) tgbotapi.MessageConfig {
	var text string
	if gifType == giphy.GIF {
		text = constant.PhraseForGifMessage
	} else {
		text = constant.PhraseForStickerMessage
	}

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = tgKeyboard
	return msg
}

func (bh *BotHandler) getAnimation(chatId int64, phrase string) (tgbotapi.AnimationConfig, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	accountInfo, ok := bh.cacheCl.GetAccountInfo(chatId)

	if !ok {
		return tgbotapi.AnimationConfig{}, &GifTypeNotSpecifiedError{}
	}

	var gif string
	gifs, phr := accountInfo.GetGifsByPhrase(phrase)
	if len(gifs) > 0 {
		idx := rand.Intn(len(gifs))
		gif = gifs[idx]
		gifs = append(gifs[:idx], gifs[idx+1:]...)
		accountInfo.UpdateGifs(phr, removeRepeats(gifs))
		bh.cacheCl.Set(chatId, accountInfo)
	} else {
		links, phr := bh.getGifLinks(accountInfo, cache.Phrase{FirstLang: phrase})

		if links.Size() == 0 {
			return tgbotapi.AnimationConfig{}, &GifsNotFoundError{}
		} else {
			idx := rand.Intn(links.Size())
			gif = links.Get(idx)
			links.Remove(idx)
			accountInfo.UpdateGifs(phr, removeRepeats(links.Array()))
			bh.cacheCl.Set(chatId, accountInfo)
		}
	}

	return tgbotapi.NewAnimation(chatId, tgbotapi.FileURL(gif)), nil
}

func (bh *BotHandler) getGifLinks(info cache.AccountInfo, phrase cache.Phrase) (*concurent.Slice[string], cache.Phrase) {
	done := make(chan struct{})
	links := concurent.NewSlice[string]()

	go bh.searchGifs(done, giphy.SearchGifRequest{Phrase: phrase.FirstLang, GifType: info.GifType}, links)

	translate, err := bh.translator.AutoTranslate(phrase.FirstLang)
	if err != nil {
		<-done
		return nil, cache.Phrase{}
	}

	go bh.searchGifs(done, giphy.SearchGifRequest{Phrase: translate, GifType: info.GifType}, links)
	phrase.SecondLang = translate

	<-done
	<-done

	return links, phrase
}

func (bh *BotHandler) searchGifs(done chan struct{}, searchRequest giphy.SearchGifRequest, links *concurent.Slice[string]) {
	gifs, err := bh.giftApi.GetGifList(searchRequest)
	if err != nil {
		log.Warn(err)
		return
	}

	for _, gif := range gifs {
		links.Append(string(gif.Url))
	}

	done <- struct{}{}
}

func messageAndErrorNil(chatId int64, text string) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = tgKeyboard
	return msg, nil
}

func removeRepeats(array []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, val := range array {
		set[val] = struct{}{}
	}

	return set
}
