package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	cli "gopkg.in/urfave/cli.v1"
)

// RunCmds - run a command with injected secrets
var RunCmds = cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "run a command with secrets in the environment",
	Action:  runHandler,
}

func runHandler(ctx *cli.Context) error {
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

	// we append all secrets as follows
	// cmd.Env = append(cmd.Env, "MY_VAR=some_value")

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
