package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainskiys/gif-bot/internal/client/giphy"
	"github.com/ukrainskiys/gif-bot/internal/config"
	"strconv"
)

type Cache struct {
	client *redis.Client
}

func NewClient(conf config.RedisConfig) (*Cache, error) {
	cl := &Cache{redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + strconv.Itoa(conf.Port),
		Password: conf.Password,
		DB:       0,
	})}

	cl.client.Set("test", "123", 0)

	return cl, cl.client.Ping().Err()
}

func (c *Cache) SetNewTypeForAccount(chatId int64, gifType giphy.GifType) {
	account, ok := c.GetAccountInfo(chatId)
	if ok {
		account.GifType = gifType
		c.Set(chatId, account)
	} else {
		c.Set(chatId, AccountInfo{GifType: gifType, GifsCache: map[string][]string{}})
	}
}

func (c *Cache) GetAccountInfo(chatId int64) (AccountInfo, bool) {
	result, err := c.client.Get(toKey(chatId)).Result()
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

func (c *Cache) Set(chatId int64, info AccountInfo) {
	marshal, err := json.Marshal(info)
	if err != nil {
		log.Error(err)
	}
	c.client.Set(toKey(chatId), marshal, 0)
}

func toKey(chatId int64) string {
	return strconv.FormatInt(chatId, 10)
}
