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
func (smgr *SecretsMgr) DecryptSecret(ciphertext string, usrOrSvcPriv *rsa.PrivateKey) (string, error) {
	sec, err := secret.DecodeSimpleSecret(ciphertext)
	if err != nil {
		return "", fmt.Errorf("could not decode PEM secret %s", err)
	}
	privID := keys.GetFingerprint(&usrOrSvcPriv.PublicKey)

	parts := [][]byte{}
	for _, sh := range sec.Shards {
		if sh.KeyID == smgr.padlFile.Data.SharedKey {
			decryptedSharedShard, err := smgr.client.DecryptSecret(sh.Value, sh.KeyID)
			if err != nil {
				return "", fmt.Errorf("could not decrypt shared shard: %s", err)
			}
			parts = append(parts, []byte(decryptedSharedShard))
		} else if sh.KeyID == privID {
			decryptedUserShard, err := sh.Decrypt(usrOrSvcPriv)
			if err != nil {
				return "", fmt.Errorf("could not decrypt user shard: %s", err)
			}
			parts = append(parts, decryptedUserShard.Value)
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
	// precache necessary encryption keys in the filesystem
	pubs, err := smgr.PrecachePubs()
	if err != nil {
		return "", fmt.Errorf("could not precache public keys: %s", err)
	}
	// establish secret object
	s := secret.Secret{Shards: []*secret.EncryptedShard{}}
	// we begin by splitting the plaintext into two top-level shares
	topLevelParts, err := shamir.Split([]byte(plaintext), 2, 2)
	if err != nil {
		return "", fmt.Errorf("could not split plaintext secret: %s", err)
	}
	// we encrypt one of them with the shared public key
	sharedShard, err := encryptPart(topLevelParts[0], pubs[smgr.padlFile.Data.SharedKey])
	if err != nil {
		return "", fmt.Errorf("could not encrypt shared shard: %s", err)
	}
	s.Shards = append(s.Shards, sharedShard)
	// we encrypt the other top level shard N times (with each of the N user/service keys)
	memberKeys, serviceKeys := smgr.padlFile.Data.MemberKeys, smgr.padlFile.Data.ServiceKeys
	for _, k := range append(memberKeys, serviceKeys...) {
		usrShard, err := encryptPart(topLevelParts[1], pubs[k])
		if err != nil {
			return "", fmt.Errorf("could not encrypt shared shard: %s", err)
		}
		s.Shards = append(s.Shards, usrShard)
	}
	// we then PEM encode the secret data
	return s.EncodeSimple(), nil
}

func encryptPart(part []byte, pub *rsa.PublicKey) (*secret.EncryptedShard, error) {
	plainShard, err := secret.NewShard(part)
	if err != nil {
		return nil, fmt.Errorf("could not build shard: %s", err)
	}
	encShard, err := plainShard.Encrypt(pub)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt shard: %s", err)
	}
	return encShard, nil
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
func (smgr *SecretsMgr) PrecachePubs() (map[string]*rsa.PublicKey, error) {
	pubs := make(map[string]*rsa.PublicKey)
	// for all keys (user keys + service account keys + shared team key)
	keysToFetch := append(
		smgr.padlFile.Data.MemberKeys,
		append(smgr.padlFile.Data.ServiceKeys, smgr.padlFile.Data.SharedKey)...,
	)
	for _, k := range keysToFetch {
		// try to get pub from filesystem
		pubPEM, err := smgr.keyManager.GetPub(k)
		if err == nil {
			pubRSA, err := keys.DecodePubKeyPEM([]byte(pubPEM))
			if err == nil {
				pubs[k] = pubRSA
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
		pubs[k] = pubRSA

		// store it to the file system
		if err := smgr.keyManager.PutPub(pub.ID, pub.PEM); err != nil {
			return nil, fmt.Errorf("could not put pub %s in file system: %s", k, err)
		}
	}

	return pubs, nil
}
