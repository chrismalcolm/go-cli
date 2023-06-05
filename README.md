# CLI
Go CLI tool, used to create a CLI only with a config file and a struct with the methods.

## How it works
1) Define config in a config file
2) Create a struct which has methods as defined in the config file
3) Apply it all together

## Example
1) Sample yaml file
```yaml
---
initFunc: Setup
exitFunc: Leave
exitCmd: leave
helpCmd: help
prompt: ">>> "
commands:
  - label: show
    arguments:
      - label: ""
        help: "Show will show the stuff you want"
        execFunc: MyFunc
        options:
          - label: readOnly
            short: -r
            long: --read-only
            help: "Short desc"
            variable:
              label: var4
              required: false
              defaul: "das"
          - label: desc
            short: -d
            long: --description
            help: "Short desc"
            variable:
              label: var4
              required: true
              defaul: "das"
      - label: daily-tasks
        help: "Show will show the other stuff you want"
        execFunc: MyFunc2
        options:
          - label: readOnly
            short: -r
            long: --read-only
            help: "Short desc"
            variable:
              label: var4
              required: true
              defaul: "das"
```

2) Program struct
```go
type Program struct {
}

func (program *Program) Setup(cli.Flags) []byte {
	return []byte("Starting\n")
}

func (program *Program) Leave(cli.Flags) []byte {
	return []byte("Leaving\n")
}

func (program *Program) MyFunc(flags cli.Flags) []byte {
	return []byte(fmt.Sprintf("MyFunc called flags %v\n", flags))
}

func (program *Program) MyFunc2(flags cli.Flags) []byte {
	return []byte(fmt.Sprintf("MyFunc2 called flags %v\n", flags))
}
```

3) Run in main.go
```go
func main() {

	config, err := cli.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	app, err := cli.New(config).Using(&Program{})
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
```

## Ideas for future
* Use of arrow keys to scroll up and down previous inputs/outputs
* Allow split config file into many config files, with hierarchical structure