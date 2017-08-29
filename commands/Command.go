package commands

type Command interface {
	CommandIdentifier() string
	CommandDescription() string
	execute(arguments string)
}

var commandList = []commandCtor{}

type commandCtor func() Command

func Register(newfund commandCtor) {
	commandList = append(commandList, newfund)
}

func findCommand(command string) (a Command) {
	for _, item := range commandList {
		a = item()
		if (a.CommandIdentifier() == command) {
			return a
		}
	};
	return nil
}

func ExecuteCommand(command,arguments string) {
	findCommand(command).execute(arguments)
}
