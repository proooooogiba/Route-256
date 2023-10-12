package main

import (
	"homework-3/internal/commands/invoker"
	"os"
)

func main() {
	commandInvoker, withCommand := invoker.NewCommandInvoker(os.Args)
	if withCommand {
		commandInvoker.ExecuteCommand()
		return
	}

}
