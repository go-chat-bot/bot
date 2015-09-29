package telegram

import (
	"log"
	"strconv"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/go-chat-bot/bot"
)

var (
	tg *tgbotapi.BotAPI
)

func responseHandler(target, message, sender string) {
	id, err := strconv.Atoi(target)
	if err != nil {
		log.Println(err)
		return
	}
	msg := tgbotapi.NewMessage(id, message)
	msg.ReplyToMessageID = msg.ReplyToMessageID

	tg.SendMessage(msg)
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

	err = tg.UpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	bot.New(&bot.Handlers{
		Response: responseHandler,
	})

	for update := range tg.Updates {
		target := strconv.Itoa(update.Message.Chat.ID)
		sender := strconv.Itoa(update.Message.From.ID)
		bot.MessageReceived(target, update.Message.Text, sender)
	}
}
