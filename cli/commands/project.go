package commands

import (
	"fmt"
	"strings"

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
		{
			Name: "create",
			Flags: []cli.Flag{
				asMandatory(nameFlag),
				asMandatory(descriptionFlag),
				pathFlag,
				// Defaulting to yml
				jsonFlag,
			},
			Before: createProjectValidator,
			Action: createProjectHandler,
		},
	},
}

func projectListHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func createProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, descriptionFlag)
}

func createProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}
	p := ctx.String(name(pathFlag))
	pname := ctx.String(name(nameFlag))
	descr := ctx.String(name(descriptionFlag))
	json := ctx.Bool(name(jsonFlag))

	if p == "" {
		if json {
			p = "./padlfile.json"
		} else {
			p = "./padlfile.yaml"
		}
	}

	if !strings.Contains(p, ".json") && !strings.Contains(p, ".yaml") {
		return fmt.Errorf("provide file name and type")
	}

	pf, err := c.CreateProject(pname, descr)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err)
	}

	err = pf.Write(p)
	if err != nil {
		return fmt.Errorf("unable to write padl file: %s", err)
	}
	return nil
}
