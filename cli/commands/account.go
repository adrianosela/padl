package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v1"
)

// AccountCmds - manage accounts and auth flows
var AccountCmds = cli.Command{
	Name:    "account",
	Aliases: []string{"a"},
	Usage:   "manage accounts and auth flows",
	Subcommands: []cli.Command{
		{
			Name:  "register",
			Usage: "register a new account on a padl server",
			Flags: []cli.Flag{
				emailFlag,
				passwordFlag,
			},
			Action: registerAccountHandler,
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
	},
}

func registerAccountHandler(ctx *cli.Context) error {
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
	// TODO: password encrypt
	keyMgr, err := keymgr.NewFSManager(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("could not establish key manager: %s", err)
	}

	fName := fmt.Sprintf("%s.priv", strings.ReplaceAll(keys.GetFingerprint(pub), ":", ""))
	if err = keyMgr.PutKey(fName, string(keys.EncodePrivKeyPEM(priv))); err != nil {
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
		return err
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
