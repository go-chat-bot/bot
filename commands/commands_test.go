package commands

import (
	"fmt"
	. "github.com/fabioxgn/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	commands = make(map[string]CommandFunc)
}

func (s *MySuite) TestRegisterCommand(c *C) {
	command := "teste"
	ok := false
	fn := func(args []string) string { ok = true; return "" }
	RegisterCommand(command, fn)

	cmd := commands[command]
	cmd([]string{""})

	c.Check(ok, Equals, true)
}

func (s *MySuite) TestHandleExistingCommand(c *C) {
	arg1 := "echo"
	cmdString := "cmd"
	cmd := &Command{
		Command: cmdString,
		Args:    []string{arg1},
	}
	fn := func(args []string) string { return ArgsToString(args) }

	RegisterCommand(cmdString, fn)

	msg := ""
	privMsgFunc := func(c string, m string) {
		msg = m
	}

	HandleCmd(cmd, "", privMsgFunc)

	c.Assert(msg, Equals, arg1)
}

func (s *MySuite) TestNoCommandsAvailable(c *C) {
	cmd := &Command{Command: "cmd"}

	msg := []string{}
	fn := func(c string, m string) {
		msg = append(msg, m)
	}

	HandleCmd(cmd, "", fn)

	c.Assert(msg, HasLen, 2)
	c.Check(msg[1], Equals, noCommandsAvailable)
}

func (s *MySuite) TestHandleCommandNotFound(c *C) {
	channel := ""
	msg := []string{}
	fn := func(c string, m string) {
		channel = c
		msg = append(msg, m)
	}

	cmd := &Command{}
	cmd.Command = "allyourbase"

	expectedChannel := "#go-bot"
	expectedMsg := fmt.Sprintf(commandNotAvailable, cmd.Command)

	HandleCmd(cmd, expectedChannel, fn)

	c.Check(channel, Equals, expectedChannel)
	c.Assert(msg, HasLen, 2)
	c.Check(msg[0], Equals, expectedMsg)
}
