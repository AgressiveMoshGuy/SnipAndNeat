package application

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Environment  string        `yaml:"environment"`
	StartTimeout time.Duration `yaml:"start_timeout"`
	StopTimeout  time.Duration `yaml:"stop_timeout"`
	Telegram     Telegram      `yaml:"telegram"`
}

type Telegram struct {
	URL           string `yaml:"url"`
	TelegramToken string `yaml:"telegram_token"`
}

func (c Config) CreateFromFile(path string) (Config, error) {
	cfg := Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
