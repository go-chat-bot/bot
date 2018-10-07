# go-bot

[![Circle CI](https://circleci.com/gh/go-chat-bot/bot/tree/master.svg?style=svg)](https://circleci.com/gh/go-chat-bot/bot/tree/master) [![GoDoc](https://godoc.org/github.com/go-chat-bot/bot?status.png)](https://godoc.org/github.com/go-chat-bot/bot) [![Coverage Status](https://coveralls.io/repos/github/go-chat-bot/bot/badge.svg?branch=master)](https://coveralls.io/github/go-chat-bot/bot?branch=master) ![Go report](https://goreportcard.com/badge/github.com/go-chat-bot/bot)

IRC, Slack & Telegram bot written in [Go][go] using [go-ircevent][go-ircevent] for IRC connectivity, [nlopes/slack](https://github.com/nlopes/slack) for Slack and [Syfaro/telegram-bot-api](https://github.com/Syfaro/telegram-bot-api) for Telegram.

![2016-01-17 11 21 38 036](https://cloud.githubusercontent.com/assets/1084729/12377689/5bf7d5f2-bd0d-11e5-87d9-525481f01c3a.gif)

## Plugins

Please see the [plugins repository](https://github.com/go-chat-bot/plugins) for a complete list of plugins.

You can also write your own, it's really simple.

## Compiling and testing the bot and plugins (Debug)

This project uses the new [Go 1.11 modules](https://github.com/golang/go/wiki/Modules) if you have Go 1.11 installed, just clone the project and follow the instructions bellow, when you build Go will automatically download all dependencies.

To test the bot, use the [debug](https://github.com/go-chat-bot/bot/tree/master/debug) console app.

- Clone this repository or use `go get github.com/go-chat-bot`
- Build everything: `go build ./...`
- Build and execute the debug app:
  -  `cd debug`
  -  `go build`
  -  `./debug`
- This will open a console where you can type commands
- Type `!help` to see the list of available commands

### Testing your plugin

- Add your plugin to `debug/main.go` import list
- Build the debug app
- Execute it and test with the interactive console

## Protocols

### Slack

To deploy your go-bot to Slack, you need to:

* [Create a new bot user](https://my.slack.com/services/new/bot) integration on Slack and get your token
* Import the package `github.com/go-chat-bot/bot/slack`
* Import the commands you would like to use
* Call `slack.Run(token)`

Here is a full example reading the Slack token from the `SLACK_TOKEN` env var:

```Go
package main

import (
    "os"

    "github.com/go-chat-bot/bot/slack"
    _ "github.com/go-chat-bot/plugins/catfacts"
    _ "github.com/go-chat-bot/plugins/catgif"
    _ "github.com/go-chat-bot/plugins/chucknorris"
    // Import all the commands you wish to use
)

func main() {
    slack.Run(os.Getenv("SLACK_TOKEN"))
}
```

### IRC

To deploy your own go-bot to IRC, you need to:

* Import the package `github.com/go-chat-bot/bot/irc`
* Import the commands you would like to use
* Fill the Config struct
* Call `irc.Run(config)`

Here is a full example:
```Go
package main

import (
	"github.com/go-chat-bot/bot/irc"
	_ "github.com/go-chat-bot/plugins/catfacts"
	_ "github.com/go-chat-bot/plugins/catgif"
	_ "github.com/go-chat-bot/plugins/chucknorris"
	// Import all the commands you wish to use
	"os"
	"strings"
)

func main() {
	irc.Run(&irc.Config{
		Server:   os.Getenv("IRC_SERVER"),
		Channels: strings.Split(os.Getenv("IRC_CHANNELS"), ","),
		User:     os.Getenv("IRC_USER"),
		Nick:     os.Getenv("IRC_NICK"),
		Password: os.Getenv("IRC_PASSWORD"),
		UseTLS:   true,
		Debug:    os.Getenv("DEBUG") != "",})
}
```

To join channels with passwords just put the password after the channel name separated by a space:

    Channels: []string{"#mychannel mypassword", "#go-bot"}

### Telegram

To deploy your go-bot to Telegram, you need to:

* Follow Telegram instructions to [create a new bot user](https://core.telegram.org/bots#3-how-do-i-create-a-bot) and get your token
* Import the package `github.com/go-chat-bot/bot/telegram`
* Import the commands you would like to use
* Call `telegram.Run(token, debug)`

Here is a full example reading the telegram token from the `TELEGRAM_TOKEN` env var:

```Go
package main

import (
    "os"

    "github.com/go-chat-bot/bot/telegram"
    _ "github.com/go-chat-bot/plugins/catfacts"
    _ "github.com/go-chat-bot/plugins/catgif"
    _ "github.com/go-chat-bot/plugins/chucknorris"
    // Import all the commands you wish to use
)

func main() {
    telegram.Run(os.Getenv("TELEGRAM_TOKEN"), os.Getenv("DEBUG") != "")
}
```

### Rocket.chat

To deploy your go-bot to Rocket.chat, you need to:

* Import the package `github.com/go-chat-bot/bot/rocket`
* Import the commands you would like to use
* Call `rocket.Run(config)`

Here is a full example:

```Go
package main

import (
	"os"

	"github.com/go-chat-bot/bot/rocket"
	_ "github.com/go-chat-bot/plugins/godoc"
	_ "github.com/go-chat-bot/plugins/catfacts"
	_ "github.com/go-chat-bot/plugins/catgif"
	_ "github.com/go-chat-bot/plugins/chucknorris"
)

func main() {
	config := &rocket.Config{
		Server:   os.Getenv("ROCKET_SERVER"),
		Port:     os.Getenv("ROCKET_PORT"),
		User:     os.Getenv("ROCKET_USER"),
		Email:    os.Getenv("ROCKET_EMAIL"),
		Password: os.Getenv("ROCKET_PASSWORD"),
		UseTLS:   false,
		Debug:    os.Getenv("DEBUG") != "",
	}
	rocket.Run(config)
}
```

## Deploying your own bot

To see an example project on how to deploy your bot, please see my own configuration:

- **IRC**: https://github.com/fabioxgn/go-bot-heroku
- **Slack**: https://github.com/fabioxgn/go-bot-slack

[go]: http://golang.org
[go-ircevent]: https://github.com/thoj/go-ircevent
[slack]: http://slack.com
[giphy.com]: http://giphy.com
[thecatapi.com]: http://thecatapi.com
