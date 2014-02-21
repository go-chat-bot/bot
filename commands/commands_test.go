package commands

import (
	// "fmt"
	. "github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }
type MySuite struct{}
var (
	_ = Suite(&MySuite{})
   defaultCmd = "mycommand"
   defaultArg = "arg1 arg2"
	defaultCommand = &Command{
		Command: defaultCmd,
		FullArg: defaultArg,
	}
   defaultCommandFn = func(cmd *Command) (string, error) {
		return cmd.FullArg, nil
	}
)

func (s *MySuite) SetUpTest(c *C) {
	c.Skip("skipping these tests temporarily")
	// helps = make(map[string]Manual)
	commands = make(map[string]CommandFunc)
}

func (s *MySuite) TestRegisterCommand(c *C) {
	RegisterCommand(defaultCmd, defaultCommandFn)
	cmdFn := commands[defaultCmd]
	c.Check(cmdFn, NotNil)
}

func (s *MySuite) TestHandleExistingCommand(c *C) {
	RegisterCommand(defaultCmd, defaultCommandFn)
	resultError := HandleCmd(defaultCommand, irccon)
	c.Check(resultError, IsNil)
}

// func (s *MySuite) TestNoCommandsAvailable(c *C) {
// 	cmd := &Command{Command: "cmd"}

// 	msg := []string{}
// 	// fn := func(c string, m string) {
// 	// 	msg = append(msg, m)
// 	// }

// 	HandleCmd(cmd, irccon)

// 	c.Assert(msg, HasLen, 2)
// 	c.Check(msg[1], Equals, noCommandsAvailable)
// }

// func (s *MySuite) TestHandleCommandNotFound(c *C) {
// 	channel := ""
// 	msg := []string{}
// 	// fn := func(c string, m string) {
// 	// 	channel = c
// 	// 	msg = append(msg, m)
// 	// }

// 	cmd := &Command{}
// 	cmd.Command = "allyourbase"

// 	expectedChannel := "#go-bot"
// 	expectedMsg := fmt.Sprintf(commandNotAvailable, cmd.Command)

// 	HandleCmd(cmd, irccon)

// 	c.Check(channel, Equals, expectedChannel)
// 	c.Assert(msg, HasLen, 2)
// 	c.Check(msg[0], Equals, expectedMsg)
// }

// func (s *MySuite) TestHandleInvalidCommand(c *C) {
// 	cmd1 := "cmd1"
// 	cmd2 := "cmd2"
// 	RegisterCommand(cmd1, nil)
// 	RegisterCommand(cmd2, nil)

// 	cmd3 := "cmd3"
// 	cmd := &Command{Command: cmd3}

// 	msg := []string{}
// 	// fn := func(c string, m string) {
// 	// 	msg = append(msg, m)
// 	// }

// 	HandleCmd(cmd, irccon)

// 	c.Check(msg, HasLen, 2)
// 	c.Check(msg[0], Equals, fmt.Sprintf(commandNotAvailable, cmd3))
// 	c.Check(msg[1], Equals, fmt.Sprintf("%s: %s, %s", availableCommands, cmd1, cmd2))
// }
