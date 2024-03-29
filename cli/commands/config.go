package commands

import (
	"fmt"
	"os"

	"github.com/adrianosela/padl/cli/config"
	"github.com/olekukonko/tablewriter"
	cli "gopkg.in/urfave/cli.v1"
)

// ConfigCmds is the CLI command object for the config operation
var ConfigCmds = cli.Command{
	Name:    "config",
	Aliases: []string{"c"},
	Usage:   "Manage configuration and CLI settings",
	Subcommands: []cli.Command{
		{
			Name:  "set",
			Usage: "create configuration file with given options",
			Flags: []cli.Flag{
				asMandatory(urlFlag),
				ConfigFlag,
			},
			Before: configSetValidator,
			Action: configSetHandler,
		},
		{
			Name:  "show",
			Usage: "show contents of set configuration file",
			Flags: []cli.Flag{
				ConfigFlag,
			},
			Action: configShowHandler,
		},
	},
}

func configSetValidator(ctx *cli.Context) error {
	return assertSet(ctx, urlFlag)
}

func configSetHandler(ctx *cli.Context) error {
	if err := config.SetConfig(&config.Config{
		HostURL: ctx.String(name(urlFlag)),
	}, ctx.String(name(pathFlag))); err != nil {
		return fmt.Errorf("could not set configuration: %s", err)
	}
	return nil
}

func configShowHandler(ctx *cli.Context) error {
	path := ctx.String(name(pathFlag))
	c, err := config.GetConfig(path)
	if err != nil {
		return fmt.Errorf("could not retrive configuration from %s: %s", path, err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"HOST_URL", c.HostURL})
	if c.User != "" {

		table.Append([]string{"USER", c.User})
	}
	if c.Token != "" {
		table.Append([]string{"AUTH_TOKEN", "SECRET-TOKEN-PRESENT"})
	}
	table.Render()
	return nil
}
