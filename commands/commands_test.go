package commands

import (
	"errors"
	"fmt"
	. "github.com/motain/gocheck"
	"reflect"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	Mock IRCConnectionMock
}

var (
	_ = Suite(&MySuite{})
)

func (s *MySuite) SetUpTest(c *C) {
	//c.Skip("skipping these tests temporarily")
	// helps = make(map[string]Manual)
	commands = make(map[string]CommandFunc)
	s.Mock = IRCConnectionMock{}
}

func (s *MySuite) TestRegisterCommand(c *C) {
	fn := func(cmd *Command) (string, error) { return "", nil }

	RegisterCommand("cmd", fn)

	c.Check(reflect.ValueOf(commands["cmd"]).Pointer(), Equals, reflect.ValueOf(fn).Pointer())
}

func (s *MySuite) TestErrorExecutingCommand(c *C) {
	cmd := &Command{
		Command: "cmd",
		Channel: "#go-bot",
	}

	cmdError := errors.New("Error")
	RegisterCommand(cmd.Command, func(c *Command) (string, error) {
		return "", cmdError

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

func (s *MySuite) TestHandleExistingCommand(c *C) {
	cmd := &Command{
		Command: "cmd",
		Channel: "#go-bot",
	}
	expectedMsg := []string{"msg"}

	cmdFuncCalled := false
	RegisterCommand(cmd.Command, func(c *Command) (string, error) {
		cmdFuncCalled = true
		return expectedMsg[0], nil
	})

	printedMsg := []string{}
	channel := ""
	s.Mock.PrivMsgFunc = func(c string, m string) {
		channel = c
		printedMsg = append(printedMsg, m)
	}

	err := HandleCmd(cmd, s.Mock)

	c.Check(err, IsNil)
	c.Check(cmdFuncCalled, Equals, true)

	c.Check(channel, Equals, cmd.Channel)
	c.Check(printedMsg[0], Equals, expectedMsg[0])
}

func (s *MySuite) TestNoCommandsAvailable(c *C) {
	c.Skip("TestNoCommandsAvailable")
	cmd := &Command{Command: "cmd"}

	msg := []string{}
	s.Mock.PrivMsgFunc = func(target, message string) {
		msg = append(msg, message)
	}

	HandleCmd(cmd, s.Mock)

	c.Assert(msg, HasLen, 2)
	c.Check(msg[1], Equals, noCommandsAvailable)
}

func (s *MySuite) TestHandleCommandNotFound(c *C) {
	cmd1 := &Command{
		Command: "cmd",
	}

	err := HandleCmd(cmd1, nil)

	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, fmt.Sprintf(commandNotAvailable, "cmd"))
}
