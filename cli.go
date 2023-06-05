package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
)

var whitespaceCharacters = " \n\r\t"

// App is the CLI application
type App struct {
	config *Config
	writer *bufio.Writer
	reader *bufio.Reader
	sigint chan os.Signal
	active bool
}

// New creates a new App from the given config
func New(config *Config) (app *App) {
	return &App{
		config: config,
		writer: bufio.NewWriter(os.Stdout),
		reader: bufio.NewReader(os.Stdin),
		sigint: make(chan os.Signal, 1),
		active: true,
	}
}

// Using gets the App to use the methods from program
func (app *App) Using(program interface{}) (*App, error) {
	if err := app.config.withProgram(program); err != nil {
		return app, err
	}
	return app, nil
}

// Run runs the CLI
func (app *App) Run() {

	// Write CLI initial input
	initOutput := app.config.init(Flags{})
	if err := app.write([]byte(initOutput)); err != nil {
		log.Fatal(err)
	}

	// Run the CLI until terminaled by ctl-C or user exit
	go func() {

		for app.active {

			// Write CLI prompt
			if err := app.write([]byte(app.config.Prompt)); err != nil {
				log.Fatal(err)
			}

			// Get input from CLI
			input, err := app.read()
			if err != nil {
				log.Fatal(err)
			}

			// Get output from cli
			output := app.getOutput(input)

			// Write output
			if err := app.write(output); err != nil {
				log.Fatal(err)
			}
		}

		app.sigint <- os.Kill
	}()

	// Interrupt on ctl-C
	signal.Notify(app.sigint, os.Interrupt)
	<-app.sigint
}

// write writes bytes to the CLI
func (app *App) write(b []byte) error {

	// Attempt to write bytes to writer
	if _, err := app.writer.Write(b); err != nil {
		return err
	}

	// Flush write buffer to display to the screen
	if err := app.writer.Flush(); err != nil {
		return err
	}

	return nil
}

// read reads input from the CLI
func (app *App) read() (str string, err error) {

	// Attempt tog et input from user
	str, err = app.reader.ReadString('\n')
	if err != nil {
		return str, err
	}

	return str, nil
}

// getOutput extracts the command from the input and runs the correct executable
func (app *App) getOutput(input string) []byte {

	// Trim any whitespace from the input
	input = strings.Trim(input, whitespaceCharacters)
	if input == "" {
		return []byte{}
	}

	// Exit the CLI if the ExitCmd is the input
	if input == app.config.ExitCmd {
		app.prepareExit()
		return app.config.exit(Flags{})
	}

	// If input ends with help coomand, remove help command from input
	// and return the help output instead.
	if strings.HasSuffix(input, app.config.HelpCmd) {
		input = strings.TrimRight(input[:len(input)-len(app.config.HelpCmd)], whitespaceCharacters)
		return app.getHelpOutput(input)
	}

	// Extract the command and reamining input after removing the input
	command, remainingInput, err := app.extractCommand(input)
	if err != nil {
		return []byte(fmt.Sprintf("%v\n", err))
	}

	// Get the argument and flags
	argument, optionsInput, err := app.extractArgument(remainingInput, command)
	if err != nil {
		return []byte(fmt.Sprintf("%v\n", err))
	}

	// Attempt to extraxt the flags from the options input
	flags, err := app.extractFlags(optionsInput, argument)
	if err != nil {
		return []byte(fmt.Sprintf("%v\n", err))
	}

	// Return the output from the executable
	return argument.executable(flags)
}

// getHelpOutput extracts the help command output.
// The input here should be the original input but
// with the help command removed.
func (app *App) getHelpOutput(input string) []byte {

	// If there is no input left, original command must've been
	// just the help command. Hence, run the global help command.
	if input == "" {
		return app.config.help(Flags{})
	}

	// Extract the command and reamining input after removing the input
	command, remainingInput, err := app.extractCommand(input)
	if err != nil {
		return []byte(fmt.Sprintf("%v\n", err))
	}

	// If there is no remaining input, the original command must've
	// been a single command followed by the help command.
	if remainingInput == "" {
		return command.help(Flags{})
	}

	// Get the argument and flags
	argument, _, err := app.extractArgument(remainingInput, command)
	if err != nil {
		return []byte(fmt.Sprintf("%v\n", err))
	}

	// Return the argument version of the help.
	return argument.help(Flags{})
}

// extractCommand extracts the command string from the input.
// It also returns the remaining input, which is the original
// output with the preceeding command string removed.
func (app *App) extractCommand(input string) (command Command, remainingInput string, err error) {

	// Extract the command label
	var commandLabel string
	index := strings.Index(input, " ")
	if index == -1 {
		commandLabel = input
	} else {
		commandLabel = input[:index]
		remainingInput = strings.TrimLeft(input[index:], whitespaceCharacters)
	}

	// Search for the command from the config
	for _, cmd := range app.config.Commands {
		if cmd.Label == commandLabel {
			return cmd, remainingInput, nil
		}
	}

	// Return an error if unable to find the command in the config
	return command, remainingInput, fmt.Errorf("unable to find command \"%s\"", commandLabel)
}

// extractArgument extracts the argument.
// It also returns the options input, which is
// the input with the argument label removed.
func (app *App) extractArgument(remainingInput string, command Command) (argument Argument, optionsInput string, err error) {

	// Attempt to find an argument that is in the remaining input
	var foundArg bool
	for _, arg := range command.Arguments {

		// If argument label is empty, this represents a command with no arguments.
		// However, this will be overriden if we continue iterating and find a match
		// with an argument label that isn't empty.
		if arg.Label == "" {
			optionsInput = remainingInput
			argument = arg
			foundArg = true
			continue
		}

		// If the argument label is not in the remaining input, continue
		index := strings.Index(remainingInput, arg.Label)
		if index == -1 {
			continue
		}

		// Validation for the input after the argument
		after := remainingInput[index+len(arg.Label):]
		if after != "" && !strings.HasPrefix(after, " ") {
			continue
		}

		// Set the options input as the remaining input with the argument label removed
		// and break out of the loop
		optionsInput = remainingInput[:index] + " " + after
		argument = arg
		foundArg = true
		break
	}

	// If no matching argument has been found, return an error
	if !foundArg {
		return argument, optionsInput, fmt.Errorf("invalid use of the \"%s\" command, no valid argument provided", command.Label)
	}

	return argument, optionsInput, nil
}

// extractFlags extracts the flags from the options input
func (app *App) extractFlags(optionsInput string, argument Argument) (flags Flags, err error) {

	// Set the default flag metadata for the flags
	metadata := make(map[string]flagMetadata, 0)
	for _, option := range argument.Options {
		metadata[option.Label] = flagMetadata{
			isset:    false,
			hasVar:   false,
			variable: "",
		}
	}

	// Loop though the option flags remaining in the input
	// For each flag found, re-configure the flag metadata.
	var expectingValue bool
	optionsStrings := strings.Split(optionsInput, " ")
	for i, s := range optionsStrings {

		// Ignore empty strings
		if s == "" {
			continue
		}

		// Check if using a option short or long name
		if strings.HasPrefix(s, "-") {
			var variable string
			var shortVersion bool

			// Loop though all options
			for _, option := range argument.Options {
				if option.Short != s && option.Long != s {
					continue
				}
				shortVersion = option.Short == s

				// If the option requires no variable, re-configure the flag metadata
				// for this option as isset = true
				if option.Variable == nil {
					metadata[option.Label] = flagMetadata{
						isset:    true,
						hasVar:   false,
						variable: "",
					}
					break
				}

				// The option requires a variable, for now set the variable as default value
				variable = option.Variable.Default

				// For short version, syntax will be -<char> <variable> e.g. (-a read).
				// We will use expectingValue = true to ignore the non-flag text in the
				// next loop iteration.
				if shortVersion {
					if i+1 < len(optionsStrings) {
						variable = optionsStrings[i+1]
					} else if option.Variable.Required {
						return flags, fmt.Errorf("missing variable \"%s\" for option \"%s\"", option.Variable.Label, option.Label)
					}
					metadata[option.Label] = flagMetadata{
						isset:    true,
						hasVar:   true,
						variable: variable,
					}
					expectingValue = true
					break
				}

				// For long version, syntax will be --<chars>=<variable> e.g. (--append=true).
				if index := strings.Index(s, "="); index != -1 {
					variable = s[index:]
				} else if option.Variable.Required {
					return flags, fmt.Errorf("required option \"%s\" missing required variable \"%s\"", option.Label, option.Variable.Label)
				}
				metadata[option.Label] = flagMetadata{
					isset:    true,
					hasVar:   true,
					variable: variable,
				}
				break
			}
		} else if !expectingValue {
			return flags, fmt.Errorf("invalid text \"%s\" detected", s)
		}
	}

	return Flags{mapping: metadata}, nil
}

// prepareExit prepare the CLI to exit after the next output has been sent
func (app *App) prepareExit() {
	app.active = false
}
