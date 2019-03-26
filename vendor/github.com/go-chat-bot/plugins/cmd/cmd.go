package cmd

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/go-chat-bot/bot"
)

var (
	disableCmds   = map[string]bool{"shutdown": true, "reboot": true, "init": true, "rm": true, "top": true, "htop": true, "iotop": true}
	errDisableCmd = errors.New("command is disabled")
)

func cmd(command *bot.Cmd) (string, error) {
	if _, ok := disableCmds[command.Args[0]]; ok {
		return "", errDisableCmd
	}
	cmd := exec.Command("/bin/bash", "-c", command.RawArgs)
	data, err := cmd.CombinedOutput()
	return string(data), err
}

func cmdV3(command *bot.Cmd) (result bot.CmdResultV3, err error) {
	result = bot.CmdResultV3{Message: make(chan string), Done: make(chan bool)}
	if _, ok := disableCmds[command.Args[0]]; ok {
		err = errDisableCmd
		return
	}

	cmd := exec.Command("/bin/bash", "-c", command.RawArgs)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err = cmd.Start()
	if err != nil {
		return
	}
	done := false
	go func() {
		cmd.Wait()
		done = true
		result.Done <- true
	}()
	go func() {
		for {
			line, _ := b.ReadString('\n')
			if line != "" {
				result.Message <- line
			}
			if done {
				break
			}

		}
	}()
	return
}

func init() {
	bot.RegisterCommand(
		"cmd",
		"run cmd on system",
		"pwd",
		cmd)
	bot.RegisterCommandV3(
		"cmdv3",
		"run cmd on system",
		"pwd",
		cmdV3)
}
