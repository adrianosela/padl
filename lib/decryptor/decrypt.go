package decryptor

import (
	"fmt"

	"github.com/adrianosela/padl/api/client"
	"github.com/adrianosela/padl/lib/keymgr"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/adrianosela/padl/lib/padlfile"
	"github.com/adrianosela/padl/lib/secret"
	"github.com/adrianosela/padl/lib/shamir"
)

// SecretsMgr encrypts/decrypts padlfile secrets
type SecretsMgr struct {
	client     *client.Padl
	keyManager keymgr.Manager
	padlFile   *padlfile.File
}

// NewSecretsMgr is the constructor for the SecretsMgr
func NewSecretsMgr(client *client.Padl, keyMgr keymgr.Manager, pf *padlfile.File) *SecretsMgr {
	return &SecretsMgr{
		client:     client,
		keyManager: keyMgr,
		padlFile:   pf,
	}
}

// DecryptPadlFileSecrets uses the network and the file system to decrypt
// the contents of a padlfile
func (smgr *SecretsMgr) DecryptPadlFileSecrets(filesystemKeyID string) (map[string]string, error) {
	decrypted := make(map[string]string)

	for varName, encrypted := range smgr.padlFile.Data.Variables {
		sec, err := secret.DecodePEM(encrypted)
		if err != nil {
			return nil, fmt.Errorf("could not decode padlfile var: %s, body %s. %s", varName, encrypted, err)
		}

		parts := [][]byte{}

		for _, sh := range sec.Shards {
			if sh.KeyID == smgr.padlFile.Data.SharedKey {
				decryptedSharedShard, err := smgr.client.DecryptSecret(sh.Value, sh.KeyID)
				if err != nil {
					return nil, fmt.Errorf("could not decrypt shared shard for var: %s. %s", varName, err)
				}
				parts = append(parts, []byte(decryptedSharedShard))
			} else if sh.KeyID == filesystemKeyID {
				priv, err := smgr.keyManager.GetPriv(filesystemKeyID)
				if err != nil {
					return nil, fmt.Errorf("could not get private key from filesystem: %s", err)
				}
				k, err := keys.DecodePrivKeyPEM([]byte(priv))
				if err != nil {
					return nil, fmt.Errorf("could not decode RSA private key: %s", err)
				}
				decryptedUserShard, err := keys.DecryptMessage([]byte(sh.Value), k)
				if err != nil {
					return nil, fmt.Errorf("could not decrypt user shard for var: %s. %s", varName, err)
				}
				parts = append(parts, decryptedUserShard)
			}
		}

		if len(parts) < 2 {
			return nil, fmt.Errorf("could not decrypt necessary parts for var: %s", varName)
		}
		plain, err := shamir.Combine(parts)
		if err != nil {
			return nil, fmt.Errorf("could not shamir.Combine decrypted parts for var: %s. %s", varName, err)
		}
		decrypted[varName] = string(plain)
	}

	return decrypted, nil
}

// EncryptPadlfileSecrets uses the network and the file system to encrypt
// the contents of a padlfile
func (smgr *SecretsMgr) EncryptPadlfileSecrets() (map[string]string, error) {
	encrypted := make(map[string]string)

	// FIXME: actually encrypt
	for varName, plain := range smgr.padlFile.Data.Variables {
		// NOP
		encrypted[varName] = plain
	}

	return encrypted, nil
}
