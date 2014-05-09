# go-bot

[![Build Status](https://travis-ci.org/fabioxgn/go-bot.png?branch=master)](https://travis-ci.org/fabioxgn/go-bot)

Simple IRC bot written in [Go][go] using [go-ircevent][go-ircevent] for the IRC connectivity.

[go]: golang.org
[go-ircevent]: https://github.com/thoj/go-ircevent

# Config is read from the environment:

    export IRC_SERVER=irc.freenode.org:7000
    export IRC_CHANNELS="#go-bot,#lightirc"
    export IRC_USER=go-bot
    export IRC_NICK=go-bot

To join channels with passwords just put the password after the channel name separated by a space:

    export IRC_CHANNELS="#mychannel mypassword,#go-bot"