package kms

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/lib/keys"
)

// Key represents a key managed by padl
type Key struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Users map[string]int `json:"users"`
	PEM   string         `json:"pem"`
}

// NewKey is the constructor for the padl Key object
func NewKey(key *rsa.PrivateKey, creator, name string) (*Key, error) {
	if key == nil {
		return nil, errors.New("key cannot be nil")
	}
	return &Key{
		ID:   keys.GetFingerprint(&key.PublicKey),
		Name: name,
		PEM:  string(keys.EncodePrivKeyPEM(key)),
		Users: map[string]int{
			creator: project.PrivilegeLvlOwner,
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
