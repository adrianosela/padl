package secretsmgr

import (
	"crypto/rsa"
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
	pubs, err := smgr.PrecachePubs()
	if err != nil {
		return nil, fmt.Errorf("could not precache public keys: %s", err)
	}

	encrypted := make(map[string]string)

	for varName, plaintext := range smgr.padlFile.Data.Variables {
		// we need as many parts as we have keys
		parts, err := shamir.Split([]byte(plaintext), len(smgr.padlFile.Data.MemberKeys)+1, 2)
		if err != nil {
			return nil, fmt.Errorf("could not split plaintext secret for %s: %s", varName, err)
		}

		s := secret.Secret{Shards: []*secret.EncryptedShard{}}

		for i, part := range parts {
			plainShard, err := secret.NewShard(part)
			if err != nil {
				return nil, fmt.Errorf("could not build shard for %s: %s", varName, err)
			}
			encShard, err := plainShard.Encrypt(pubs[i])
			if err != nil {
				return nil, fmt.Errorf("could not encrypt shard for %s: %s", varName, err)
			}
			s.Shards = append(s.Shards, encShard)
		}

		padlPEMSecret, err := s.EncodePEM()
		if err != nil {
			return nil, fmt.Errorf("could not PEM encode secret for %s: %s", varName, err)
		}

		// add encrypted/encoded secret to padlfile
		encrypted[varName] = padlPEMSecret
	}

	return encrypted, nil
}

// PrecachePubs gets (and caches) all public keys needed to encrypt
// a padlfile
func (smgr *SecretsMgr) PrecachePubs() ([]*rsa.PublicKey, error) {
	pubs := []*rsa.PublicKey{}
	// for all keys (user keys + shared team key)
	for _, k := range append(smgr.padlFile.Data.MemberKeys, smgr.padlFile.Data.SharedKey) {
		// try to get pub from filesystem
		pubPEM, err := smgr.keyManager.GetPub(k)
		if err == nil {
			pubRSA, err := keys.DecodePubKeyPEM([]byte(pubPEM))
			if err == nil {
				pubs = append(pubs, pubRSA)
				continue // key already in fs
			}
			// fall back to server
		}

		// get public key from padl server
		pub, err := smgr.client.GetPublicKey(k)
		if err != nil {
			return nil, fmt.Errorf("could not get key %s from padl server: %s", k, err)
		}

		pubRSA, err := pub.PubRSA()
		if err != nil {
			return nil, fmt.Errorf("unable to materialize padl pub key onto RSA public key: %s", err)
		}
		pubs = append(pubs, pubRSA)

		// store it to the file system
		if err := smgr.keyManager.PutPub(pub.ID, pub.PEM); err != nil {
			return nil, fmt.Errorf("could not put pub %s in file system: %s", k, err)
		}
	}

	return pubs, nil
}
