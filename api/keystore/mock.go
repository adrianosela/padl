package keystore

import (
	"github.com/adrianosela/padl/api/kms"
)

// MockKeystore is an in-memory implementation of
// the Keystore interface
type MockKeystore struct {
	keys map[string]*kms.Key
}

// NewMockKeystore is the constructor for MockKeystore
func NewMockKeystore() *MockKeystore {
	return &MockKeystore{
		keys: make(map[string]*kms.Key),
	}
}

// PutKey adds a key to the keystore
func (db *MockKeystore) PutKey(k *kms.Key) error {
	if _, ok := db.keys[k.ID]; ok {
		return ErrKeyExists
	}
	db.keys[k.ID] = k
	return nil
}

// GetKey gets a key by id from the keystore
func (db *MockKeystore) GetKey(id string) (*kms.Key, error) {
	if k, ok := db.keys[id]; ok {
		return k, nil
	}
	return nil, ErrKeyNotFound
}

// DeleteKey deletes a key from the keystore
func (db *MockKeystore) DeleteKey(id string) error {
	if _, ok := db.keys[id]; !ok {
		return ErrKeyNotFound
	}
	delete(db.keys, id)
	return nil
}

// UpdateKey updates the key in the keystore
func (db *MockKeystore) UpdateKey(k *kms.Key) error {
	if _, ok := db.keys[k.ID]; !ok {
		return ErrKeyNotFound
	}
	db.keys[k.ID] = k
	return nil
}
