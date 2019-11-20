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
func (smgr *SecretsMgr) DecryptPadlFileSecrets(userPriv *rsa.PrivateKey) (map[string]string, error) {
	decrypted := make(map[string]string)

	for varName, encrypted := range smgr.padlFile.Data.Variables {
		plain, err := smgr.DecryptSecret(encrypted, userPriv)
		if err != nil {
			return nil, fmt.Errorf("could not decrypt secret for var %s: %s", varName, err)
		}
		decrypted[varName] = string(plain)
	}

	return decrypted, nil
}

// DecryptSecret ecrypts a single pem encoded secret
func (smgr *SecretsMgr) DecryptSecret(ciphertext string, userPriv *rsa.PrivateKey) (string, error) {
	sec, err := secret.DecodePEM(ciphertext)
	if err != nil {
		return "", fmt.Errorf("could not decode PEM secret %s", err)
	}

	parts := [][]byte{}
	for _, sh := range sec.Shards {
		if sh.KeyID == smgr.padlFile.Data.SharedKey {
			decryptedSharedShard, err := smgr.client.DecryptSecret(sh.Value, sh.KeyID)
			if err != nil {
				return "", fmt.Errorf("could not decrypt shared shard: %s", err)
			}
			parts = append(parts, []byte(decryptedSharedShard))
		} else if sh.KeyID == keys.GetFingerprint(&userPriv.PublicKey) {
			decryptedUserShard, err := keys.DecryptMessage([]byte(sh.Value), userPriv)
			if err != nil {
				return "", fmt.Errorf("could not decrypt user shard: %s", err)
			}
			parts = append(parts, decryptedUserShard)
		}
	}

	if len(parts) < 2 {
		return "", fmt.Errorf("could not decrypt necessary parts for var")
	}
	plain, err := shamir.Combine(parts)
	if err != nil {
		return "", fmt.Errorf("could not shamir.Combine decrypted parts: %s", err)
	}

	return string(plain), nil
}

// EncryptSecret encrypts a single secret
func (smgr *SecretsMgr) EncryptSecret(plaintext string) (string, error) {
	pubs, err := smgr.PrecachePubs()
	if err != nil {
		return "", fmt.Errorf("could not precache public keys: %s", err)
	}

	// we need as many parts as we have keys
	parts, err := shamir.Split([]byte(plaintext), len(smgr.padlFile.Data.MemberKeys)+1, 2)
	if err != nil {
		return "", fmt.Errorf("could not split plaintext secret: %s", err)
	}

	s := secret.Secret{Shards: []*secret.EncryptedShard{}}

	for i, part := range parts {
		plainShard, err := secret.NewShard(part)
		if err != nil {
			return "", fmt.Errorf("could not build shard: %s", err)
		}
		encShard, err := plainShard.Encrypt(pubs[i])
		if err != nil {
			return "", fmt.Errorf("could not encrypt shard: %s", err)
		}
		s.Shards = append(s.Shards, encShard)
	}

	padlPEMSecret, err := s.EncodePEM()
	if err != nil {
		return "", fmt.Errorf("could not PEM encode secret: %s", err)
	}

	return padlPEMSecret, nil
}

// EncryptPadlfileSecrets uses the network and the file system to encrypt
// the contents of a padlfile
func (smgr *SecretsMgr) EncryptPadlfileSecrets() (map[string]string, error) {
	var err error
	encrypted := make(map[string]string)
	for varName, plaintext := range smgr.padlFile.Data.Variables {
		if encrypted[varName], err = smgr.EncryptSecret(plaintext); err != nil {
			return nil, fmt.Errorf("could not encrypt var %s: %s", varName, err)
		}
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
