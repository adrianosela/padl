package commands

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v1"
)

// KMSCmds - manage padl-managed private and public keys
var KMSCmds = cli.Command{
	Name:    "kms",
	Aliases: []string{"k"},
	Usage:   "manage private and public keys",
	Subcommands: []cli.Command{
		{
			Name:  "public",
			Usage: "get an RSA public key for inspection/encryption",
			Flags: []cli.Flag{
				jsonFlag,
				asMandatory(idFlag),
			},
			Before: publicKeyValidator,
			Action: publicKeyHandler,
		},
		{
			Name:  "decrypt",
			Usage: "decrypt secret with given padl key id",
			Flags: []cli.Flag{
				asMandatory(idFlag),
				asMandatory(secretFlag),
			},
			Before: decryptSecretValidator,
			Action: decryptSecretHandler,
		},
	},
}

func publicKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
}

func decryptSecretValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag, secretFlag)
}

func publicKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	k, err := c.GetPublicKey(ctx.String(name(idFlag)))
	if err != nil {
		return fmt.Errorf("could not get public key: %s", err)
	}

	if ctx.Bool(name(jsonFlag)) {
		return printJSON(&k)
	}

	fmt.Println(k.PEM)
	return nil
}

func decryptSecretHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	kid := ctx.String(name(idFlag))
	secret := ctx.String(name(secretFlag))

	message, err := c.DecryptSecret(secret, kid)
	if err != nil {
		return fmt.Errorf("could not decrypt secret: %s", err)
	}

	fmt.Println(message)

	return nil
}
