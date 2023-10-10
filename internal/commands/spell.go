package commands

import (
	"fmt"
	"strings"
)

type SpellCommand struct {
	Word string
}

func (sc *SpellCommand) Execute() {
	fmt.Print("Spelling: ")

	chars := strings.Join(strings.Split(sc.Word, ""), " ")

	fmt.Print(chars)
}
