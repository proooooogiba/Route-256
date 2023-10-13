package commands

type Command interface {
	GetCommandName() string
	GetArguments(args []string) error
	Execute()
}
