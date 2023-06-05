package cli

// Argument is the word, words or set of consecutive characters,
// that follow the command. If the command has no arguments,
// Label should be set as an empty string.
type Argument struct {

	// When the argument is parsed from the input, the argument
	// text must match the label in order to invoke the command.
	Label string `yaml:"label"`

	// If applicable, any options that are required for the
	// command.
	Options []Option `yaml:"options"`

	// The function performed when this command is invoked.
	// The options will be passed to this function as Flags.
	ExecFunc   string `yaml:"execFunc"`
	executable func(Flags) []byte

	// (optional) help message for this argument.
	HelpMsg string `yaml:"help"`

	// This function returns a help message for this argument.
	help func(Flags) []byte
}
