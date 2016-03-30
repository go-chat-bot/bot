// Package telegram implements Telegram handlers for github.com/go-chat-bot/bot
package telegram

import (
	"log"
	"strconv"

	"github.com/go-chat-bot/bot"
	"gopkg.in/telegram-bot-api.v3"
)

var (
	tg *tgbotapi.BotAPI
)

func responseHandler(target string, message string, sender *bot.User) {
	id, err := strconv.ParseInt(target, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	msg := tgbotapi.NewMessage(id, message)
	msg.ReplyToMessageID = msg.ReplyToMessageID

	tg.Send(msg)
}

// Run executes the bot and connects to Telegram using the provided token. Use the debug flag if you wish to see all trafic logged
func Run(token string, debug bool) {
	var err error
	tg, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	tg.Debug = debug

	log.Printf("Authorized on account %s", tg.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tg.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	b := bot.New(&bot.Handlers{
		Response: responseHandler,
	})
	b.Disable([]string{"url"})

	for update := range updates {
		target := strconv.FormatInt(update.Message.Chat.ID, 10)
		sender := strconv.Itoa(update.Message.From.ID)
		b.MessageReceived(target, update.Message.Text, &bot.User{Nick: sender})
	}
}
