package commands

import (
	"encoding/json"
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
			Name:  "create",
			Usage: "create a new project",
			Flags: []cli.Flag{
				asMandatory(nameFlag),
				asMandatory(descriptionFlag),
				pathFlag,
				// Defaulting to yml
				withDefault(fmtFlag, ".yaml"),
			},
			Before: createProjectValidator,
			Action: createProjectHandler,
		},
		{
			Name:  "get",
			Usage: "finds an existing project user is a memeber of",
			Flags: []cli.Flag{
				asMandatory(idFlag),
			},
			Before: getProjectValidator,
			Action: getProjectHandler,
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

func getProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
}

func createProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}
	path := ctx.String(name(pathFlag))
	pname := ctx.String(name(nameFlag))
	descr := ctx.String(name(descriptionFlag))
	format := ctx.String(name(fmtFlag))

	if path == "" {
		if format == ".json" {
			path = "./.padlfile.json"
		} else {
			path = "./.padlfile.yaml"
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
	fmt.Printf("project %s:%s initialized successfully!\n", pname, pf.Data.Project)
	return nil
}

func getProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	id := ctx.String(name(idFlag))

	project, err := c.GetProject(id)
	if err != nil {
		return fmt.Errorf("error finding project: %s", err)
	}

	prettyJSON, err := json.MarshalIndent(project, "", "    ")
	if err != nil {
		return fmt.Errorf("Failed to generate json: %s", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
	return nil
}
