package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/via04/speech-writer-bot/pkg/converter"
	"github.com/via04/speech-writer-bot/pkg/utils"
	"github.com/via04/speech-writer-bot/pkg/witai"
)

// This value is maximum for free version of assemblyai
const MAX_CONNECTIONS = 5
const MAX_TIME = time.Second * 30

// Only one thread can print messages, so we need this structure to use it in channels
type message struct {
	upd  tgbotapi.Update
	text string
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic("nil: no bot")
	}
	ctx := context.Background()
	messageChan := make(chan message, MAX_CONNECTIONS)
	fmt.Fprintf(os.Stdout, "Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 3
	updates := bot.GetUpdatesChan(updateConfig)
	go func() {
		// poll messages channel and send answers to user
		for message := range messageChan {
			msg := tgbotapi.NewMessage(message.upd.Message.Chat.ID, message.text)
			bot.Send(msg)
		}
	}()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			// if we've got a string in a message
			switch update.Message.Text {
			case "/start":
				go func(update tgbotapi.Update) {
					// start goroutine to send help text to user
					helpText := "Hello!\n This bot is used to read voice messages." +
						"It uses Assembly AI to convert audio to text."
					newMesage := message{upd: update, text: helpText}
					messageChan <- newMesage
				}(update)
			default:
				go func(update tgbotapi.Update) {
					helpText := `Enter /start to start a bot or just send voice message`
					newMessage := message{upd: update, text: helpText}
					messageChan <- newMessage
				}(update)
			}
		}
		if update.Message.Voice != nil {
			// if we've got voice message
			go func(update tgbotapi.Update) {
				//start goroutine to take care of getting text from audio
				//create needed temp folder for files
				dataFolder, err := os.Getwd()
				fmt.Fprintln(os.Stdout, dataFolder)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				dataFolder = filepath.Join(dataFolder, "data")
				fmt.Fprintln(os.Stdout, dataFolder)
				os.Mkdir(dataFolder, os.ModePerm)
				fileId := update.Message.Voice.FileID
				fileInfo, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileId})
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				// get file name
				fname := fileInfo.FilePath
				parts := strings.Split(fname, "/")
				fname = parts[len(parts)-1]
				fmt.Fprintln(os.Stdout, fname)
				// get direct url for file
				url, err := bot.GetFileDirectURL(fileId)
				fmt.Fprintln(os.Stdout, url)
				ctxTimeout, cancel := context.WithTimeout(ctx, MAX_TIME)
				defer cancel()
				if err != nil {
					fmt.Fprint(os.Stderr, "voice url: something went wrong, cannot get url")
				}
				// download file to created folder
				fullNamein := filepath.Join(dataFolder, fname)
				fnameoutWExt, err := utils.GetNameNoExt(fname)
        if err != nil {
          fmt.Fprintln(os.Stderr, err.Error())
        }
				fullNameout := filepath.Join(dataFolder, fnameoutWExt+".wav")
				utils.Download(ctxTimeout, fullNamein, url)
				// convert file to work with witai to wav format
				conv, err := converter.New(fname, fnameoutWExt+".wav")
				defer conv.Purge(true)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				err = conv.ConvertToWav()
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				// get text from witai
				text, err := witai.GetText(ctxTimeout, fullNameout)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				// messageChan <- newMessage
				newMessage := message{upd: update, text: text}
				messageChan <- newMessage
				select {
				case <-ctxTimeout.Done():
					fmt.Fprint(os.Stderr, "connection timeout: no answer from assemblyai")
				case <-ctx.Done():
					panic("unexpected behaviour")
				default:
					// do not block
				}
			}(update)
		}
	}
}
