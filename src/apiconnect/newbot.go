package apiconnect

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	newBot *tgbotapi.BotAPI
}

var weatherBot Bot

func MakeNewBot(token string) error {
	var err error
	weatherBot.newBot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println("FATAL ERROR MakeNewBot")
		return err
	}
	return nil
}

func GetBot() *tgbotapi.BotAPI {
	return weatherBot.newBot
}
