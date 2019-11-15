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

// AddUser adds a user to a key with the specified privilege level.
// Note that this operation is stateless. No error is returned if the
// user is already part of the key's users.
func (k *Key) AddUser(email string, priv privilege.Level) {
	k.Users[email] = priv
}

// RemoveUser removes a user from the key
func (k *Key) RemoveUser(email string) error {
	if _, ok := k.Users[email]; !ok {
		return fmt.Errorf("user not in key")
	}
	delete(k.Users, email)
	return nil
}

// IsVisibleTo returns true if a given user (email) has the minimum
// required privilege level on a key, i.e. to check if a user has
// privilege to perform a subsequent action
func (k *Key) IsVisibleTo(required privilege.Level, email string) bool {
	if priv, ok := k.Users[email]; ok {
		return priv >= required
	}
	return false
}

// HideSecret simply changes the Key object such that
// the (secret) private key is no longer visible
func (k *Key) HideSecret() {
	k.PEM = "RSA PRIVATE KEY HIDDEN"
}
