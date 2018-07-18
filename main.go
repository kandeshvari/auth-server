package main

import (
	"sync"

	"github.com/kandeshvari/auth-server/cli"
	"github.com/kandeshvari/auth-server/cli/server"
	"github.com/kandeshvari/auth-server/cli/user"
	"github.com/kandeshvari/auth-server/config"

	_ "github.com/go-sql-driver/mysql"
)

type CXT struct {
	wg   sync.WaitGroup
	conf *config.Config
}

var commands = cli.CommandsMap{
	"help":    &helpCmd{},
	"version": &versionCmd{},
	"server":  &server.Cmd{},
	"user":    &user.Cmd{},
}

func main() {
	cli.Exec(commands)
}
