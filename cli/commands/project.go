package commands

import (
	cli "gopkg.in/urfave/cli.v1"
)

// ProjectCmds - manage projects
var ProjectCmds = cli.Command{
	Name:    "project",
	Aliases: []string{"p"},
	Usage:   "manage projects",
	Subcommands: []cli.Command{
		{
			Name:  "list",
			Usage: "list all projects user is a member of",
			Flags: []cli.Flag{
				jsonFlag,
			},
			Action: projectListHandler,
		},
	},
}

func projectListHandler(ctx *cli.Context) error {
	// TODO
	return nil
}
