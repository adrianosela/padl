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
	// config flags
	pathFlag = cli.StringFlag{
		Name:  "path, p",
		Usage: "override default config file path",
	}
	urlFlag = cli.StringFlag{
		Name:  "url, u",
		Usage: "host URL to use",
	}

	// option flags
	jsonFlag = cli.BoolFlag{
		Name:  "json, j",
		Usage: "print raw json -- don't pretty print",
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
