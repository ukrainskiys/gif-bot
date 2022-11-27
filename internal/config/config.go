package config

import (
	"errors"
	"fmt"
	"github.com/ukrainskiys/gif-bot/internal/constant"
	"gopkg.in/yaml.v2"
	"os"
)

type AppConfig struct {
	Giphy        string       `yaml:"giphy"`
	YandexConfig YandexConfig `yaml:"yandex"`
	RedisConfig  RedisConfig  `yaml:"redis"`
}

type YandexConfig struct {
	Api struct {
		Translate string `yaml:"translate"`
		Detect    string `json:"detect"`
		Tokens    string `yaml:"tokens"`
	} `yaml:"api"`
	FolderId string `yaml:"folder-id"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

func NewConfig() (*AppConfig, error) {
	bytes, err := os.ReadFile(constant.ConfigName)
	if err != nil {
		return nil, err
	}

	config := &AppConfig{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing YAML config: %s\n", constant.ConfigName))
	}

	return config, nil
}
