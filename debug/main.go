package main

import (
	"bufio"
	"fmt"
	"os"

	bot "github.com/bnfinet/go-chat-bot"
	_ "github.com/go-chat-bot/plugins-br/cnpj"
	_ "github.com/go-chat-bot/plugins-br/cotacao"
	_ "github.com/go-chat-bot/plugins-br/cpf"
	_ "github.com/go-chat-bot/plugins-br/dilma"
	_ "github.com/go-chat-bot/plugins-br/lula"
	_ "github.com/go-chat-bot/plugins-br/megasena"
	_ "github.com/go-chat-bot/plugins/9gag"
	_ "github.com/go-chat-bot/plugins/catfacts"
	_ "github.com/go-chat-bot/plugins/catgif"
	_ "github.com/go-chat-bot/plugins/chucknorris"
	_ "github.com/go-chat-bot/plugins/cmd"
	_ "github.com/go-chat-bot/plugins/crypto"
	_ "github.com/go-chat-bot/plugins/encoding"
	_ "github.com/go-chat-bot/plugins/example"
	_ "github.com/go-chat-bot/plugins/gif"
	_ "github.com/go-chat-bot/plugins/godoc"
	_ "github.com/go-chat-bot/plugins/guid"
	_ "github.com/go-chat-bot/plugins/jira"
	_ "github.com/go-chat-bot/plugins/puppet"
	_ "github.com/go-chat-bot/plugins/treta"
	_ "github.com/go-chat-bot/plugins/uptime"
	_ "github.com/go-chat-bot/plugins/url"
	_ "github.com/go-chat-bot/plugins/web"
)

func responseHandler(target string, message string, sender *bot.User) {
	if message == "" {
		return
	}
	fmt.Println(fmt.Sprintf("%s: %s", sender.Nick, message))
}

func main() {
	b := bot.New(&bot.Handlers{
		Response: responseHandler,
	},
		&bot.Config{
			Protocol: "debug",
			Server:   "debug",
		},
	)

	fmt.Println("Type a command or !help for available commands...")

	for {
		r := bufio.NewReader(os.Stdin)

		input, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		b.MessageReceived(
			&bot.ChannelData{
				Protocol:  "debug",
				Server:    "",
				Channel:   "console",
				IsPrivate: true,
			},
			&bot.Message{Text: input},
			&bot.User{ID: "id", RealName: "Debug Console", Nick: "bot", IsBot: false})
	}
}
