package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"github.com/ukrainskiys/gif-bot/internal/services/giphy"
	"strconv"
)

type Service struct {
	client *redis.Client
}

func NewService(conf config.RedisConfig) (*Service, error) {
	cl := &Service{redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + strconv.Itoa(conf.Port),
		Password: conf.Password,
		DB:       0,
	})}

	return cl, cl.client.Ping().Err()
}

func (s *Service) SetNewTypeForAccount(chatId int64, gifType giphy.GifType) {
	account, ok := s.GetAccountInfo(chatId)
	if ok {
		account.GifType = gifType
		s.Set(chatId, account)
	} else {
		s.Set(chatId, AccountInfo{GifType: gifType, GifsCache: map[string][]string{}})
	}
}

func (s *Service) GetAccountInfo(chatId int64) (AccountInfo, bool) {
	result, err := s.client.Get(toKey(chatId)).Result()
	if err != nil {
		return AccountInfo{}, false
	}

	var account AccountInfo
	if err = json.Unmarshal([]byte(result), &account); err != nil {
		return AccountInfo{GifsCache: map[string][]string{}}, false
	} else {
		return account, true
	}
}

func (s *Service) Set(chatId int64, info AccountInfo) {
	marshal, err := json.Marshal(info)
	if err != nil {
		log.Error(err)
	}
	s.client.Set(toKey(chatId), marshal, 0)
}

func toKey(chatId int64) string {
	return strconv.FormatInt(chatId, 10)
}
