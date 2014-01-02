# go-bot

[![Build Status](https://travis-ci.org/fabioxgn/go-bot.png?branch=master)](https://travis-ci.org/fabioxgn/go-bot)

Simple IRC bot written in [Go][go] using [go-ircevent][go-ircevent] for the IRC connectivity.

[go]: golang.org
[go-ircevent]: https://github.com/thoj/go-ircevent

# Sample config

    {  
	    "Server": "irc.freenode.net:7000",
	    "Channels": ["#go-bot"],
	    "User": "go-bot",		
	    "Nick": "go-bot",
	    "Cmd": "!go-bot",
	    "UseTLS": true
    }

# TODO

- Connect to multiple channels
- Run commands in parallel
