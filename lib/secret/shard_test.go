package secret

import (
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/adrianosela/padl/lib/keys"
	"github.com/stretchr/testify/assert"
)

func TestNewShard(t *testing.T) {
	tests := []struct {
		testName    string
		shardValue  []byte
		shardsTotal int
		expectErr   bool
		expectedErr string
	}{
		{
			testName:    "positive test",
			shardValue:  []byte{0x80, 0x80, 0x80, 0x80},
			shardsTotal: 5,
			expectErr:   false,
		},
		{
			testName:    "empty value test",
			shardValue:  []byte{},
			shardsTotal: 5,
			expectErr:   true,
			expectedErr: ErrMsgEmptyValue,
		},
	}
	for _, test := range tests {
		s, err := NewShard(test.shardValue)
		if test.expectErr {
			assert.Nil(t, s, test.testName)
			assert.EqualError(t, err, test.expectedErr, test.testName)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, s.Value, test.shardValue, test.testName)
		}
	}
}

func TestEncrypt(t *testing.T) {

	badKey := &rsa.PublicKey{}

	_, goodKey, err := keys.GenerateRSAKeyPair(2048)
	if err != nil {
		assert.FailNow(t, "could not generate test key")
	}

	tests := []struct {
		testName    string
		shard       *Shard
		key         *rsa.PublicKey
		expectErr   bool
		expectedErr string
	}{
		{
			testName: "positive test",
			shard: &Shard{
				Value: []byte{0x80, 0x80, 0x80, 0x80},
			},
			key:       goodKey,
			expectErr: false,
		},
		{
			testName: "empty value test",
			shard: &Shard{
				Value: []byte{},
			},
			key:         goodKey,
			expectErr:   true,
			expectedErr: ErrMsgEmptyValue,
		},
		{
			testName: "bad key test",
			shard: &Shard{
				Value: []byte{0x80, 0x80, 0x80, 0x80},
			},
			key:         badKey,
			expectErr:   true,
			expectedErr: fmt.Sprintf("%s: %s: %s", ErrMsgCouldNotEncrypt, "crypto/rsa", "missing public modulus"),
		},
	}

	for _, test := range tests {
		es, err := test.shard.Encrypt(test.key)
		if test.expectErr {
			assert.Nil(t, es, test.testName)
			assert.EqualError(t, err, test.expectedErr, test.testName)
		} else {
			assert.Nil(t, err)
			assert.NotEqual(t, es.Value, test.shard.Value, test.testName)
			assert.Equal(t, es.KeyID, keys.GetFingerprint(test.key), test.testName)
		}
	}
}

func TestDecrypt(t *testing.T) {
	// a good key pair to encrypt/decrypt successfully
	goodPriv, goodPub, err := keys.GenerateRSAKeyPair(2048)
	if err != nil {
		assert.FailNow(t, "could not generate test key")
	}

	// a different key to test attempting to decrypt with the wrong key
	differentPriv, _, err := keys.GenerateRSAKeyPair(2048)
	if err != nil {
		assert.FailNow(t, "could not generate different test key")
	}

	// an invalid key to attempt decrypting with
	badPriv := &rsa.PrivateKey{}

	mockSecret := "this is a secret"

	goodShard := &Shard{
		Value: []byte(mockSecret),
	}

	goodEncryptedShard, err := goodShard.Encrypt(goodPub)
	if err != nil {
		assert.FailNow(t, "could not encrypt mock shard")
	}

	tests := []struct {
		testName    string
		encShard    *EncryptedShard
		key         *rsa.PrivateKey
		expectErr   bool
		expectedErr string
	}{
		{
			testName:  "positive test",
			encShard:  goodEncryptedShard,
			key:       goodPriv,
			expectErr: false,
		},
		{
			testName:    "incorerct key test",
			encShard:    goodEncryptedShard,
			key:         differentPriv,
			expectErr:   true,
			expectedErr: ErrMsgIncorrectDecryptionKey,
		},
		{
			testName:    "bad key test",
			encShard:    goodEncryptedShard,
			key:         badPriv,
			expectErr:   true,
			expectedErr: ErrMsgIncorrectDecryptionKey,
		},
		{
			testName: "bad encoding test",
			encShard: &EncryptedShard{
				Value: "thisisnotbase64",
				KeyID: keys.GetFingerprint(goodPub),
			},
			key:         goodPriv,
			expectErr:   true,
			expectedErr: fmt.Sprintf("%s: %s", ErrMsgCouldNotDecode, "illegal base64 data at input byte 12"),
		},
		{
			testName: "not encrypted test",
			encShard: &EncryptedShard{
				Value: "dGhpc2lzYmFzZTY0", // "thisisbase64" in b64
				KeyID: keys.GetFingerprint(goodPub),
			},
			key:         goodPriv,
			expectErr:   true,
			expectedErr: fmt.Sprintf("%s: %s: %s", ErrMsgCouldNotDecrypt, "crypto/rsa", "decryption error"),
		},
	}

	for _, test := range tests {
		s, err := test.encShard.Decrypt(test.key)
		if test.expectErr {
			assert.Nil(t, s, test.testName)
			assert.EqualError(t, err, test.expectedErr, test.testName)
		} else {
			assert.Nil(t, err, test.testName)
			assert.NotEqual(t, s.Value, test.encShard.Value, test.testName)
			assert.Equal(t, string(s.Value), mockSecret, test.testName)
		}
	}
}
