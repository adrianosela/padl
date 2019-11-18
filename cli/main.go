package main

import (
	"fmt"
	"os"

	"github.com/adrianosela/padl/cli/commands"
	cli "gopkg.in/urfave/cli.v1"
)

var version string // injected at build-time

var appflags = []cli.Flag{
	commands.ConfigFlag,
	commands.VerboseFlag,
}

var appcmds = []cli.Command{
	commands.ConfigCmd,
	commands.AccountCmds,
	commands.ProjectCmds,
	commands.KMSCmds,
	commands.RunCmds,
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Usage = "secrets management as-a-service"
	app.Flags = appflags
	app.Commands = appcmds
	app.CommandNotFound = func(c *cli.Context, command string) {
		c.App.Run([]string{"help"})
		fmt.Printf("\ncommand \"%s\" does not exist\n", command)
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
