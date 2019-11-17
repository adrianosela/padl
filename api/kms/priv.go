package kms

import (
	"crypto/rsa"
	"fmt"

	"github.com/adrianosela/padl/lib/keys"
)

// PrivateKey represents a private key managed by padl
type PrivateKey struct {
	ID      string `json:"id"`
	Project string `json:"project"`
	PEM     string `json:"pem"`
}

// NewPrivateKey is the constructor for the padl PrivateKey object
func NewPrivateKey(bits int, project string) (*PrivateKey, error) {
	priv, pub, err := keys.GenerateRSAKeyPair(bits)
	if err != nil {
		return nil, fmt.Errorf("could not generate RSA key: %s", err)
	}
	return &PrivateKey{
		ID:      keys.GetFingerprint(pub),
		Project: project,
		PEM:     string(keys.EncodePrivKeyPEM(priv)),
	}, nil
}

// PrivRSA returns the materialized RSA private key corresponding
// to the padl-managed Key object
func (k *PrivateKey) PrivRSA() (*rsa.PrivateKey, error) {
	return keys.DecodePrivKeyPEM([]byte(k.PEM))
}

// PubRSA returns the materialized RSA public key corresponding
// to the padl-managed Key object
func (k *PrivateKey) PubRSA() (*rsa.PublicKey, error) {
	priv, err := keys.DecodePrivKeyPEM([]byte(k.PEM))
	if err != nil {
		return nil, fmt.Errorf("could not decode key PEM: %s", err)
	}
	return &priv.PublicKey, nil
}

// Pub returns the materialized padl-managed
// Public Key object corresponding to this privatekey
func (k *PrivateKey) Pub() (*PublicKey, error) {
	priv, err := keys.DecodePrivKeyPEM([]byte(k.PEM))
	if err != nil {
		return nil, fmt.Errorf("could not decode key PEM: %s", err)
	}
	return &PublicKey{
		ID:  k.ID,
		PEM: string(keys.EncodePubKeyPEM(&priv.PublicKey)),
	}, nil
}

// HideSecret simply changes the Key object such that
// the (secret) private key is no longer visible
func (k *PrivateKey) HideSecret() {
	k.PEM = "RSA PRIVATE KEY HIDDEN"
}
