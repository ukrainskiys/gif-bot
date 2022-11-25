package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type AppConfig struct {
	Giphy  string `yaml:"giphy"`
	Yandex Yandex `yaml:"yandex"`
}

type Yandex struct {
	Api struct {
		Translate string `yaml:"translate"`
		Detect    string `json:"detect"`
		Tokens    string `yaml:"tokens"`
	} `yaml:"api"`
	FolderId string `yaml:"folder-id"`
}

func NewConfig() (*AppConfig, error) {
	const fileName = "config.yml"

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	config := &AppConfig{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing YAML config: %s\n", fileName))
	}

	return config, nil
}
