package bot

import (
	"errors"
	"fmt"
	check "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type CmdSuite struct {
	Mock ConnectionMock
}

var (
	_ = check.Suite(&CmdSuite{})
)

func (s *CmdSuite) SetUpTest(c *check.C) {
	Commands = make(map[string]*CustomCommand)
	s.Mock = ConnectionMock{}
}

func (s *CmdSuite) TestRegisterCommand(c *check.C) {
	fn := func(cmd *Cmd) (string, error) { return "", nil }

	cmd := &CustomCommand{
		Cmd:     "cmd",
		CmdFunc: fn,
	}
	RegisterCommand(cmd)

	c.Check(Commands["cmd"], check.Equals, cmd)
}

func (s *CmdSuite) TestErrorExecutingCommand(c *check.C) {
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

	c.Check(err, check.Equals, cmdError)

	c.Check(channel, check.Equals, cmd.Channel)
	c.Check(msg[0], check.Equals, fmt.Sprintf(errorExecutingCommand, cmd.Command, cmdError.Error()))
}

func (s *CmdSuite) TestHandleExistingCommand(c *check.C) {
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

	c.Check(err, check.IsNil)
	c.Check(channel, check.Equals, cmd.Channel)
	c.Check(printedMsg[0], check.Equals, expectedMsg[0])
}

func (s *CmdSuite) TestHandleCommandNotFound(c *check.C) {
	cmd1 := &Cmd{
		Command: "cmd",
	}

	err := HandleCmd(cmd1, nil)

	c.Check(err, check.NotNil)
	c.Check(err.Error(), check.Equals, fmt.Sprintf(commandNotAvailable, "cmd"))
}
