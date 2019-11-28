package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/adrianosela/padl/lib/keys"
	cli "gopkg.in/urfave/cli.v1"
)

// ProjectCmds - manage projects
var ProjectCmds = cli.Command{
	Name:    "project",
	Aliases: []string{"p"},
	Usage:   "Manage projects",
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
			Name:  "delete",
			Usage: "delete a padl project",
			Flags: []cli.Flag{
				asMandatory(projectFlag),
			},
			Before: deleteProjectValidator,
			Action: deleteProjectHandler,
		},
		{
			Name:  "get",
			Usage: "get a padl project by name",
			Flags: []cli.Flag{
				asMandatory(projectFlag),
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
			Name:  "service-account",
			Usage: "manage secrets for project",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create a service account in the project",
					Flags: []cli.Flag{
						asMandatory(projectFlag),
						asMandatory(nameFlag),
						jsonFlag,
					},
					Before: addServiceAccountValidator,
					Action: projectAddServiceAccountHandler,
				},
				{
					Name:  "remove",
					Usage: "delete a service account from the project",
					Flags: []cli.Flag{
						asMandatory(projectFlag),
						asMandatory(nameFlag),
					},
					Before: removeServiceAccountValidator,
					Action: projectRemoveServiceAccountHandler,
				},
			},
		},
		{

			Name:  "user",
			Usage: "manage users for project",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a user to a project",
					Flags: []cli.Flag{
						asMandatory(projectFlag),
						asMandatory(emailFlag),
						withDefaultInt(privFlag, 0),
					},
					Before: addUserValidator,
					Action: addUserHandler,
				},
				{
					Name:  "remove",
					Usage: "remove a user from a project",
					Flags: []cli.Flag{
						asMandatory(projectFlag),
						asMandatory(emailFlag),
					},
					Before: removeUserValidator,
					Action: removeUserHandler,
				},
			},
		},
	},
}

func addUserValidator(ctx *cli.Context) error {
	return assertSet(ctx, projectFlag, emailFlag, privFlag)
}

func addServiceAccountValidator(ctx *cli.Context) error {
	return assertSet(ctx, projectFlag, nameFlag)
}

func removeServiceAccountValidator(ctx *cli.Context) error {
	return assertSet(ctx, projectFlag, nameFlag)
}

func createProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, descriptionFlag)
}

func deleteProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, projectFlag)
}

func getProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, projectFlag)
}

func removeUserValidator(ctx *cli.Context) error {
	return assertSet(ctx, projectFlag, emailFlag)
}

func createProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	pname := ctx.String(name(nameFlag))
	bits := ctx.Int(name(bitsFlag))
	descr := ctx.String(name(descriptionFlag))
	format := ctx.String(name(fmtFlag))

	path := padlfilePath(ctx.String(name(pathFlag)), format)

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

	projectName := ctx.String(name(projectFlag))

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
	tableStringsMap(table, "SERVICE ACCOUNTS", project.ServiceAccounts)

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

func projectAddServiceAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(projectFlag))
	keyName := ctx.String(name(nameFlag))

	priv, pub, err := keys.GenerateRSAKeyPair(4096)
	if err != nil {
		return fmt.Errorf("could not generate key pair: %s", err)
	}

	resp, err := c.CreateServiceAccount(projectName, keyName, string(keys.EncodePubKeyPEM(pub)))
	if err != nil {
		return fmt.Errorf("error creating service account: %s", err)
	}

	data := struct {
		PrivateKey string `json:"private_key"`
		JWT        string `json:"jwt"`
	}{
		PrivateKey: string(keys.EncodePrivKeyPEM(priv)),
		JWT:        resp.Token,
	}

	if ctx.Bool(name(jsonFlag)) {
		return printJSON(data)
	}

	fmt.Println("---------------------- IMPORTANT NOTE ----------------------")
	fmt.Println(">> Both the RSA private key and and auth token are secret <<")
	fmt.Println(">> If either is disclosed you MUST delete the svc account <<")
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("\nSERVICE ACCOUNT PRIVATE KEY:\n%s", data.PrivateKey)
	fmt.Printf("\nSERVICE ACCOUNT AUTH TOKEN:\n%s\n", data.JWT)

	return nil
}

func projectRemoveServiceAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(projectFlag))
	keyName := ctx.String(name(nameFlag))

	err = c.RemoveServiceAccount(projectName, keyName)
	if err != nil {
		return fmt.Errorf("error removing service account: %s", err)
	}
	return nil
}

func addUserHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(projectFlag))
	email := ctx.String(name(emailFlag))
	privLevel := ctx.Int(name(privFlag))

	if err := c.AddUserToProject(projectName, email, privLevel); err != nil {
		return fmt.Errorf("error adding user: %s", err)
	}
	fmt.Printf("user %s added to project %s successfully!\n", email, projectName)
	return nil
}

func removeUserHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(projectFlag))
	email := ctx.String(name(emailFlag))

	if err = c.RemoveUserFromProject(projectName, email); err != nil {
		return fmt.Errorf("error removing user: %s", err)
	}
	fmt.Printf("user %s removed from project %s successfully!\n", email, projectName)
	return nil
}

func deleteProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(projectFlag))
	if err := c.DeleteProject(projectName); err != nil {
		return fmt.Errorf("error deleting project: %s", err)
	}
	fmt.Printf("project %s deleted successfully!\n", projectName)
	return nil
}
