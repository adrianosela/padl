package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/adrianosela/padl/lib/padlfile"
	"github.com/adrianosela/padl/lib/secretsmgr"
	cli "gopkg.in/urfave/cli.v1"
)

// RunCmds - run a command with injected secrets
var RunCmds = cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "Run a command with secrets in the environment",
	Flags: []cli.Flag{
		asMandatory(nameFlag),
		withDefault(fmtFlag, "yaml"),
		privateKeyFlag, // set by BeforeFunc
		pathFlag,
	},
	Before: checkCanModifyPadlFile,
	Action: runHandler,
}

func runHandler(ctx *cli.Context) error {
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
	// get key manager
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

	secretsMap, err := secMgr.DecryptPadlFileSecrets(rsa)
	if err != nil {
		return fmt.Errorf("could not decrypt padlfile secrets: %s", err)
	}

	var cmd *exec.Cmd
	if len(os.Args) > 3 {
		cmd = exec.Command(os.Args[2], os.Args[3:]...)
	} else if len(os.Args) > 2 {
		cmd = exec.Command(os.Args[2])
	} else {
		return fmt.Errorf("no command provided")
	}

	// copy parent environment
	cmd.Env = os.Environ()
	// attach decrypted secret to the cmd's environment
	for k, v := range secretsMap {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
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
