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
			Usage: "create a server managed key-pair",
			Flags: []cli.Flag{
				jsonFlag,
				asMandatory(nameFlag),
				asMandatory(descriptionFlag),
				withDefaultInt(bitsFlag, 2048),
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
		{
			Name:  "add-user",
			Usage: "add a user to a key-pair's access control list",
			Flags: []cli.Flag{
				jsonFlag,
				asMandatory(idFlag),
				asMandatory(emailFlag),
				withDefaultInt(privFlag, 0),
			},
			Before: addUserToKeyValidator,
			Action: addUserToKeyHandler,
		},
		{
			Name:  "remove-user",
			Usage: "remove a user from a key-pair's access control list",
			Flags: []cli.Flag{
				jsonFlag,
				asMandatory(idFlag),
				asMandatory(emailFlag),
			},
			Before: rmUserToKeyValidator,
			Action: rmUserToKeyHandler,
		},
	},
}

func createKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, nameFlag, descriptionFlag)
}

func privateKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
}

func publicKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag)
}

func addUserToKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag, emailFlag)
}

func rmUserToKeyValidator(ctx *cli.Context) error {
	return assertSet(ctx, idFlag, emailFlag)
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

func addUserToKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	id := ctx.String(name(idFlag))
	email := ctx.String(name(emailFlag))
	privilege := ctx.Int(name(privFlag))

	k, err := c.AddUserToPrivateKey(id, email, privilege)
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

func rmUserToKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	id := ctx.String(name(idFlag))
	email := ctx.String(name(emailFlag))

	k, err := c.RemoveUserFromPrivateKey(id, email)
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
