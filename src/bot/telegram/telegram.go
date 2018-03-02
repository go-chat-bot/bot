// Package telegram implements Telegram handlers for github.com/go-chat-bot/bot
package telegram

import (
	"log"
	"strconv"
	"strings"

	"bot/bot"

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
		target := &bot.ChannelData{
			Protocol:  "telegram",
			Server:    "telegram",
			Channel:   strconv.FormatInt(update.Message.Chat.ID, 10),
			IsPrivate: update.Message.Chat.IsPrivate()}
		name := []string{update.Message.From.FirstName, update.Message.From.LastName}
		message := &bot.Message{
			Text: update.Message.Text,
		}

		b.MessageReceived(target, message, &bot.User{
			ID:       strconv.Itoa(update.Message.From.ID),
			Nick:     update.Message.From.UserName,
			RealName: strings.Join(name, " ")})
	}
}
