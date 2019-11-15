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
			Name:  "create",
			Usage: "create a server managed private key for decryption",
			Flags: []cli.Flag{
				jsonFlag,
				asMandatory(nameFlag),
				asMandatory(descriptionFlag),
				asMandatoryInt(bitsFlag),
			},
			Before: createKeyValidator,
			Action: createKeyHandler,
		},
		{
			Name:  "private",
			Usage: "get an RSA private key for inspection/decryption",
			Flags: []cli.Flag{
				jsonFlag,
				asMandatory(idFlag),
			},
			Before: privateKeyValidator,
			Action: privateKeyHandler,
		},
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

func createKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, descriptionFlag, bitsFlag)
}

func privateKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
}

func publicKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
}

func createKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	kname := ctx.String(name(nameFlag))
	descr := ctx.String(name(descriptionFlag))
	bits := ctx.Int(name(bitsFlag))

	k, err := c.CreateKey(kname, descr, bits)
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

	fmt.Println(k.ID)
	return nil
}

func privateKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	k, err := c.GetPrivateKey(ctx.String(name(idFlag)))
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
