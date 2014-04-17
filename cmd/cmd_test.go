package cmd

import (
	"errors"
	"fmt"
	"github.com/fabioxgn/go-bot/irc"
	. "github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CmdSuite struct {
	Mock irc.ConnectionMock
}

var (
	_ = Suite(&CmdSuite{})
)

func (s *CmdSuite) SetUpTest(c *C) {
	Commands = make(map[string]*CustomCommand)
	s.Mock = irc.ConnectionMock{}
}

func (s *CmdSuite) TestRegisterCommand(c *C) {
	fn := func(cmd *Cmd) (string, error) { return "", nil }

	cmd := &CustomCommand{
		Cmd:     "cmd",
		CmdFunc: fn,
	}
	RegisterCommand(cmd)

	c.Check(Commands["cmd"], Equals, cmd)
}

func (s *CmdSuite) TestErrorExecutingCommand(c *C) {
	cmd := &Cmd{
		Command: "cmd",
		Channel: "#go-bot",
	}

	cmdError := errors.New("Error")
	fn := func(c *Cmd) (string, error) { return "", cmdError }

	RegisterCommand(&CustomCommand{
		Cmd:     cmd.Command,
		CmdFunc: fn,
	})

	msg := []string{}
	channel := ""
	s.Mock.PrivMsgFunc = func(c string, m string) {
		channel = c
		msg = append(msg, m)
	}

	err := HandleCmd(cmd, s.Mock)

	c.Check(err, Equals, cmdError)

	c.Check(channel, Equals, cmd.Channel)
	c.Check(msg[0], Equals, fmt.Sprintf(errorExecutingCommand, cmd.Command, cmdError.Error()))
}

func (s *CmdSuite) TestHandleExistingCommand(c *C) {
	cmd := &Cmd{
		Command: "cmd",
		Channel: "#go-bot",
	}
	expectedMsg := []string{"msg"}
	fn := func(c *Cmd) (string, error) {
		return expectedMsg[0], nil
	}

	RegisterCommand(&CustomCommand{
		Cmd:     "cmd",
		CmdFunc: fn,
	})

	printedMsg := []string{}
	channel := ""
	s.Mock.PrivMsgFunc = func(c string, m string) {
		channel = c
		printedMsg = append(printedMsg, m)
	}

	err := HandleCmd(cmd, s.Mock)

	c.Check(err, IsNil)
	c.Check(channel, Equals, cmd.Channel)
	c.Check(printedMsg[0], Equals, expectedMsg[0])
}

func (s *CmdSuite) TestHandleCommandNotFound(c *C) {
	cmd1 := &Cmd{
		Command: "cmd",
	}

	err := HandleCmd(cmd1, nil)

	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, fmt.Sprintf(commandNotAvailable, "cmd"))
}
