package commands

import (
	"fmt"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/adrianosela/padl/lib/padlfile"
	"github.com/adrianosela/padl/lib/secretsmgr"
	cli "gopkg.in/urfave/cli.v1"
)

// PadlfileCmds - manage padlfile
var PadlfileCmds = cli.Command{
	Name:    "file",
	Aliases: []string{"f"},
	Usage:   "manage padlfile",
	Subcommands: []cli.Command{
		{
			Name:  "pull",
			Usage: "update padlfile to match server state",
			Flags: []cli.Flag{
				pathFlag,
			},
			// TODO: before?
			Action: padlfilePullHandler,
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
					Before: padlfileSetSecretValidator,
					Action: padlfileSetSecretHandler,
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
					Before: padlfileShowSecretValidator,
					Action: padlfileShowSecretHandler,
				},
				{
					Name:  "remove",
					Usage: "delete a secret from a project",
					Flags: []cli.Flag{
						asMandatory(nameFlag),
						withDefault(fmtFlag, "yaml"),
						privateKeyFlag, // set by BeforeFunc
						pathFlag,
					},
					Before: func(ctx *cli.Context) error {
						if err := checkCanModifyPadlFile(ctx); err != nil {
							return err
						}
						return padlfileRemoveSecretValidator(ctx)
					},
					Action: padlfileRemoveSecretHandler,
				},
			},
		},
	},
}

func padlfileSetSecretValidator(ctx *cli.Context) error {
	if err := checkCanModifyPadlFile(ctx); err != nil {
		return err
	}
	return assertSet(ctx, nameFlag, secretFlag)
}

func padlfileShowSecretValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag)
}

func padlfileRemoveSecretValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag)
}

func padlfilePullHandler(ctx *cli.Context) error {
	// TODO
	return nil
}

func padlfileSetSecretHandler(ctx *cli.Context) error {
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

func padlfileShowSecretHandler(ctx *cli.Context) error {
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

func padlfileRemoveSecretHandler(ctx *cli.Context) error {
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
