package cli

import (
	"fmt"
	"strings"
)

// validate performs a validation check on a Command
func (cmd Command) validate() error {

	// Label must be a non-empty string
	if cmd.Label == "" {
		return fmt.Errorf("empty command label detected")
	}

	// Label must not contain any whitespace characters
	if strings.ContainsAny(cmd.Label, " \n\r\t") {
		return fmt.Errorf("invalid command label \"%s\", invalid whitespace characters detected", cmd.Label)
	}

	// There must be at least one argument
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("command \"%s\" requires at least one argument", cmd.Label)
	}

	// Arguments must all be valid and do not repeat
	labels := make(map[string]bool)
	for _, arg := range cmd.Arguments {

		// Arguments must all be valid
		if err := arg.validate(); err != nil {
			return fmt.Errorf("command \"%s\", %s", cmd.Label, err)
		}

		// Argument labels must not repeat
		if _, alreadyExists := labels[arg.Label]; alreadyExists {
			return fmt.Errorf("command \"%s\", multiple occurrences of the argument label \"%s\"", cmd.Label, arg.Label)
		}
		labels[arg.Label] = true
	}

	return nil
}

// validate performs a validation check on an Argument
func (arg Argument) validate() error {

	// Label must not contain any special whitespace characters
	if strings.ContainsAny(arg.Label, "\n\r\t") {
		return fmt.Errorf("invalid argument label \"%s\", invalid special whitespace characters detected", arg.Label)
	}

	// Label must not start or end with spaces
	if strings.TrimLeft(arg.Label, " ") != arg.Label {
		return fmt.Errorf("invalid argument label \"%s\", spaces detected at start", arg.Label)
	}
	if strings.TrimRight(arg.Label, " ") != arg.Label {
		return fmt.Errorf("invalid argument label \"%s\", spaces detected at end", arg.Label)
	}

	// Options must all be valid and must not repeat labels, shorts or longs
	labels := make(map[string]bool)
	shorts := make(map[string]bool)
	longs := make(map[string]bool)
	for _, opt := range arg.Options {

		// Options must all be valid
		if err := opt.validate(); err != nil {
			return fmt.Errorf("argument \"%s\", %s", arg.Label, err)
		}

		// Option labels must not repeat
		if _, alreadyExists := labels[opt.Label]; alreadyExists {
			return fmt.Errorf("argument \"%s\", multiple occurrences of the option label \"%s\"", arg.Label, opt.Label)
		}
		labels[opt.Label] = true

		// Option shorts must not repeat
		if opt.Short != "" {
			if _, alreadyExists := shorts[opt.Short]; alreadyExists {
				return fmt.Errorf("argument \"%s\", multiple occurrences of the option short \"%s\"", arg.Label, opt.Short)
			}
			shorts[opt.Short] = true
		}

		// Option longs must not repeat
		if opt.Long != "" {
			if _, alreadyExists := longs[opt.Long]; alreadyExists {
				return fmt.Errorf("argument \"%s\", multiple occurrences of the option long \"%s\"", arg.Label, opt.Long)
			}
			longs[opt.Long] = true
		}
	}

	return nil
}

// validate performs a validation check on an Option.
func (opt Option) validate() error {

	// Label must be a non-empty string
	if opt.Label == "" {
		return fmt.Errorf("empty option label detected")
	}

	// Label must not contain any special whitespace characters
	if strings.ContainsAny(opt.Label, "\n\r\t") {
		return fmt.Errorf("invalid option label \"%s\", invalid whitespace characters detected", opt.Label)
	}

	// At least one of Short or Long must be provided
	if opt.Short == "" && opt.Long == "" {
		return fmt.Errorf("at least one of option short or option long must be provided")
	}

	// If applicable, Short must be single dash (–) followed by a signle
	// character, with no invalid characters.
	if opt.Short != "" {

		// Short must start with a single dash (-)
		if !strings.HasPrefix(opt.Short, "-") {
			return fmt.Errorf("invalid option short \"%s\", must start with a single dash (-)", opt.Short)
		}

		// Short must only be two characters long
		if len(opt.Short) != 2 {
			return fmt.Errorf("invalid option short \"%s\", must be a single dash (-) followed by a single character", opt.Short)
		}

		// Short must not contain any whitespace characters
		if strings.ContainsAny(opt.Short[1:], " \n\r\t") {
			return fmt.Errorf("invalid option short \"%s\", whitespace characters detected", opt.Short)
		}

		// Short must not contain any other invalid characters
		if strings.ContainsAny(opt.Short[1:], "[]{}()-=") {
			return fmt.Errorf("invalid option short \"%s\", invalid characters detected", opt.Short)
		}
	}

	// If applicable, Long must be double dash (–-) followed by a
	// descriptive name, with no invalid characters.
	if opt.Long != "" {

		// Long must start with a double dash (--)
		if !strings.HasPrefix(opt.Long, "--") {
			return fmt.Errorf("invalid option long \"%s\", must start with a double dash (--)", opt.Long)
		}

		// Long must longer than two characters
		if len(opt.Long) == 2 {
			return fmt.Errorf("invalid option long \"%s\", must be longer than two characters", opt.Long)
		}

		// Long must not contain any special whitespace characters
		if strings.ContainsAny(opt.Long, "\n\r\t") {
			return fmt.Errorf("invalid option long \"%s\", special whitespace characters detected", opt.Long)
		}

		// Long must not contain any other invalid characters
		if strings.ContainsAny(opt.Long, "[]{}()=") {
			return fmt.Errorf("invalid option long \"%s\", invalid characters detected", opt.Long)
		}
	}

	// The variable must be valid
	if opt.Variable != nil {
		if err := opt.Variable.validate(); err != nil {
			return fmt.Errorf("option \"%s\", %s", opt.Label, err)
		}
	}

	return nil
}

// validate performs a validation check on an Argument
func (va Variable) validate() error {

	// Label must be a non-empty string
	if va.Label == "" {
		return fmt.Errorf("empty variable label detected")
	}

	// Label must not contain any special whitespace characters
	if strings.ContainsAny(va.Label, "\n\r\t") {
		return fmt.Errorf("invalid variable label \"%s\", invalid whitespace characters detected", va.Label)
	}

	return nil
}
