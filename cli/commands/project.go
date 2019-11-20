package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/adrianosela/padl/lib/padlfile"
	"github.com/adrianosela/padl/lib/secretsmgr"
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
			Name:  "delete",
			Usage: "delete a padl project",
			Flags: []cli.Flag{
				asMandatory(nameFlag),
			},
			Before: deleteProjectValidator,
			Action: deleteProjectHandler,
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
			Name:  "secret",
			Usage: "manage secrets for project",
			Subcommands: []cli.Command{
				{
					Name:  "set",
					Usage: "set a secret in a project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
						asMandatory(secretFlag),
						withDefault(fmtFlag, "yaml"),
						privateKeyFlag, // set by BeforeFunc
						pathFlag,
					},
					Before: checkCanModifyPadlFile,
					Action: projectSetSecretHandler,
				},
				{
					Name:  "show",
					Usage: "see a specific secret in a project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
						withDefault(fmtFlag, "yaml"),
						privateKeyFlag, // set by BeforeFunc
						pathFlag,
					},
					Before: checkCanModifyPadlFile,
					Action: projectShowSecretHandler,
				},
				{
					Name:  "remove",
					Usage: "delete a secret from a project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
						asMandatory(secretFlag),
						withDefault(fmtFlag, "yaml"),
						privateKeyFlag, // set by BeforeFunc
						pathFlag,
					},
					Before: checkCanModifyPadlFile,
					Action: projectRemoveSecretHandler,
				},
			},
		},
		{
			Name:  "deploykey",
			Usage: "manage secrets for project",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a deploy key to the project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
						asMandatory(keyNameFlag),
						descriptionFlag,
					},
					Before: checkCanModifyPadlFile,
					Action: projectAddDeployKeyHandler,
				},
				{
					Name:  "remove",
					Usage: "remove a deploy key from the project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
						asMandatory(keyNameFlag),
					},
					Before: checkCanModifyPadlFile,
					Action: projectRemoveDeployKeyHandler,
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
						asMandatory(nameFlag),
						asMandatory(emailFlag),
						asMandatoryInt(privFlag),
					},
					Before: addUserValidator,
					Action: addUserHandler,
				},
				{
					Name:  "remove",
					Usage: "remove a user from a project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
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
	return assertSet(ctx, nameFlag, emailFlag, privFlag)
}

func createProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, descriptionFlag)
}

func deleteProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag)
}

func getProjectValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag)
}

func removeUserValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, emailFlag)
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

func projectSetSecretHandler(ctx *cli.Context) error {
	sName := ctx.String(name(nameFlag))
	plaintext := ctx.String(name(secretFlag))
	format := ctx.String(name(fmtFlag))
	path := padlfilePath(ctx.String(name(pathFlag)), format)

	// get client
	pc, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not get client: %s", err)
	}
	// read padlfile
	pf, err := padlfile.ReadPadlfile(path)
	if err != nil {
		return fmt.Errorf("could not read padlfile: %s", err)
	}

	// get key panager
	keyMgr, err := keymgr.NewFSManager(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("could not establish key manager: %s", err)
	}
	secMgr := secretsmgr.NewSecretsMgr(pc, keyMgr, pf)
	// encrypt secret and add to padlfile
	encrypted, err := secMgr.EncryptSecret(plaintext)
	if err != nil {
		return fmt.Errorf("could not encrypt secret %s: %s", sName, err)
	}
	pf.Data.Variables[sName] = encrypted
	if err = pf.Write(path); err != nil {
		return fmt.Errorf("could not write padlfile: %s", err)
	}
	fmt.Println("padlfile updated!")
	return nil
}

func projectShowSecretHandler(ctx *cli.Context) error {
	sName := ctx.String(name(nameFlag))
	format := ctx.String(name(fmtFlag))
	priv := ctx.String(name(privateKeyFlag))
	path := padlfilePath(ctx.String(name(pathFlag)), format)

	// get client
	pc, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not get client: %s", err)
	}
	// read padlfile
	pf, err := padlfile.ReadPadlfile(path)
	if err != nil {
		return fmt.Errorf("could not read padlfile: %s", err)
	}

	if _, ok := pf.Data.Variables[sName]; !ok {
		return fmt.Errorf("secret %s not in padlfile", sName)
	}
	// get key panager
	keyMgr, err := keymgr.NewFSManager(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("could not establish key manager: %s", err)
	}
	secMgr := secretsmgr.NewSecretsMgr(pc, keyMgr, pf)
	// decrypted secret and print it
	rsa, err := keys.DecodePrivKeyPEM([]byte(priv))
	if err != nil {
		return fmt.Errorf("could not materialize user private key: %s", err)
	}
	decrypted, err := secMgr.DecryptSecret(pf.Data.Variables[sName], rsa)
	if err != nil {
		return fmt.Errorf("could not decrypt secret %s: %s", sName, err)
	}
	fmt.Println(decrypted)
	return nil
}

func projectRemoveSecretHandler(ctx *cli.Context) error {
	sName := ctx.String(name(nameFlag))
	format := ctx.String(name(fmtFlag))
	path := padlfilePath(ctx.String(name(pathFlag)), format)

	// read padlfile
	pf, err := padlfile.ReadPadlfile(path)
	if err != nil {
		return fmt.Errorf("could not read padlfile: %s", err)
	}

	if _, ok := pf.Data.Variables[sName]; !ok {
		return fmt.Errorf("secret %s not in padlfile", sName)
	}
	// delete var
	delete(pf.Data.Variables, sName)
	// write padlfile
	if err = pf.Write(path); err != nil {
		return fmt.Errorf("could not write padlfile: %s", err)
	}
	fmt.Println("padlfile updated!")
	return nil
}

func projectAddDeployKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(nameFlag))
	keyName := ctx.String(name(keyNameFlag))
	description := ctx.String(name(descriptionFlag))

	resp, err := c.CreateDeployKey(projectName, keyName, description)
	if err != nil {
		return fmt.Errorf("error creating a deploy key: %s", err)
	}

	fmt.Println(resp.Token)

	//TODO Generate a user "deploy" key

	//TODO Generate a file to store token and key

	return nil
}

func projectAddRemoveKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(nameFlag))
	keyName := ctx.String(name(keyNameFlag))

	ok, err := c.RemoveDeployKey(projectName, keyName)
	if err != nil {
		return fmt.Errorf("error reomiving a deploy key: %s", err)
	}

	fmt.Println(ok)

	return nil
}

func projectRemoveDeployKeyHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func addUserHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(nameFlag))
	email := ctx.String(name(emailFlag))
	privLevel := ctx.Int(name(privFlag))

	ok, err := c.AddUserToProject(projectName, email, privLevel)
	if err != nil {
		return fmt.Errorf("error adding user: %s", err)
	}
	fmt.Println(ok)
	return nil
}

func removeUserHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(nameFlag))
	email := ctx.String(name(emailFlag))

	ok, err := c.RemoveUserFromProject(projectName, email)
	if err != nil {
		return fmt.Errorf("error removing user: %s", err)
	}
	fmt.Println(ok)
	return nil
}

func deleteProjectHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	projectName := ctx.String(name(nameFlag))
	ok, err := c.DeleteProject(projectName)
	if err != nil {
		return fmt.Errorf("error deleting project: %s", err)
	}
	fmt.Println(ok)
	return nil
}

func padlfilePath(path, fmt string) string {
	if path == "" {
		if fmt == "yaml" {
			return "./.padlfile.yaml"
		}
		return "./.padlfile.json"
	}
	return path
}
