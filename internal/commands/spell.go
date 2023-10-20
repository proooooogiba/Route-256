package commands

import (
	"errors"
	"fmt"
	"strings"
)

type SpellCommand struct {
	Word string
}

func (sc *SpellCommand) GetArguments(args []string) error {
	if (len(args)) < 1 {
		return errors.New("please specify a word to spell")
	}
	sc.Word = args[0]
	return nil
}

func (sc *SpellCommand) GetCommandName() string {
	return "spell"
}

func (sc *SpellCommand) Execute() {
	fmt.Print("Spelling: ")
	chars := strings.Join(strings.Split(sc.Word, ""), " ")

	fmt.Print(chars)
}
