package kms

import (
	"crypto/rsa"
	"fmt"

	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/lib/keys"
)

// Key represents a key managed by padl
type Key struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Users       map[string]privilege.Level `json:"users"`
	PEM         string                     `json:"pem"`
}

// NewKey is the constructor for the padl Key object
func NewKey(bits int, creator, name, descr string) (*Key, error) {
	priv, pub, err := keys.GenerateRSAKeyPair(bits)
	if err != nil {
		return nil, fmt.Errorf("could not generate RSA key: %s", err)
	}
	return &Key{
		ID:          keys.GetFingerprint(pub),
		Name:        name,
		Description: descr,
		PEM:         string(keys.EncodePrivKeyPEM(priv)),
		Users: map[string]privilege.Level{
			creator: privilege.PrivilegeLvlOwner,
		},
	}, nil
}

// Priv returns the materialized RSA private key corresponding
// to the padl-managed Key object
func (k *Key) Priv() (*rsa.PrivateKey, error) {
	return keys.DecodePrivKeyPEM([]byte(k.PEM))
}

// Pub returns the materialized RSA public key corresponding
// to the padl-managed Key object
func (k *Key) Pub() (*rsa.PublicKey, error) {
	priv, err := keys.DecodePrivKeyPEM([]byte(k.PEM))
	if err != nil {
		return &priv.PublicKey, nil
	}
	return nil, fmt.Errorf("could not decode key PEM: %s", err)
}
