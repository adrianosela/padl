package kms

import (
	"crypto/rsa"
	"fmt"

	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/lib/keys"
)

// PrivateKey represents a private key managed by padl
type PrivateKey struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Users       map[string]privilege.Level `json:"users"`
	PEM         string                     `json:"pem"`
}

// NewPrivateKey is the constructor for the padl PrivateKey object
func NewPrivateKey(bits int, creator, name, descr string) (*PrivateKey, error) {
	priv, pub, err := keys.GenerateRSAKeyPair(bits)
	if err != nil {
		return nil, fmt.Errorf("could not generate RSA key: %s", err)
	}
	return &PrivateKey{
		ID:          keys.GetFingerprint(pub),
		Name:        name,
		Description: descr,
		PEM:         string(keys.EncodePrivKeyPEM(priv)),
		Users: map[string]privilege.Level{
			creator: privilege.PrivilegeLvlOwner,
		},
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

// AddUser adds a user to a key with the specified privilege level.
// Note that this operation is stateless. No error is returned if the
// user is already part of the key's users.
func (k *PrivateKey) AddUser(email string, priv privilege.Level) {
	k.Users[email] = priv
}

// RemoveUser removes a user from the key
func (k *PrivateKey) RemoveUser(email string) error {
	if _, ok := k.Users[email]; !ok {
		return fmt.Errorf("user not in key")
	}
	delete(k.Users, email)
	return nil
}

// IsVisibleTo returns true if a given user (email) has the minimum
// required privilege level on a key, i.e. to check if a user has
// privilege to perform a subsequent action
func (k *PrivateKey) IsVisibleTo(required privilege.Level, email string) bool {
	if priv, ok := k.Users[email]; ok {
		return priv >= required
	}
	return false
}

// HideSecret simply changes the Key object such that
// the (secret) private key is no longer visible
func (k *PrivateKey) HideSecret() {
	k.PEM = "RSA PRIVATE KEY HIDDEN"
}
