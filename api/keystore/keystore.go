package keystore

import (
	"errors"

	"github.com/adrianosela/padl/api/kms"
)

var (
	// ErrKeyExists is the error type returned
	// when the PutKey() method is called with
	// a key that already exists in the keystore
	ErrKeyExists = errors.New("key already exists")

	// ErrKeyNotFound is returned when GetKey(),
	// UpdateKey(), or DeleteKey() is called with
	// an id (or key with id in the case of update)
	// that is not found in the keystore
	ErrKeyNotFound = errors.New("key not found")
)

// Keystore represents all the necessary operations
// to build a key storage solution
type Keystore interface {
	PutPrivKey(*kms.PrivateKey) error
	GetPrivKey(string) (*kms.PrivateKey, error)
	UpdatePrivKey(*kms.PrivateKey) error
	DeletePrivKey(string) error

	PutPubKey(*kms.PublicKey) error
	GetPubKey(string) (*kms.PublicKey, error)
	DeletePubKey(string) error
}
