package commands

import (
	"fmt"
	"strings"
)

const (
	helpDescripton = "Description: %s"
	helpUsage = "Usage: %s%s %s"
)

func Help(cmd *Command) (msg string, err error) {
	notice := ""

	if cmd.FullArg != "" {
		if h,ok := helps[cmd.FullArg]; ok {
			if h.helpDescripton != "" {
				notice = fmt.Sprintf(helpDescripton, h.helpDescripton)
				irccon.Notice(cmd.Nick, notice)
			}
			if h.helpUse != "" {
				notice = fmt.Sprintf(helpUsage, cmd.Prefix, cmd.FullArg, h.helpUse)
				irccon.Notice(cmd.Nick, notice)
			}
		}
	} else {
		cmds := make([]string, 0)
		for k := range commands {
			cmds = append(cmds, k)
		}
		notice = fmt.Sprintf("%s: %s", availableCommands, strings.Join(cmds, ", "))
		irccon.Notice(cmd.Nick, notice)
	}

	return
}

func init() {
	RegisterCommand("help", Help)

	man := Manual{
		helpDescripton: "Show this help",
		helpUse: "[command]",
	}
	RegisterHelp("help", man)
}