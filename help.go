package cli

import (
	"fmt"
	"strings"
)

// createHelp is a function for generating the global help function
func (config Config) createHelp() func(Flags) []byte {
	return func(flags Flags) (output []byte) {
		output = make([]byte, 0)
		for _, command := range config.Commands {
			output = append(output, command.help(flags)...)
		}
		return output
	}
}

// createHelp is a function for generating the command help function
func (cmd Command) createHelp() func(Flags) []byte {
	return func(_ Flags) []byte {
		return []byte(cmd.helpCmd())
	}
}

// createHelp is a function for generating the argument help function
func (cmd Command) createArgHelp(arg Argument) func(Flags) []byte {
	return func(_ Flags) []byte {
		return []byte(cmd.helpArg(arg))
	}
}

// helpCmd returns information on the usage of the command
func (cmd Command) helpCmd() string {
	desc := fmt.Sprintf(
		"\nUsage: %s\n\n%s %s\n",
		cmd.Label,
		cmd.Label,
		describeArguments(cmd.Arguments),
	)
	for _, arg := range cmd.Arguments {
		desc += fmt.Sprintf("%s %s", cmd.Label, arg.helpArg())
	}
	return desc
}

// helpArg returns information on the usage of the argument using the command
func (cmd Command) helpArg(arg Argument) string {
	return fmt.Sprintf(
		"\nUsage: %s %s\n\n%s %s",
		cmd.Label,
		arg.friendlyName(),
		cmd.Label,
		arg.helpArg(),
	)
}

// help returns information on the usage of the argument
func (arg Argument) helpArg() string {
	return fmt.Sprintf(
		"%s %s\n",
		arg.friendlyName(),
		describeOptions(arg.Options),
	)
}

// describeArguments describes the arguments using command syntax convention
func describeArguments(arguments []Argument) string {

	// In command syntax convention, arguments are displayed
	// in a list separated by (|)
	var optional bool
	labels := make([]string, 0)
	for _, argument := range arguments {
		if argument.Label == "" {
			optional = true
			continue
		}
		labels = append(labels, argument.Label)
	}
	desc := strings.Join(labels, "|")

	// If there was a argument with an empty label
	// it means that the command can be inkoved
	// without any options, hence the arguments are
	// optional and need to be displayed inside
	// square brackets to reflect this.
	if optional {
		desc = fmt.Sprintf("[%s]", desc)
	}
	desc += "\n"

	// Format for the padding
	var longestLabelLength int
	var label string
	for _, argument := range arguments {
		label = argument.friendlyName()
		if longestLabelLength < len(label) {
			longestLabelLength = len(label)
		}
	}
	paddingStr := fmt.Sprintf("%%-%ds", longestLabelLength)

	// List each argument with its help message and correct padding
	for _, argument := range arguments {
		label = argument.friendlyName()
		desc += fmt.Sprintf("\t"+paddingStr+" %s\n", label, argument.HelpMsg)
	}
	return desc
}

// friendlyName returns the friendly name for the argument.
// Namely, it returns the argument label unless it is empty,
// in that case it returns "(no arguments)".
func (arg Argument) friendlyName() string {
	if arg.Label == "" {
		return "(no arguments)"
	}
	return arg.Label
}

// describeOptions describes the options using command syntax convention
func describeOptions(options []Option) string {

	// In command syntax convention, options are split into 4 distinct categories.
	// They are split on whether they have a short name or not and split on
	// whether they are required. The required short will be displayed first,
	// followed by the required long. The optional short and optional long follow
	// after that, with both of these sections being incased in square braces ([])
	// to show that they are optional. The optional short are also combined, such
	// that there is one dash (-) followed by all the optional short name
	// characters.
	var reqShort, reqLong, optShort, optLong string
	for _, option := range options {
		if option.Variable != nil && option.Variable.Required {
			if option.Short != "" {
				reqShort += fmt.Sprintf("%s %s ", option.Short, option.Variable.Label)
			} else {
				reqLong += fmt.Sprintf("%s=%s ", option.Long, option.Variable.Label)
			}
		} else {
			if option.Short != "" {
				if optShort == "" {
					optShort = option.Short
				} else {
					optShort += option.Short[1:]
				}
			} else {
				optLong += fmt.Sprintf("%s ", option.Long)
			}
		}
	}

	var desc string
	reqShort = strings.TrimRight(reqShort, " ")
	if reqShort != "" {
		desc += reqShort + " "
	}
	reqLong = strings.TrimRight(reqLong, " ")
	if reqLong != "" {
		desc += reqLong + " "
	}
	optShort = strings.TrimRight(optShort, " ")
	if optShort != "" {
		desc += optShort + " "
	}
	optLong = strings.TrimRight(optLong, " ")
	if optLong != "" {
		desc += optLong
	}
	desc = strings.TrimRight(desc, " ")
	desc += "\n"

	// Format for the padding
	var longestLongLength int
	for _, option := range options {
		if longestLongLength < len(option.Long) {
			longestLongLength = len(option.Long)
		}
	}

	// List each option with its help message and correct padding
	paddingStr := fmt.Sprintf("%%-%ds", longestLongLength)
	for _, option := range options {
		desc += fmt.Sprintf("\t%s "+paddingStr+" %s\n", option.Short, option.Long, option.HelpMsg)
	}
	return desc
}
