package invoker

import (
	"fmt"
	"homework-3/internal/commands"
)

type CommandInvoker struct {
	Args     []string
	commands []commands.Command
}

func NewCommandInvoker(Args []string) (CommandInvoker, bool) {
	if len(Args) < 2 {
		return CommandInvoker{}, false
	}

	invoker := CommandInvoker{Args: Args}
	invoker.AddCommand(&commands.HelpCommand{CommandName: "help"})
	invoker.AddCommand(&commands.SpellCommand{CommandName: "spell"})
	invoker.AddCommand(&commands.FormatCommand{CommandName: "fmt"})

	return invoker, true
}

func (ci *CommandInvoker) AddCommand(command commands.Command) {
	ci.commands = append(ci.commands, command)
}

func (ci *CommandInvoker) ExecuteCommand() {
	commandName := ci.Args[1]
	for _, command := range ci.commands {
		if command.GetCommandName() == commandName {
			err := command.GetArguments(ci.Args[2:])
			if err != nil {
				fmt.Println(err)
				return
			}
			command.Execute()
			return
		}
	}

	fmt.Println("Command isn't supported")
}
