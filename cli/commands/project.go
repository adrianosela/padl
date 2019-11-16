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
	path := ctx.String(name(pathFlag))
	pname := ctx.String(name(nameFlag))
	descr := ctx.String(name(descriptionFlag))
	json := ctx.Bool(name(jsonFlag))

	if path == "" {
		if json {
			path = "./padlfile.json"
		} else {
			path = "./padlfile.yaml"
		}
	}

	if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".yaml") {
		return fmt.Errorf("invalid file extension, must be one of { \".yaml\", \".json\" }")
	}

	pf, err := c.CreateProject(pname, descr)
	if err != nil {
		return fmt.Errorf("error creating project: %s", err)
	}

	err = pf.Write(path)
	if err != nil {
		return fmt.Errorf("unable to write padl file: %s", err)
	}
	return nil
}
