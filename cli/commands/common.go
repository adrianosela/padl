package commands

import (
	"github.com/adrianosela/padl/api/client"
	"github.com/adrianosela/padl/cli/config"
	cli "gopkg.in/urfave/cli.v1"
)

func getClient(ctx *cli.Context) (*client.Padl, error) {
	config, err := config.GetConfig(ctx.GlobalString("config"))
	if err != nil {
		return nil, err
	}
	return client.NewPadlClient(config.HostURL, nil)
}
