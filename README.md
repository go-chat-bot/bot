# go-bot

[![Build Status](https://travis-ci.org/fabioxgn/go-bot.png?branch=master)](https://travis-ci.org/fabioxgn/go-bot) [![GoDoc](https://godoc.org/github.com/fabioxgn/go-bot?status.png)](https://godoc.org/github.com/fabioxgn/go-bot) [![Coverage Status](https://img.shields.io/coveralls/fabioxgn/go-bot.svg)](https://coveralls.io/r/fabioxgn/go-bot?branch=master)

Nice to meet you! I'm a IRC bot written in [Go][go] using [go-ircevent][go-ircevent] for IRC connectivity.

I can be deployed to [heroku][heroku] and used in [slack][slack] to overpower slackbot.

## #go-bot @ irc.freenode.org

I'm always hanging out as **go-bot** in the channel **#go-bot @ irc.freenode.org**

To see what I can do, type **!help** in the channel or send me a private message.

## My awesome commands

### Active commands

* **!gif**: Posts a random gif url from [giphy.com][giphy.com]. Try it with: **!gif cat**
* **!catgif**: Posts a random cat gif url from [thecatapi.com][thecatapi.com]
* **!godoc**: Searches packages in godoc.org. Try it with: **!godoc net/http**
* **!puppet**: Allows you to send messages through the bot: Try it with: **!puppet say #go-bot Hello!**

### Passive commands (triggers)

I also have some commands, which are triggered by keywords, urls, etc

* **url**: Detects url and gets it's title (very naive implementation, works sometimes)
* **catfacts**: Tells a random cat fact based on some cat keywords
* **jira**: Detects jira issue numbers and posts the url (necessary to configure the JIRA URL)
* **chucknorris**: Shows a random chuck norris quote every time the word "chuck" is mentioned

### Brazilian commands (pt-br)

I also have some brazilian commands which only apply to Brazil:

* **megasena**: Gera um número da megasena ou mostra o último resultado
* **cotacao**: Informa a cotação atual do Dólar e Euro
* **dilma** (passive): Diz alguma frase da Dilma quando a palavra "dilma" é citada

### Example commands

If you wish to write your own commands, start with the 2 example commands, they are in the example directory.

## Joining and parting channels

If you want me to join your channel, send me a private message with:

    !join #channel pass

If I'm boring you, just send a **!part** command in a channel I'm in.

## Deploying to heroku

To see an example project on how to deploy it to heroku, please see my own configuration:

https://github.com/fabioxgn/go-bot-heroku

## Deploying your own clone of me

To deploy your own go-bot, you need to:

* Import the package bot
* Import the commands you would like to use
* Fill the Config struct
* Call Bot.Run(config)

Here is a full example:
```Go
	import (
		"github.com/fabioxgn/go-bot"
		_ "github.com/fabioxgn/go-bot/commands/catfacts"
		_ "github.com/fabioxgn/go-bot/commands/catgif"
		_ "github.com/fabioxgn/go-bot/commands/chucknorris"
		// Import all the commands you wish to use
		"log"
		"strings"
	)

	func main() {
		bot.Run(&bot.Config{
			Server:   os.Getenv("IRC_SERVER"),
			Channels: strings.Split(os.Getenv("IRC_CHANNELS"), ","),
			User:     os.Getenv("IRC_USER"),
			Nick:     os.Getenv("IRC_NICK"),
			Password: os.Getenv("IRC_PASSWORD"),
			UseTLS:   true,
			Debug:    os.Getenv("DEBUG") != "",}
		)
	}
```

To join channels with passwords just put the password after the channel name separated by a space:

    Channels: []string{"#mychannel mypassword", "#go-bot"}

[go]: http://golang.org
[go-ircevent]: https://github.com/thoj/go-ircevent
[slack]: http://slack.com
[heroku]: http://heroku.com
[giphy.com]: http://giphy.com
[thecatapi.com]: http://thecatapi.com
