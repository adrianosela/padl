package commands

import (
	"github.com/adrianosela/padl/api/client"
	"github.com/adrianosela/padl/cli/config"
	cli "gopkg.in/urfave/cli.v1"
)

func getClient(ctx *cli.Context) (*client.Padl, error) {
	cPath := ctx.GlobalString("config")
	config, err := config.GetConfig(cPath)
	if err != nil {
		return nil, err
	}
	return client.NewPadlClient(config.HostURL, config.AuthTK, nil)
}
