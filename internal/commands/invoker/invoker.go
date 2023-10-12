package invoker

import (
	"fmt"
	"homework-3/internal/commands"
	"os"
)

type CommandInvoker struct {
	Args []string
}

func NewCommandInvoker(Args []string) (CommandInvoker, bool) {
	if len(Args) < 2 {
		return CommandInvoker{}, false
	}
	return CommandInvoker{Args: Args}, true
}

func (ci *CommandInvoker) ExecuteCommand() {
	var command commands.Command
	switch os.Args[1] {
	case "help":
		command = &commands.HelpCommand{}
	case "spell":
		if len(os.Args) < 3 {
			fmt.Println("Please specify a word to spell")
			return
		}
		command = &commands.SpellCommand{Word: os.Args[2]}
	case "fmt":
		if len(os.Args) < 3 {
			fmt.Println("Please specify a .txt document to format")
			return
		}
		command = &commands.FormatCommand{FileName: os.Args[2]}
	default:
		fmt.Println("Command isn't supported")
		return
	}
	command.Execute()
}
