package cache

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"sync"
)

const usersFileName = "users.json"

type UserCache struct {
	users map[int64]any
}

func (uc *UserCache) MarshalJSON() ([]byte, error) {
	users := make(map[string][]int64)
	users["users"] = make([]int64, len(uc.users))
	i := 0
	for userId := range uc.users {
		users["users"][i] = userId
		i++
	}
	return json.Marshal(users)
}

func (uc *UserCache) UnmarshalJSON(b []byte) error {
	users := make(map[string][]int64)
	err := json.Unmarshal(b, &users)
	if err != nil {
		return err
	}
	uc.users = make(map[int64]any)
	for _, userId := range users["users"] {
		uc.users[userId] = nil
	}
	return nil
}

func NewUserCache() *UserCache {
	return &UserCache{users: make(map[int64]any)}
}

func (uc *UserCache) UpdateUserCache(u *tgbotapi.User) {
	if !u.IsBot {
		oldLen := len(uc.users)
		uc.users[u.ID] = nil
		if len(uc.users) > oldLen {
			uc.writeUsers()
		}
	}
}

func (uc *UserCache) loadUsers() {
	file, err := os.Open(usersFileName)
	if err != nil {
		return
	}

	var userCache UserCache
	err = json.NewDecoder(file).Decode(&userCache)
}

func (uc *UserCache) writeUsers() {
	mx := sync.Mutex{}
	go func() {
		mx.Lock()
		defer mx.Unlock()
		file, err := os.Create(usersFileName)
		if err != nil {
			return
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(uc)
	}()
}
