package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Enviroment   string        `env:"ENVIRONMENT"`
	StartTimeout time.Duration `env:"START_TIMEOUT"`
	StopTimeout  time.Duration `env:"STOP_TIMEOUT"`

	Telegram struct {
		Url           string `env:"TELEGRAM_URL"`
		TelegramToken string `env:"TELEGRAM_TOKEN"`
	}

	EmailServer struct {
		Host      string `env:"EMAIL_SMTP_SERVER"`
		Port      string `env:"EMAIL_SMTP_PORT" default:"25"`
		UseAuth   bool   `env:"EMAIL_SMTP_USE_AUTH" default:"true"`
		Identity  string `env:"EMAIL_SMTP_IDENTITY"`
		Username  string `env:"EMAIL_SMTP_USERNAME"`
		Password  string `env:"EMAIL_SMTP_PASSWORD"`
		Recipient string `env:"EMAIL_RECIPIENT"`
	}

	Scheduler struct {
		CheckStateWaitTime               time.Duration `env:"CHECK_STATE_WAIT_TIME" default:"600s"`
		PeriodTimeToCheckDropTransaction time.Duration `env:"PERIOD_TIME_TO_CHECK_DROP_TRANSACTION" default:"5s"`

		MailEnabled      bool          `env:"MAIL_ENABLED" default:"false"`
		MailRepeatPeriod time.Duration `env:"MAIL_REPEAT_PERIOD" default:"2h"`
	}

	Server struct {
		Address      string        `env:"ADDRESS"`
		StartTimeout time.Duration `env:"SERVER_START_TIMEOUT"`
		ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
		WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	}

	OzonClientConfig struct {
		URL          string        `env:"OZON_URL"`
		APIKey       string        `env:"OZON_API_KEY"`
		ClientID     string        `env:"OZON_CLIENT_ID"`
		Timeout      time.Duration `env:"OZON_TIMEOUT" default:"10s"`
		MemcacheHost string        `env:"MEMCACHE_HOST"`
	}

	// Config конфигурация подключения к PostgreSQL
	DBConfig struct {
		Host            string        `env:"DB_HOST"`
		Schema          string        `env:"DB_SCHEMA" default:"public"`
		User            string        `env:"DB_USER"`
		Password        string        `env:"DB_PASSWORD"`
		Name            string        `env:"DB_NAME"`
		Port            int           `env:"DB_PORT"`
		MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" default:"25"`
		MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" default:"50"`
		StartTimeout    time.Duration `env:"DB_START_TIMEOUT" default:"20s"`
		ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" default:"30s"`
	}
}

func CreateFromFile() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// String implementation fmt.Stringer
func (conf *Config) StringDB() string {
	return fmt.Sprintf(
		"sqlite3://%s.db?_auth&_auth_user=%s&_auth_pass=%s",
		conf.DBConfig.Name,
		conf.DBConfig.User,
		conf.DBConfig.Password,
	)
}
