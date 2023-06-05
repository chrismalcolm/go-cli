package cli

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

// Config represents the all the config options required to run the CLI.
type Config struct {

	// The output to the CLI to prompt input from the user.
	Prompt string `yaml:"prompt"`

	// The commands that are configured.
	Commands []Command `yaml:"commands"`

	// The function performed when the CLI is intialised.
	// The output from this function will appear before
	// any other output in the CLI.
	InitFunc string `yaml:"initFunc"`
	init     func(Flags) []byte

	// The function performed when the CLI is terminated.
	// This function's output will be the last output to
	// appear in the CLI before it closes.
	ExitFunc string `yaml:"exitFunc"`
	exit     func(Flags) []byte

	// The function performed when the user requests help.
	// This is a built in function that is automatically
	// created when the config is initialised.
	help func(Flags) []byte

	// The CLI command used to trigger an exit.
	ExitCmd string `yaml:"exitCmd"`

	// The CLI command used to print a help message.
	HelpCmd string `yaml:"helpCmd"`
}

// LoadConfig extracts the config from the given yaml
// file and unmarshals it into a Config.
// Any errors reading the file or unmarshaling the file
// will be returned.
func LoadConfig(filename string) (config *Config, err error) {

	// Attempt to read the file
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// Attempt to unmarshal the yaml file into config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return config, err
	}

	// Validation check on the exit command
	if config.ExitCmd == "" {
		return config, fmt.Errorf("missing/empty exit command \"exitCmd\"")
	}

	// Validation check on the help command
	if config.HelpCmd == "" {
		return config, fmt.Errorf("missing/empty help command \"helpCmd\"")
	}

	// Config needs to have at least one command
	if len(config.Commands) == 0 {
		return config, fmt.Errorf("missing/empty commands \"commands\"")
	}

	// Validation check on the commands
	for _, command := range config.Commands {
		if command.Label == config.ExitCmd {
			return config, fmt.Errorf("command cannot share same label as exit command \"%s\"", config.ExitCmd)
		}
		if command.Label == config.HelpCmd {
			return config, fmt.Errorf("command cannot share same label as help command \"%s\"", config.HelpCmd)
		}
		if err := command.validate(); err != nil {
			return config, err
		}
	}

	// Generate placeholder and help commands
	config.init = func(Flags) []byte { return []byte("") }
	config.exit = func(Flags) []byte { return []byte("") }
	for i, command := range config.Commands {
		for j, argument := range command.Arguments {
			config.Commands[i].Arguments[j].executable = command.createExecutable(argument)
			config.Commands[i].Arguments[j].help = command.createArgHelp(argument)
		}
		config.Commands[i].help = command.createHelp()
	}
	config.help = config.createHelp()

	return config, nil
}

// WithProgram maps the execFuncs defined in the config
// to methods with the same name in program.
// If the program does not have a method with the same name,
// or the method is not of the executable type
// (func(Flags) []byte), then an error will be returned.
func (config *Config) withProgram(program interface{}) (err error) {

	// Apply the init method.
	config.init, err = getExecutable(program, config.InitFunc)
	if err != nil {
		return err
	}

	// Apply the exit method.
	config.exit, err = getExecutable(program, config.ExitFunc)
	if err != nil {
		return err
	}

	// Apply the argument methods.
	for i, command := range config.Commands {
		for j, argument := range command.Arguments {
			if argument.ExecFunc != "" {
				config.Commands[i].Arguments[j].executable, err = getExecutable(program, argument.ExecFunc)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getExecutable attempts to return the method from the program from the given funcName.
// If the method doesn't exist or is not of the correct type (func(Flags) []byte), an
// error will be returned.
func getExecutable(program interface{}, funcName string) (action func(Flags) []byte, err error) {

	// Create and run a panic-safe func. This func will attempt to do the following:
	// - Use the reflect package to get the action method called funcName.
	// - Convert that reflecr.Value into an interface
	// - Type cast the interface into the correct type for an action method (func(Flags) []byte)
	// If any of these stages fail, a panic will be raised. This is caught by the recover()
	// in the defer statment, so that if a panic occurres, we will not terminate.
	// If the panic is caught, this function will return <nil> as the return value.
	action = func() func(Flags) []byte {
		defer func() {
			recover()
		}()
		return reflect.ValueOf(program).MethodByName(funcName).Interface().(func(Flags) []byte)
	}()

	// Raise an error if enable to find the method funcName
	if action == nil {
		return action, fmt.Errorf("unable to find method \"%s\" for type \"%s\"", funcName, reflect.TypeOf(program))
	}

	return action, nil
}
