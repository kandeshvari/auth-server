package user

import (
	"crypto"
	"fmt"
	"os"

	"github.com/kandeshvari/auth-server/auth"
	"github.com/kandeshvari/auth-server/cli"
	"github.com/kandeshvari/auth-server/config"
	"github.com/kandeshvari/auth-server/logger"
	flag "github.com/ogier/pflag"
	//"github.com/kandeshvari/auth-server/table_utils"
	"github.com/olekukonko/tablewriter"
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
	return fmt.Sprintf(`Usage: auth-server user <sub-action> 

Run auth server user management cli.

sub-actions:

ls
	list all users

add <username> <password>
	add new user

chpass <username> <password>
	change user password

rm <username>
	remove user
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

	// choose main loggers
	ctx.log, err = logger.GetLoggerErr("auth")
	if err != nil {
		ctx.log = logger.GetLogger("default")
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

	switch args[0] {
	case "add":
		return c.addAction(args)
	case "chpass":
		return c.chpassAction(args)
	case "ls":
		return c.lsAction(args)
	case "rm":
		return c.rmAction(args)
	default:
		return cli.ErrArgs
	}
	return nil
}

/* Actions */

// Add new user
func (c *Cmd) addAction(args []string) error {
	if len(args) != 3 {
		return cli.ErrArgs
	}

	hash := crypto.SHA512_224.New()
	hash.Write([]byte(args[2]))
	passwordHash := fmt.Sprintf("%x", hash.Sum(nil))

	err := ctx.auth.CreateUser(args[1], passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cmd) chpassAction(args []string) error {
	return nil
}

func (c *Cmd) lsAction(args []string) error {
	users, err := ctx.auth.ListUsers()
	if err != nil {
		return err
	}

	data := make([][]string, 0, 1024) // rows with images

	for _, user := range users {
		data = append(data, []string{fmt.Sprintf("%d", user.Id), user.Login})
	}

	table := tablewriter.NewWriter(os.Stdout)
	//table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetHeader([]string{"ID", "LOGIN"})
	//sort.Sort(table_utils.ByName(data))
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (c *Cmd) rmAction(args []string) error {
	if len(args) != 2 {
		return cli.ErrArgs
	}

	user, err := ctx.auth.GetUserByLogin(args[1])
	if err != nil {
		return err
	}

	err = ctx.auth.DeleteUser(user.Id)
	if err != nil {
		return err
	}

	return nil
}
