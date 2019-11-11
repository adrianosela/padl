package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

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

	priv, pub, err := keys.GenerateRSAKeyPair(2048)
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
	privPem := keys.EncodePrivKeyPEM(priv)
	err = ioutil.WriteFile(fmt.Sprintf("%s.priv", email), privPem, 0644)
	if err != nil {
		return fmt.Errorf("could not write private key to file: %s", err)
	}

	fmt.Printf("registered user %s successfully!\n", email)
	return nil
}
