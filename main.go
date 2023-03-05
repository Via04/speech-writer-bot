package main

import (
	"fmt"
	"os"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic("nil: no bot")
	}
	bot.Debug = true
	fmt.Fprintf(os.Stdout, "Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 3
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			// if we've got a string in a message
			switch update.Message.Text {
			case "/start":
				helpText := "Hello!\n This bot is used to read audio messages." +
					"It uses Assembly AI to convert audio to text."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
				bot.Send(msg)
			default:
				helpText := `Enter /start to start a bot or just send audio message`
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
				bot.Send(msg)
			}
		}
	}
}
