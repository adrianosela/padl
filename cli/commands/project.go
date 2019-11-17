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
			Name:  "create",
			Usage: "create a new padl project",
			Flags: []cli.Flag{
				asMandatory(nameFlag),
				asMandatory(descriptionFlag),
				withDefaultInt(bitsFlag, 2048),
				pathFlag,
				withDefault(fmtFlag, "yaml"),
			},
			Before: createProjectValidator,
			Action: createProjectHandler,
		},
		{
			Name:  "get",
			Usage: "get a padl project by name",
			Flags: []cli.Flag{
				asMandatory(idFlag),
				asMandatory(nameFlag),
			},
			Before: getProjectValidator,
			Action: getProjectHandler,
		},
		{
			Name:  "list",
			Usage: "get all your padl projects",
			Flags: []cli.Flag{
				jsonFlag,
			},
			Action: projectListHandler,
		},
	},
}

func createProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, descriptionFlag)
}

func getProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag)
}

func createProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	path := ctx.String(name(pathFlag))
	pname := ctx.String(name(nameFlag))
	bits := ctx.Int(name(bitsFlag))
	descr := ctx.String(name(descriptionFlag))
	format := ctx.String(name(fmtFlag))

	if path == "" {
		if format == "yaml" {
			path = "./.padlfile.yaml"
		} else {
			path = "./.padlfile.json"
		}
	}

	if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".yaml") {
		return fmt.Errorf("invalid file extension, must be one of { \".yaml\", \".json\" }")
	}

	pf, err := c.CreateProject(pname, descr, bits)
	if err != nil {
		return fmt.Errorf("error creating project: %s", err)
	}

	err = pf.Write(path)
	if err != nil {
		return fmt.Errorf("unable to write padl file: %s", err)
	}
	fmt.Printf("project %s initialized successfully!\n", pname)
	return nil
}

func getProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(nameFlag))

	project, err := c.GetProject(projectName)

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

func projectListHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	jf := ctx.Bool(name(jsonFlag))

	projects, err := c.ListProjects()
	if err != nil {
		return fmt.Errorf("error fetching projects: %s", err)
	}

	if jf {
		byt, err := json.Marshal(&projects)
		if err != nil {
			return fmt.Errorf("error printing json: %s", err)
		}
		fmt.Println(string(byt))
	} else {
		prettyJSON, err := json.MarshalIndent(projects, "", "    ")
		if err != nil {
			return fmt.Errorf("error pretty printing json: %s", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	}

	return nil
}
