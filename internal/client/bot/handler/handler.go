package handler

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/client/giphy"
	"github.com/ukrainskiys/gif-bot/internal/client/translation"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"github.com/ukrainskiys/gif-bot/pkg/concurent"

	"math/rand"
)

var (
	tgKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(giphy.GIF.String()),
			tgbotapi.NewKeyboardButton(giphy.STICKER.String()),
		),
	)

	conditions = make(map[int64]giphy.GifType)
)

type BotHandler struct {
	giftApi    *giphy.Client
	translator *translation.Client
}

func NewBotHandler(conf *config.AppConfig) (*BotHandler, error) {
	translator, err := translation.NewClient(conf.Yandex)
	if err != nil {
		return nil, err
	}

	giftApi, err := giphy.NewClient(conf.Giphy)
	if err != nil {
		return nil, err
	}

	handler := &BotHandler{
		giftApi,
		translator,
	}
	return handler, nil
}

func (bh *BotHandler) HandleMessage(message *tgbotapi.Message) (tgbotapi.Chattable, error) {
	switch message.Text {
	case giphy.GIF.String(), giphy.STICKER.String():
		gifType := giphy.ParseType(message.Text)
		conditions[message.Chat.ID] = gifType
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
	gifType, ok := conditions[chatId]
	if !ok {
		return tgbotapi.AnimationConfig{}, &GifTypeNotSpecifiedError{}
	}

	done := make(chan struct{})
	links := concurent.NewSlice[string](0)

	go bh.searchGifs(done, giphy.SearchGifRequest{Phrase: phrase, GifType: gifType}, links)

	translate, err := bh.translator.AutoTranslate(phrase)
	if err != nil {
		<-done
		return tgbotapi.AnimationConfig{}, err
	}

	go bh.searchGifs(done, giphy.SearchGifRequest{Phrase: translate, GifType: gifType}, links)

	<-done
	<-done

	if links.Size() == 0 {
		return tgbotapi.AnimationConfig{}, &GifsNotFoundError{}
	} else {
		return tgbotapi.NewAnimation(chatId, tgbotapi.FileURL(links.Get(rand.Intn(links.Size())))), nil
	}
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
