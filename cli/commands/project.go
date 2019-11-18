package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
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
				asMandatory(nameFlag),
				jsonFlag,
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
		{
			Name:  "add-secret",
			Usage: "add a secret to a project",
			Flags: []cli.Flag{
				jsonFlag,
			},
			Before: checkCanModifyPadlFile,
			Action: projectAddSecretHandler,
		},
		{
			Name:  "update-secret",
			Usage: "update a secret in a project",
			Flags: []cli.Flag{
				jsonFlag,
			},
			Before: checkCanModifyPadlFile,
			Action: projectUpdateSecretHandler,
		},
		{
			Name:  "remove-secret",
			Usage: "delete a secret to a project",
			Flags: []cli.Flag{
				jsonFlag,
			},
			Before: checkCanModifyPadlFile,
			Action: projectRemoveSecretHandler,
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

	if ctx.Bool(name(jsonFlag)) {
		return printJSON(&project)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Append([]string{"NAME", project.Name})
	table.Append([]string{"DESCRIPTION", project.Description})
	table.Append([]string{"KEY", project.ProjectKey})

	tablePrivsMap(table, "MEMBERS", project.Members)
	tableStringsMap(table, "DEPLOY KEYS", project.DeployKeys)

	table.Render()
	return nil
}

func projectListHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projects, err := c.ListProjects()
	if err != nil {
		return fmt.Errorf("error fetching projects: %s", err)
	}

	if ctx.Bool(name(jsonFlag)) {
		return printJSON(&projects)
	}

	if len(projects.Projects) == 0 {
		fmt.Println("no projects to show :(")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"NAME", "DESCRIPTION"})
	for _, proj := range projects.Projects {
		table.Append([]string{proj.Name, proj.Description})
	}
	table.Render()

	return nil
}

func projectAddSecretHandler(ctx *cli.Context) error {
	return nil
}

func projectUpdateSecretHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func projectRemoveSecretHandler(ctx *cli.Context) error {
	// TODO
	return nil
}
