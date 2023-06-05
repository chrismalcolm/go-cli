package cli

import (
	"fmt"
)

// Command is the first word or set of consecutive characters.
type Command struct {

	// When the command is parsed from the input, the command
	// text must match the label in order to invoke the command.
	Label string `yaml:"label"`

	// Any arguments that are required for the command.
	// For a command without any arguments, the argument label
	// should be an empty string.
	Arguments []Argument `yaml:"arguments"`

	// This function returns a help message for this command.
	help func(Flags) []byte
}

// createExecutable creates a placeholder executable method
func (cmd Command) createExecutable(arg Argument) func(Flags) []byte {
	return func(Flags) []byte {
		return []byte(fmt.Sprintf("\"%s\" is not configured\n", arg.ExecFunc))
	}
}
