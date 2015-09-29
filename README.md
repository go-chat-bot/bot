# go-bot

[![Circle CI](https://circleci.com/gh/go-chat-bot/bot/tree/master.svg?style=svg)](https://circleci.com/gh/go-chat-bot/bot/tree/master) [![GoDoc](https://godoc.org/github.com/go-chat-bot/bot?status.png)](https://godoc.org/github.com/go-chat-bot/bot) [![Coverage Status](https://img.shields.io/coveralls/go-chat-bot/bot.svg)](https://coveralls.io/r/go-chat-bot/bot?branch=master)

Nice to meet you! I'm a IRC, Slack & Telegram bot written in [Go][go] using [go-ircevent][go-ircevent] for IRC connectivity, [nlopes/slack](https://github.com/nlopes/slack) for Slack and [Syfaro/telegram-bot-api](https://github.com/Syfaro/telegram-bot-api) for Telegram.

I can be deployed to [heroku][heroku] and used in [slack][slack] to overpower slackbot.

## #go-bot @ irc.freenode.org

I'm always hanging out as **go-bot** in the channel **#go-bot @ irc.freenode.org**

To see what I can do, type **!help** in the channel or send me !help in private message.

## My awesome commands

### Active commands

* **!gif**: Posts a random gif url from [giphy.com][giphy.com]. Try it with: **!gif cat**
* **!catgif**: Posts a random cat gif url from [thecatapi.com][thecatapi.com]
* **!godoc**: Searches packages in godoc.org. Try it with: **!godoc net/http**
* **!puppet**: Allows you to send messages through the bot: Try it with: **!puppet say #go-bot Hello!**

### Passive commands (triggers)

Passive commands are triggered by keywords, urls, etc.

These commands differ from the active commands as they are executed for every match in any text. Ex: The Chuck Norris command, replies with a Chuck Norris fact every time the words "chuck" or "norris" are mentioned on a channel.

* **url**: Detects url and gets it's title (very naive implementation, works sometimes)
* **catfacts**: Tells a random cat fact based on some cat keywords
* **jira**: Detects jira issue numbers and posts the url (necessary to configure the JIRA URL)
* **chucknorris**: Shows a random chuck norris quote every time the word "chuck" is mentioned

### Brazilian commands (pt-br)

I also have some brazilian commands which only makes sense to brazillians:

* **megasena**: Gera um número da megasena ou mostra o último resultado
* **cotacao**: Informa a cotação atual do Dólar e Euro
* **dilma** (passivo): Diz alguma frase da Dilma quando a palavra "dilma" é citada

### Example commands

If you wish to write your own commands, start with the 2 example commands, they are in the example directory.

## Joining and parting channels

### IRC

If you want me to join your channel, send me a private message with:

    !join #channel pass

If I'm boring you, just send a **!part** command in a channel I'm in.

### Slack

`!join` and `!part` does not work on slack (yet), so use the slack commands like `/invite` and `/remove` to invite the bot to the desired channels.

## Connecting to Slack

To deploy your go-bot to slack, you need to:

* [Create a new bot user](https://my.slack.com/services/new/bot) integration on slack and get your token
* Import the package `github.com/go-chat-bot/bot/slack`
* Import the commands you would like to use
* Call `slack.Run(token)`

Here is a full example reading the slack token from the `SLACK_TOKEN` env var:

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

## Connecting to IRC

To deploy your own go-bot, you need to:

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
		irc.Run(&bot.Config{
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

## Connecting to Telegram

To deploy your go-bot to slack, you need to:

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

## Deploying to heroku

To see an example project on how to deploy your bot, please see my own configuration:

- **IRC**: https://github.com/fabioxgn/go-bot-heroku
- **Slack**: https://github.com/fabioxgn/go-bot-slack

[go]: http://golang.org
[go-ircevent]: https://github.com/thoj/go-ircevent
[slack]: http://slack.com
[heroku]: http://heroku.com
[giphy.com]: http://giphy.com
[thecatapi.com]: http://thecatapi.com
