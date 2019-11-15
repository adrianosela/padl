package commands

import (
	"encoding/json"
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
	},
}

func publicKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
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
		byt, err := json.Marshal(&k)
		if err != nil {
			return fmt.Errorf("could not marshal response: %s", err)
		}
		fmt.Println(string(byt))
		return nil
	}

	fmt.Println(k.PEM)
	return nil
}
