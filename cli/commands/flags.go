package commands

import (
	"fmt"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	mandatoryTag = "[mandatory]"

	payloadFormatJSON = "json"
	payloadFormatList = "list"
)

var (
	// ConfigFlag is the flag for the padl config file path
	ConfigFlag = cli.StringFlag{
		Name:  "config, c",
		Usage: "override default config file path",
	}
	// VerboseFlag enables/disables process logs
	VerboseFlag = cli.BoolFlag{
		Name:  "verbose, vv",
		Usage: "print process logs",
	}
	// JWTFlag service account JWT token
	JWTFlag = cli.StringFlag{
		Name:  "service-account-token",
		Usage: "override defult token with service jwt",
	}
	// HostURL host server url
	HostURL = cli.StringFlag{
		Name:  "host-url",
		Usage: "override default host url",
	}
	urlFlag = cli.StringFlag{
		Name:  "url, u",
		Usage: "host URL to use",
	}
	pathFlag = cli.StringFlag{
		Name:  "path, p",
		Usage: "override default path",
	}
	emailFlag = cli.StringFlag{
		Name:  "email, e",
		Usage: "user email for authentication",
	}
	passwordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "user password for authentication",
	}
	idFlag = cli.StringFlag{
		Name:  "id",
		Usage: "resource id",
	}
	nameFlag = cli.StringFlag{
		Name:  "name",
		Usage: "resource name",
	}
	projectFlag = cli.StringFlag{
		Name:  "project",
		Usage: "project name",
	}
	keyNameFlag = cli.StringFlag{
		Name:  "key-name",
		Usage: "service account name",
	}
	descriptionFlag = cli.StringFlag{
		Name:  "description",
		Usage: "resource description",
	}
	bitsFlag = cli.IntFlag{
		Name:  "bits",
		Usage: "key bit size - one of { 512, 1024, 2048, 4096 }",
	}
	privFlag = cli.IntFlag{
		Name:  "privilege",
		Usage: "privilege level - { Reader: 0, Editor:1, Owner:2 }",
	}
	secretFlag = cli.StringFlag{
		Name:  "secret",
		Usage: "secret to decrypt",
	}
	jsonFlag = cli.BoolFlag{
		Name:  "json, j",
		Usage: "print raw json -- don't pretty print",
	}
	fmtFlag = cli.StringFlag{
		Name:  "fmt",
		Usage: "preferred padlfile format - one of { \"yaml\", \"json\" }",
	}
	privateKeyFlag = cli.StringFlag{
		Name:  "private-key, k",
		Usage: "provide a (user's) private key to decrypt",
	}
)

// name returns the long name of a flag
// note that the split function returns the original string in index 0
// if it does not contain the given delimiter ","
func name(f cli.Flag) string {
	return strings.Split(f.GetName(), ",")[0]
}

func withDefault(f cli.StringFlag, def string) cli.StringFlag {
	f.Value = def
	return f
}

func withDefaultInt(f cli.IntFlag, def int) cli.IntFlag {
	f.Value = def
	return f
}

func asMandatory(f cli.StringFlag) cli.StringFlag {
	f.Usage = fmt.Sprintf("%s %s", mandatoryTag, f.Usage)
	return f
}

func asMandatoryInt(f cli.IntFlag) cli.IntFlag {
	f.Usage = fmt.Sprintf("%s %s", mandatoryTag, f.Usage)
	return f
}

func asMandatoryIf(f cli.StringFlag, cond string) cli.StringFlag {
	f.Usage = fmt.Sprintf("[mandatory if %s] %s", cond, f.Usage)
	return f
}

func assertSet(ctx *cli.Context, flags ...cli.Flag) error {
	for _, f := range flags {
		if !ctx.IsSet(name(f)) {
			return fmt.Errorf("missing %s argument \"%s\"", mandatoryTag, name(f))
		}
	}
	return nil
}

func assertSetIf(ctx *cli.Context, cond func() bool, flags ...cli.Flag) error {
	if !cond() {
		return nil
	}
	return assertSet(ctx, flags...)
}
