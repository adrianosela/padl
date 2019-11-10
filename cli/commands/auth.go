package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/adrianosela/padl/api/client"

	cli "gopkg.in/urfave/cli.v1"
)

// AuthCmds - manage accounts and auth flows
var AuthCmds = cli.Command{
	Name:    "auth",
	Aliases: []string{"a"},
	Usage:   "manage accounts",
	Subcommands: []cli.Command{
		{
			Name:  "register",
			Usage: "register a new account on a padl server",
			Flags: []cli.Flag{
				asMandatory(emailFlag),
				asMandatory(pathFlag),
			},
			Before: registerAccountValidator,
			Action: registerAccountHandler,
		},
	},
}

func registerAccountValidator(ctx *cli.Context) error {
	return assertSet(ctx, emailFlag, pathFlag)
}

func registerAccountHandler(ctx *cli.Context) error {
	c, err := client.NewClient("http://localhost")
	if err != nil {
		return err
	}

	email := ctx.String(name(emailFlag))
	pubPath := ctx.String(name(pathFlag))

	pubKey, err := ioutil.ReadFile(pubPath)
	if err != nil {
		return fmt.Errorf("could not read PGP public key file %s: %s", pubPath, err)
	}

	if err := c.Register(email, string(pubKey)); err != nil {
		return fmt.Errorf("%s", err)
	}

	fmt.Printf("registered user %s successfully!\n", email)
	return nil
}
