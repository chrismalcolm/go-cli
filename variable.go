package cli

// Variable is any set of consecutive characters or word
// that follows an option.
type Variable struct {

	// Name of the variable that will be used as a key for
	// a variable mapping when used in a action function.
	// It will also be used as the variable placeholder in
	// the help messages.
	Label string `yaml:"label"`

	// Whether the option for this variable is required
	// for the command.
	Required bool `yaml:"required"`

	// (optional) The default value for the variable
	Default string `yaml:"default"`
}
