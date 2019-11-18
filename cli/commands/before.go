package commands

import (
	"fmt"

	"github.com/adrianosela/padl/cli/config"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/padlfile"
	cli "gopkg.in/urfave/cli.v1"
)

func checkCanModifyPadlFile(ctx *cli.Context) error {
	// read padlfile
	pf, err := padlfile.ReadPadlfile(ctx.String(name(pathFlag)))
	if err != nil {
		return err
	}
	// check user has decryption key
	key, kid, err := getUserKey(config.GetDefaultPath(), pf.Data.MemberKeys)
	if err != nil {
		return err
	}
	ctx.Set("key", key)
	ctx.Set("keyID", kid)
	return nil
}

func getUserKey(keyMgrPath string, keyIDs []string) (string, string, error) {
	// init new key manager at the config path
	mgr, err := keymgr.NewFSManager(keyMgrPath)
	if err != nil {
		return "", "", err
	}
	// get any private key that matches a ids given
	key, id := "", ""
	for _, k := range keyIDs {
		if key, err = mgr.GetPriv(k); err == nil {
			id = k
			break
		}
	}
	if key == "" {
		return "", "", fmt.Errorf("no valid decryption key found")
	}
	return key, id, nil
}