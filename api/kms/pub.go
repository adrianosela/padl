package kms

import (
	"crypto/rsa"
	"fmt"

	"github.com/adrianosela/padl/lib/keys"
)

// PublicKey represents a public key managed by padl
type PublicKey struct {
	ID  string `json:"id"`
	PEM string `json:"pem"`
}

// NewPublicKey is the constructor for the padl PublicKey object
func NewPublicKey(pubPEM string) (*PublicKey, error) {
	pub, err := keys.DecodePubKeyPEM([]byte(pubPEM))
	if err != nil {
		return nil, fmt.Errorf("could not decode public key PEM: %s", err)
	}
	return &PublicKey{
		ID:  keys.GetFingerprint(pub),
		PEM: pubPEM,
	}, nil
}

// PubRSA returns the materialized RSA public key corresponding
// to the padl-managed Key object
func (k *PublicKey) PubRSA() (*rsa.PublicKey, error) {
	pub, err := keys.DecodePubKeyPEM([]byte(k.PEM))
	if err != nil {
		return nil, fmt.Errorf("could not decode key PEM: %s", err)
	}
	return pub, nil
}
