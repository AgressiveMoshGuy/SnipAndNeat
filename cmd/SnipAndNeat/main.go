package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	application "SnipAndNeat/app"
	"SnipAndNeat/app/config"

	_ "github.com/mattn/go-sqlite3"
	// Import to keep it in go.mod.
	// _ "github.com/ogen-go/ogen/cmd/ogen"
)

func init() {
	flag.String("config", "./config.yaml", "The configuration file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

type InlineKeyboardButton struct {
	Text                         string  `json:"text"`
	URL                          *string `json:"url,omitempty"`                              // optional
	CallbackData                 *string `json:"callback_data,omitempty"`                    // optional
	SwitchInlineQuery            *string `json:"switch_inline_query,omitempty"`              // optional
	SwitchInlineQueryCurrentChat *string `json:"switch_inline_query_current_chat,omitempty"` // optional
}

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --clean --package oas --target ../../generated ../../openapi.yml
func main() {
	cfg, err := config.CreateFromFile()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create config")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Создание экземпляра приложения
	app := application.New(cfg)

	// Контекст с таймаутом для старта приложения
	startCtx, startCancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer startCancel()
	// Старт приложения
	if err := app.Start(startCtx); err != nil {
		log.Fatal().Msg("cannot start application")
	}

	log.Info().Msg("application started")

	quitCh := make(chan os.Signal, 1)

	// Прослушивание сигнала остановки сервиса, graceful shutdown
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	// Контекст с таймаутов для оставноки компонентов
	stopCtx, stopCancel := context.WithTimeout(context.Background(), cfg.StopTimeout)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		files, err := os.ReadDir(".")
		if err != nil {
			log.Fatal().Err(err).Msg("cannot read current directory")
		}
		for _, file := range files {
			if file.Name() != "main.go" {
				if err := os.Remove(file.Name()); err != nil {
					log.Fatal().Err(err).Msgf("cannot remove file %q", file.Name())
				}
			}
		}
		log.Fatal().Msg("cannot stop application")
	}

	log.Info().Msg("service is down")

}
