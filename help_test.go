package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type HelpSuite struct {
	Mock irc.ConnectionMock
}

var (
	_ = Suite(&HelpSuite{})
)

func (s *HelpSuite) SetUpTest(c *C) {
	cmd.Commands = make(map[string]*cmd.CustomCommand)
	s.Mock = irc.ConnectionMock{}
}

func (s *HelpSuite) TestHelpCommandNotFound(c *C) {
	channel := ""
	msg := []string{}
	s.Mock.NoticeFunc = func(target, message string) {
		channel = target
		msg = append(msg, message)
	}

	availableCommand := &cmd.CustomCommand{
		Cmd: "cmd",
	}
	cmd.RegisterCommand(availableCommand)

	command := &cmd.Cmd{
		Nick:   "nick",
		Prefix: "!",
	}
	Help(command, s.Mock)

	c.Check(channel, Equals, command.Nick)
	c.Check(msg, HasLen, 2)
	c.Check(msg[0], Equals, fmt.Sprintf(helpAboutCommand, command.Prefix))
	c.Check(msg[1], Equals, fmt.Sprintf(availableCommands, availableCommand.Cmd))
}
