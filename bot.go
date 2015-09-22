// Package bot provides a simple to use IRC bot
package bot

import (
	"math/rand"
	"time"
)

const (
	// CmdPrefix is the prefix used to identify a command.
	// !hello whould be identified as a command
	CmdPrefix = "!"
)

type connection interface {
	Privmsg(target, message string)
	GetNick() string
	Join(target string)
	Part(target string)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
