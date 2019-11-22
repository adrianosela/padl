package secret

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"

	"fmt"

	"github.com/adrianosela/padl/lib/keys"
)

const (
	// ErrMsgEmptyValue is returned when the user attempts
	// to create a shard with no data provided
	ErrMsgEmptyValue = "shard can not have empty value"

	// ErrMsgCouldNotDecrypt is returned when an error occurs
	// while attempting to encrypt a shard
	ErrMsgCouldNotEncrypt = "could not encrypt shard value"

	// ErrMsgCouldNotDecrypt is returned when an error occurs
	// while attempting to decrypt a shard
	ErrMsgCouldNotDecrypt = "could not decrypt shard value"

	// ErrMsgIncorrectDecryptionKey is returned when the user attempts to
	// decrypt an EncryptedShard with the wrong key (key id mismatch)
	ErrMsgIncorrectDecryptionKey = "the provided key does not match the shard's encryption key's fingerprint"

	// ErrMsgCouldNotDecode is returned when a shard value could not
	// be base64 encoded
	ErrMsgCouldNotDecode = "could not b64 decode shard value"
)

// Shard describes a piece of secret that has been split
// with Shamir's Secret Sharing Algorithm
type Shard struct {
	Value []byte
}

// EncryptedShard represents a shard that has been encrypted
type EncryptedShard struct {
	Value string `json:"value"`
	KeyID string `json:"key_id"`
}

// NewShard returns a populated Shard struct
func NewShard(value []byte) (*Shard, error) {
	if len(value) == 0 {
		return nil, errors.New(ErrMsgEmptyValue)
	}
	return &Shard{
		Value: value,
	}, nil
}

// Encrypt encrypts and ASCII armours a shard's value
func (s *Shard) Encrypt(k *rsa.PublicKey) (*EncryptedShard, error) {
	if len(s.Value) == 0 {
		return nil, errors.New(ErrMsgEmptyValue)
	}
	armoured, err := encryptAndArmourShamirPart(s.Value, k)
	if err != nil {
		return nil, err
	}
	return &EncryptedShard{
		Value: armoured,
		KeyID: keys.GetFingerprint(k),
	}, nil
}

func (es *EncryptedShard) Decrypt(k *rsa.PrivateKey) (*Shard, error) {
	fp := keys.GetFingerprint(&k.PublicKey)
	if es.KeyID != fp {
		return nil, errors.New(ErrMsgIncorrectDecryptionKey)
	}
	val, err := decryptAndUnarmourShamirPart(es.Value, k)
	if err != nil {
		return nil, err
	}
	return &Shard{Value: val}, nil
}

func decryptAndUnarmourShamirPart(data string, k *rsa.PrivateKey) ([]byte, error) {
	// remove ASCII armour from piece
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", ErrMsgCouldNotDecode, err)
	}
	// decrypt the raw encrypted message
	dec, err := keys.DecryptMessage(raw, k)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", ErrMsgCouldNotDecrypt, err)
	}
	return dec, nil
}

func encryptAndArmourShamirPart(data []byte, k *rsa.PublicKey) (string, error) {
	// encrypt shard value
	enc, err := keys.EncryptMessage(data, k)
	if err != nil {
		return "", fmt.Errorf("%s: %s", ErrMsgCouldNotEncrypt, err)
	}
	// ASCII armour the encrypted shard
	return base64.StdEncoding.EncodeToString(enc), nil
}
