package main

import (
	flag "github.com/ogier/pflag"

	"fmt"
	"os"
	"sort"
	"strings"
)

type helpCmd struct {
	showAll bool
}

func (c *helpCmd) Usage() string {
	return `Usage: auth-server help [--all]

Help page for the auth server.`
}

func (c *helpCmd) Flags() {
	flag.BoolVar(&c.showAll, "all", false, "Show all commands (not just interesting ones)")
}

func (c *helpCmd) Run(args []string) error {
	if len(args) > 0 {
		for _, name := range args {
			command, ok := commands[name]
			if !ok {
				fmt.Fprintf(os.Stderr, "error: unknown command: %s"+"\n", name)
			} else {
				fmt.Fprintf(os.Stdout, command.Usage()+"\n")
			}
		}
		return nil
	}

	fmt.Println("Usage: auth-server <command> [options]")
	fmt.Println()
	fmt.Println(`This is the auth server.

All of auth server's features can be driven through the various commands below.
For help with any of those, simply call them with --help.`)
	fmt.Println()
	fmt.Println("Commands:")
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if name == "help" {
			continue
		}

		command := commands[name]
		fmt.Printf("  %-16s %s\n", name, summaryLine(command.Usage()))
	}

	fmt.Println()

	return nil
}

func summaryLine(usage string) string {
	for _, line := range strings.Split(usage, "\n") {
		if strings.HasPrefix(line, "Usage:") {
			continue
		}

		if len(line) == 0 {
			continue
		}

		return strings.TrimSuffix(line, ".")
	}

	return "Missing summary."
}
