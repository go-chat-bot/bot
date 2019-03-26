// Copyright (c) 2012 VMware, Inc. - https://github.com/cloudfoundry/gosigar/blob/master/examples/uptime.go

package uptime

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"time"
	"github.com/go-chat-bot/bot"
)

func uptime(command *bot.Cmd) (msg string, err error) {
	uptime := sigar.Uptime{}
	uptime.Get()
	avg := sigar.LoadAverage{}
	avg.Get()
	msg = fmt.Sprintf("%s up %s load average: %.2f, %.2f, %.2f\n", time.Now().Format("15:04:05"),	uptime.Format(), avg.One, avg.Five, avg.Fifteen)
	return
}

func init() {
	bot.RegisterCommand(
		"uptime",
		"Sends the uptime of your server to you on the channel.",
		"",
		uptime)
}
