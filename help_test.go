package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
	. "github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	Mock irc.ConnectionMock
}

var (
	_ = Suite(&MySuite{})
)

func (s *MySuite) SetUpTest(c *C) {
	cmd.Commands = make(map[string]*cmd.CustomCommand)
	s.Mock = irc.ConnectionMock{}
}

func (s *MySuite) TestHelpCommandNotFound(c *C) {
	channel := ""
	msg := ""
	s.Mock.NoticeFunc = func(target, message string) {
		channel = target
		msg = message
	}

	availableCommand := &cmd.CustomCommand{
		Cmd: "cmd",
	}
	cmd.RegisterCommand(availableCommand)

	command := &cmd.Cmd{
		Nick: "nick",
	}
	Help(command, s.Mock)

	c.Check(channel, Equals, command.Nick)
	c.Check(msg, Equals, fmt.Sprintf(availableCommands, availableCommand.Cmd))

}
