package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v1"
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
			Action: loginAccountHandler,
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

func createAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	email := ctx.String(name(emailFlag))
	if email == "" {
		fmt.Println("Enter your email:")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return err
		}
		email = strings.TrimSpace(line)
	}

	pass := ctx.String(name(passwordFlag))
	if pass == "" {
		fmt.Println("Enter your password:")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		pass = strings.TrimSpace(string(password))
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

func loginAccountHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	path := ctx.String(name(pathFlag))

	email := ctx.String(name(emailFlag))
	if email == "" {
		fmt.Println("Enter your email:")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return err
		}
		email = strings.TrimSpace(line)
	}

	pass := ctx.String(name(passwordFlag))
	if pass == "" {
		fmt.Println("Enter your password:")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		pass = strings.TrimSpace(string(password))
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
