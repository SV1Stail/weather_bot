package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var msg tgbotapi.MessageConfig
	fmt.Println("ok")

	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вся погода", "weatherALL"),        // 1
			tgbotapi.NewInlineKeyboardButtonData("Погода по часам", "weatherHours"), // 2
			tgbotapi.NewInlineKeyboardButtonData("Краткая погода", "weather"),       // 3
		),
	)

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! I am your bot.")
	msg.ReplyMarkup = numericKeyboard
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("ERROR start %v", err)
	}
	return nil
}
