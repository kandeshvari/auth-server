package server

import (
	"fmt"

	"github.com/kandeshvari/auth-server/auth"
	"github.com/kandeshvari/auth-server/cli"
	"github.com/kandeshvari/auth-server/config"
	"github.com/kandeshvari/auth-server/logger"

	flag "github.com/ogier/pflag"
	"github.com/pkg/errors"
)

type Cmd struct {
	config string
}

type Ctx struct {
	conf *config.Config
	log  *logger.Logger
	auth *auth.Auth
}

var ctx Ctx

func (c *Cmd) Usage() string {
	return fmt.Sprintf(`Usage: auth-server server <sub-action> 

Run auth server.

sub-actions:

start           start auth server
`)
}

func (c *Cmd) Flags() {
	flag.StringVar(&c.config, "config", "/etc/auth-server/auth-server.yaml", "Path to config file")
}

func (c *Cmd) Run(args []string) error {
	var err error
	ctx.conf, err = config.LoadConfig(c.config)
	if err != nil {
		return err
	}

	err = logger.SetupLoggers(ctx.conf.Logger)
	if err != nil {
		return err
	}

	//choose main loggers
	ctx.log, err = logger.GetLoggerErr("auth")
	if err != nil {
		ctx.log = logger.GetLogger("default")
	}

	// validate mandatory config variables
	if ctx.conf.Auth.SecretKey == "" {
		return errors.New("config: non empty auth.secret_key must be provided")
	}

	// create auth connection
	ctx.auth, err = auth.NewAuth(ctx.conf.Database.Driver, ctx.conf.Database.DbConnect, ctx.conf.Auth.LoginDelay)
	if err != nil {
		return err
	}

	// check and run sub-actions
	if len(args) == 0 {
		return cli.ErrArgs
	}
	if args[0] == "start" {
		ctx.log.I("Auth server starting on %s (tls_dir=%s)", ctx.conf.Server.Address, ctx.conf.Server.TLSDir)
		err = RunAuthServer()
		if err != nil {
			ctx.log.C(err.Error())
			return err
		}

	} else {
		return cli.ErrArgs
	}
	return nil
}
