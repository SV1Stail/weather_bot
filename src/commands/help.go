package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Help(update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here is help")
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("ERROR help %v", err)
	}
	return nil

}
