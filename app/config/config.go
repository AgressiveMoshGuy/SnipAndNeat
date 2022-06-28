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
		Host            string        `yaml:"host"`
		Schema          string        `yaml:"schema" default:"public"`
		User            string        `yaml:"user"`
		Password        string        `yaml:"password"`
		Name            string        `yaml:"name"`
		Port            int           `yaml:"port"`
		MaxIdleConns    int           `yaml:"max_idle_conns" default:"25"`
		MaxOpenConns    int           `yaml:"max_open_conns" default:"50"`
		StartTimeout    time.Duration `yaml:"start_timeout" default:"20s"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" default:"30s"`
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
		"sqlite3://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		conf.DBConfig.User,
		conf.DBConfig.Password,
		conf.DBConfig.Host,
		conf.DBConfig.Port,
		conf.DBConfig.Name,
		conf.DBConfig.Schema,
	)
}
