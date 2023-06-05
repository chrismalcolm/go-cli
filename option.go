package cli

// Option is a character, set of consecutive characters,
// or a word that follows the command and any arguments.
// Options are preceded by an dash (–).
type Option struct {

	// Name of the variable that will be used as a key for
	// a flag mapping when used in a action function.
	Label string `yaml:"label"`

	// Short name, single dash (–) followed by a signle
	// character.
	Short string `yaml:"short"`

	// Long name, double dash (--) followed by a
	// descriptive name.
	Long string `yaml:"long"`

	// (optional) if this option requires a variable,
	// it should be defined here.
	Variable *Variable `yaml:"variable"`

	// (optional) help message for this option.
	HelpMsg string `yaml:"help"`
}
