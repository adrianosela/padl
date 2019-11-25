package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/adrianosela/padl/api/client"
	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/cli/config"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v1"
)

func getClient(ctx *cli.Context) (*client.Padl, error) {
	config, err := config.GetConfig(ctx.GlobalString("config"))
	jwt := ctx.GlobalString("service-account-token")
	hostURL := ctx.GlobalString("host-url")
	if hostURL == "" {
		hostURL = config.HostURL
	}
	if err != nil {
		return nil, err
	}
	return client.NewPadlClient(hostURL, config.Token, jwt, nil)
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
			t.Append([]string{header, fmt.Sprintf("%s %d", k, v)})
			headerSet = true
			continue
		}
		t.Append([]string{"", fmt.Sprintf("%s %d", k, v)})
	}
}

func promptText(prompt string, secret bool) (string, error) {
	fmt.Println(prompt)

	if secret {
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(password)), nil
	}

	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
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
