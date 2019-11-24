package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/olekukonko/tablewriter"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	defaultHostURL = "https://padl.adrianosela.com"
)

// AccountCmds - manage accounts and auth flows
var AccountCmds = cli.Command{
	Name:    "account",
	Aliases: []string{"a"},
	Usage:   "Manage accounts and authentication",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a new account on a padl server",
			Flags: []cli.Flag{
				emailFlag,
				passwordFlag,
			},
			Before: createConfigIfDoesNotExist,
			Action: createAccountHandler,
		},
		{
			Name:  "login",
			Usage: "login to an account on a padl server",
			Flags: []cli.Flag{
				emailFlag,
				passwordFlag,
				pathFlag,
			},
			Before: createConfigIfDoesNotExist,
			Action: loginAccountHandler,
		},
		{
			Name:  "rotate-key",
			Usage: "create a fresh user key and publish the public key",
			Flags: []cli.Flag{
				pathFlag,
			},
			Action: rotateKeyHandler,
		},
		{
			Name:  "show",
			Usage: "show account details based on padl token",
			Flags: []cli.Flag{
				jsonFlag,
			},
			Action: showAccountHandler,
		},
	},
}

func createConfigIfDoesNotExist(ctx *cli.Context) error {
	cPath := ctx.GlobalString(name(ConfigFlag))
	if _, err := config.GetConfig(cPath); err == nil {
		return nil
	}
	c := &config.Config{HostURL: defaultHostURL}
	if err := config.SetConfig(c, cPath); err != nil {
		return fmt.Errorf("could not set configuration: %s", err)
	}
	return nil
}

func createAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	email := ctx.String(name(emailFlag))
	if email == "" {
		if email, err = promptText("Enter your email:", false); err != nil {
			return fmt.Errorf("could not read user email")
		}
	}

	pass := ctx.String(name(passwordFlag))
	if pass == "" {
		if pass, err = promptText("Enter your password:", true); err != nil {
			return fmt.Errorf("could not read user password")
		}
	}

	priv, pub, err := keys.GenerateRSAKeyPair(4096)
	if err != nil {
		return fmt.Errorf("could not generate key pair: %s", err)
	}

	// register user
	err = c.Register(email, pass, string(keys.EncodePubKeyPEM(pub)))
	if err != nil {
		return err
	}

	// save private key in filesystem
	keyMgr, err := keymgr.NewFSManager(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("could not establish key manager: %s", err)
	}

	if err = keyMgr.PutPriv(keys.GetFingerprint(pub), string(keys.EncodePrivKeyPEM(priv))); err != nil {
		return fmt.Errorf("could not save private key: %s", err)
	}

	fmt.Printf("registered user %s successfully!\n", email)
	return nil
}

func rotateKeyHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}
	// create new key
	priv, pub, err := keys.GenerateRSAKeyPair(4096)
	if err != nil {
		return fmt.Errorf("could not generate key pair: %s", err)
	}
	// save private key in filesystem
	keyMgr, err := keymgr.NewFSManager(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("could not establish key manager: %s", err)
	}
	if err = keyMgr.PutPriv(keys.GetFingerprint(pub), string(keys.EncodePrivKeyPEM(priv))); err != nil {
		return fmt.Errorf("could not save private key: %s", err)
	}
	// rotate key
	if err = c.RotateUserKey(string(keys.EncodePubKeyPEM(pub))); err != nil {
		return err
	}
	fmt.Println("rotated user key successfully!")
	return nil
}

func loginAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	path := ctx.String(name(pathFlag))

	email := ctx.String(name(emailFlag))
	if email == "" {
		if email, err = promptText("Enter your email:", false); err != nil {
			return fmt.Errorf("could not read user email")
		}
	}

	pass := ctx.String(name(passwordFlag))
	if pass == "" {
		if pass, err = promptText("Enter your password:", true); err != nil {
			return fmt.Errorf("could not read user password")
		}
	}

	tk, err := c.Login(email, pass)
	if err != nil {
		return fmt.Errorf("could not log in: %s", err)
	}

	conf, err := config.GetConfig(path)
	if err != nil {
		return fmt.Errorf("could not get config from file system: %s", err)
	}

	conf.User = email
	conf.Token = tk
	if err = config.SetConfig(conf, path); err != nil {
		return fmt.Errorf("could not write config to file system: %s", err)
	}

	fmt.Printf("user %s logged in successfully!\n", email)
	return nil
}

func showAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	claims, err := c.Valid()
	if err != nil {
		return fmt.Errorf("could not get token details: %s", err)
	}

	if ctx.Bool(name(jsonFlag)) {
		return printJSON(&claims)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.Append([]string{"aud", claims.Audience})
	table.Append([]string{"iss", claims.Issuer})
	table.Append([]string{"sub", claims.Subject})
	table.Append([]string{"iat", strconv.FormatInt(claims.IssuedAt, 10)})
	table.Append([]string{"exp", strconv.FormatInt(claims.ExpiresAt, 10)})
	table.Append([]string{"jti", claims.Id})
	table.Render()

	return nil
}
