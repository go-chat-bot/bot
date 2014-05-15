# go-bot

[![Build Status](https://travis-ci.org/fabioxgn/go-bot.png?branch=master)](https://travis-ci.org/fabioxgn/go-bot)

IRC bot written in [Go][go] using [go-ircevent][go-ircevent] which can be deployed to heroku and used in slack.

[go]: golang.org
[go-ircevent]: https://github.com/thoj/go-ircevent

# Deploying to heroku

To see an example projeto on how to deploy it to heroku, please see my own bot:

https://github.com/fabioxgn/go-bot-heroku

To join channels with passwords just put the password after the channel name separated by a space:

    Channels: "#mychannel mypassword,#go-bot",