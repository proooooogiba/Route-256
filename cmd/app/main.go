package main

import (
	"homework-3/internal/commands"
	"homework-3/internal/commands/invoker"
	"os"
)

func main() {
	commandList := []commands.Command{
		&commands.HelpCommand{},
		&commands.FormatCommand{},
		&commands.SpellCommand{},
	}

	commandInvoker, withCommand := invoker.NewCommandInvoker(os.Args, commandList)
	if withCommand {
		commandInvoker.ExecuteCommand()
		return
	}
}
