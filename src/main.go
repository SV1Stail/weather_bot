package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/SV1Stail/weather_bot/apiconnect"
	"github.com/SV1Stail/weather_bot/callbacks"
	"github.com/SV1Stail/weather_bot/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token = ""

func main() {
	logLevel := &slog.LevelVar{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	if os.Getenv("LOGGER") == "info" {
		logLevel.Set(slog.LevelInfo)
	} else {
		logLevel.Set(slog.LevelDebug)
	}
	slog.SetDefault(logger)

	err := apiconnect.MakeNewBot(token)
	if err != nil {
		log.Fatalf("Cant connect %v", err)
	}
	bot := apiconnect.GetBot()

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message != nil && update.Message.IsCommand() {
			commands.Commands(&update)
		} else if update.Message != nil && update.Message.Text != "" {
			callbacks.GiveWeatherByCoordinats(&update)
		}
		if update.CallbackQuery != nil {
			callbacks.Callback(&update)
		}

	}
}
