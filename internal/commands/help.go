package commands

import (
	"fmt"
)

type HelpCommand struct{}

func (fc *HelpCommand) GetArguments(args []string) error {
	return nil
}

func (fc *HelpCommand) GetCommandName() string {
	return "help"
}

func (fc *HelpCommand) Execute() {
	fmt.Println("Available commands:")
	fmt.Println("- help: Show available commands")
	fmt.Println("- spell: Spell a word")
	fmt.Println("- fmt: Format text as inserts a tab before each paragraph and puts a dot at the end of sentences")
}
