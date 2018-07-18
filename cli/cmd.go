package cli

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

type Command interface {
	Usage() string
	Flags()
	Run(args []string) error
}

type CommandsMap map[string]Command

var ErrArgs = fmt.Errorf("wrong number of subcommand arguments")
var ErrUsage = fmt.Errorf("show usage")

func Exec(commands CommandsMap) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintln(os.Stderr, x)
		}
	}()
	var err error
	// check number of arguments
	if len(os.Args) < 2 {
		commands["help"].Run(nil)
		os.Exit(1)
	}

	// check action
	name := os.Args[1]
	command, ok := commands[name]
	if !ok {
		commands["help"].Run(nil)
		fmt.Fprintf(os.Stderr, "\n"+"error: unknown command: %s"+"\n", name)
		os.Exit(1)
	}
	command.Flags()

	flag.Usage = func() {
		fmt.Print(command.Usage())
		fmt.Printf("\n\n%s\n", "Options:")

		//flag.Set(os.Stdout)
		flag.PrintDefaults()
		os.Exit(0)
	}
	os.Args = os.Args[1:]
	flag.Parse()

	// run action
	err = command.Run(flag.Args())
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		if err == ErrArgs || err == ErrUsage {
			fmt.Println("")
			flag.Usage()
		}
	}
}
