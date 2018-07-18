package main

import (
	"fmt"
)

var (
	VERSION      string
	BUILD_DATE   string
	GIT_REVISION string
	GO_VERSION   string
)

type versionCmd struct {
	showAll bool
}

func (c *versionCmd) Usage() string {
	return `Usage: auth-server version

Show auth-server version.`
}

func (c *versionCmd) Flags() {
}

func (c *versionCmd) Run(args []string) error {
	fmt.Printf("Version: %s\n", VERSION)
	fmt.Printf("Git rev: %s\n", GIT_REVISION)
	fmt.Printf("Go version: %s\n", GO_VERSION)
	fmt.Printf("Built: %s\n", BUILD_DATE)

	return nil
}
