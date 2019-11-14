package commands

import (
	"bufio"
	"fmt"
	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	cli "gopkg.in/urfave/cli.v1"
	"os"
	"os/exec"
)

// RunCmds - run a command with injected secrets
var RunCmds = cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "run a command with secrets in the environment",
	Flags: []cli.Flag{
		jsonFlag,
	},
	Action: runHandler,
}

func runHandler(ctx *cli.Context) error {
	c, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize client: %s", err)
	}

	var cmd *exec.Cmd
	if len(os.Args) > 3 {
		cmd = exec.Command(os.Args[2], os.Args[3:]...)
	} else if len(os.Args) > 2 {
		cmd = exec.Command(os.Args[2])
	} else {
		return fmt.Errorf("no command provided")
	}

	// get secrets for project
	secret, err := c.GetSecrets( /*FIXME*/ )
	if err != nil {
		return fmt.Errorf("could not get secrets from server: %s", err)
	}
	kid := secret.KeyID

	mgr, err := keymgr.NewFSManager(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("could not initialize key manager: %s", err)
	}

	keyPEM, err := mgr.GetKey(kid)
	if err != nil {
		return fmt.Errorf("could not get key %s: %s", kid, err)
	}

	k, err := keys.DecodePrivKeyPEM([]byte(keyPEM))
	if err != nil {
		return fmt.Errorf("could not materialize key %s: %s", kid, err)
	}

	decrypted, err := secret.Decrypt(k)
	if err != nil {
		return fmt.Errorf("could not decrypt secrets: %s", err)
	}

	// copy parent environment
	cmd.Env = os.Environ()
	// attach decrypted secret to the cmd's environment
	for _, s := range decrypted {
		cmd.Env = append(cmd.Env, s)
	}

	return runCmdAndPipeStdout(cmd)
}

func runCmdAndPipeStdout(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	cmd.Start()
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}
