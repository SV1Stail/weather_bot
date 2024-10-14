package commands

import (
	"log/slog"

	"github.com/SV1Stail/weather_bot/apiconnect"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Commands(update *tgbotapi.Update) {
	bot := apiconnect.GetBot()
	if update.Message.Command() == "start" {
		if err := Start(update, bot); err != nil {
			slog.Error("cant use command /start", "error", err)
			return
		}
	} else if update.Message.Command() == "help" {
		if err := Help(update, bot); err != nil {
			slog.Error("cant use command /help", "error", err)
			return
		}
	}
}
