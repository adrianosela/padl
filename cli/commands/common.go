package commands

import (
	"encoding/json"
	"fmt"

	"github.com/adrianosela/padl/api/client"
	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/cli/config"
	"github.com/olekukonko/tablewriter"
	cli "gopkg.in/urfave/cli.v1"
)

func getClient(ctx *cli.Context) (*client.Padl, error) {
	config, err := config.GetConfig(ctx.GlobalString("config"))
	if err != nil {
		return nil, err
	}
	return client.NewPadlClient(config.HostURL, config.Token, nil)
}

func printJSON(i interface{}) error {
	byt, err := json.Marshal(&i)
	if err != nil {
		return fmt.Errorf("error printing json: %s", err)
	}
	fmt.Println(string(byt))
	return nil
}

func tableStringsMap(t *tablewriter.Table, header string, m map[string]string) {
	headerSet := false
	for k, v := range m {
		if !headerSet {
			t.Append([]string{header, fmt.Sprintf("%s : %s", k, v)})
			headerSet = true
			continue
		}
		t.Append([]string{"", fmt.Sprintf("%s : %s", k, v)})
	}
}

func tablePrivsMap(t *tablewriter.Table, header string, m map[string]privilege.Level) {
	headerSet := false
	for k, v := range m {
		if !headerSet {
			t.Append([]string{header, fmt.Sprintf("%s : %d", k, v)})
			headerSet = true
			continue
		}
		t.Append([]string{"", fmt.Sprintf("%s : %d", k, v)})
	}
}

func padlfilePath(path, fmt string) string {
	if path == "" {
		if fmt == "yaml" {
			return "./.padlfile.yaml"
		}
		return "./.padlfile.json"
	}
	return path
}
