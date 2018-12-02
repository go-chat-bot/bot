// Package telegram implements Telegram handlers for github.com/go-chat-bot/bot
package telegram

import (
	"log"
	"strconv"
	"strings"

	bot "github.com/bnfinet/go-chat-bot"
	tgbotapi "gopkg.in/telegram-bot-api.v3"
)

var (
	tg *tgbotapi.BotAPI
)

const (
	protocol = "telegram"
	server   = "telegram"
)

func responseHandler(target string, message string, sender *bot.User) {
	id, err := strconv.ParseInt(target, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	msg := tgbotapi.NewMessage(id, message)

	tg.Send(msg)
}

// Run executes the bot and connects to Telegram using the provided token. Use the debug flag if you wish to see all traffic logged
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
	}, &bot.Config{
		Protocol: protocol,
		Server:   server,
	},
	)
	// b.Protocol = protocol
	// b.Server = server

	b.Disable([]string{"url"})

	for update := range updates {
		target := &bot.ChannelData{
			Protocol:  protocol,
			Server:    server,
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
