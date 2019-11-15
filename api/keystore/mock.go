package keystore

import (
	"github.com/adrianosela/padl/api/kms"
)

// MockKeystore is an in-memory implementation of
// the Keystore interface
type MockKeystore struct {
	privs map[string]*kms.PrivateKey
	pubs  map[string]*kms.PublicKey
}

// NewMockKeystore is the constructor for MockKeystore
func NewMockKeystore() *MockKeystore {
	return &MockKeystore{
		privs: make(map[string]*kms.PrivateKey),
		pubs:  make(map[string]*kms.PublicKey),
	}
}

// PutPrivKey adds a private key to the keystore
func (db *MockKeystore) PutPrivKey(k *kms.PrivateKey) error {
	if _, ok := db.privs[k.ID]; ok {
		return ErrKeyExists
	}
	db.privs[k.ID] = k
	return nil
}

// GetPrivKey gets a private key by id from the keystore
func (db *MockKeystore) GetPrivKey(id string) (*kms.PrivateKey, error) {
	if k, ok := db.privs[id]; ok {
		return k, nil
	}
	return nil, ErrKeyNotFound
}

// DeletePrivKey deletes a private key from the keystore
func (db *MockKeystore) DeletePrivKey(id string) error {
	if _, ok := db.privs[id]; !ok {
		return ErrKeyNotFound
	}
	delete(db.privs, id)
	return nil
}

// UpdatePrivKey updates a private key in the keystore
func (db *MockKeystore) UpdatePrivKey(k *kms.PrivateKey) error {
	if _, ok := db.privs[k.ID]; !ok {
		return ErrKeyNotFound
	}
	db.privs[k.ID] = k
	return nil
}

// PutPubKey adds a public key to the keystore
func (db *MockKeystore) PutPubKey(k *kms.PublicKey) error {
	if _, ok := db.pubs[k.ID]; ok {
		return ErrKeyExists
	}
	db.pubs[k.ID] = k
	return nil
}

// GetPubKey gets a public key by id from the keystore
func (db *MockKeystore) GetPubKey(id string) (*kms.PublicKey, error) {
	if k, ok := db.pubs[id]; ok {
		return k, nil
	}
	return nil, ErrKeyNotFound
}

// DeletePubKey deletes a public key by id from the keystore
func (db *MockKeystore) DeletePubKey(id string) error {
	if _, ok := db.pubs[id]; !ok {
		return ErrKeyNotFound
	}
	delete(db.pubs, id)
	return nil
}
