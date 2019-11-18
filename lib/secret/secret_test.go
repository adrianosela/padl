package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePEM(t *testing.T) {
	// TODO
}

func TestDecodePEM(t *testing.T) {
	// TODO
}

func TestEncodeSimple(t *testing.T) {
	// TODO
}

func TestDecodeSimpleSecret(t *testing.T) {
	tests := []struct {
		testName     string
		toParse      string
		expectShards []*EncryptedShard
		expectErr    bool
		expectedErr  string
	}{
		{
			testName: "positive test",
			toParse:  "SOMEKEYID(SOMEVALUE)",
			expectShards: []*EncryptedShard{
				{
					KeyID: "SOMEKEYID",
					Value: "SOMEVALUE",
				},
			},
			expectErr: false,
		},
		{
			testName:    "negative test - bad format 1",
			toParse:     "SOMEKEYIDSOMEVALUE)",
			expectErr:   true,
			expectedErr: ErrMsgInvalidSimpleFmt,
		},
		{
			testName:    "negative test - bad format 2",
			toParse:     "",
			expectErr:   true,
			expectedErr: ErrMsgInvalidSimpleFmt,
		},
		{
			testName:    "negative test - bad format 3",
			toParse:     "SOMEKEY(",
			expectErr:   true,
			expectedErr: ErrMsgInvalidSimpleFmt,
		},
	}

	for _, test := range tests {
		sec, err := DecodeSimpleSecret(test.toParse)
		if test.expectErr {
			assert.Nil(t, sec, test.testName)
			assert.EqualError(t, err, test.expectedErr, test.testName)
		} else {
			assert.EqualValues(t, sec.Shards, test.expectShards, test.testName)
		}
	}
}

func TestEncodeDecodePEM(t *testing.T) {
	tests := []struct {
		testName      string
		testSecret    *Secret
		expectErr     bool
		expectedError string
	}{
		{
			testName: "positive test",
			testSecret: &Secret{
				Shards: []*EncryptedShard{
					{
						KeyID: "some key id",
						Value: "asdfghjkl",
					},
					{
						KeyID: "some key id",
						Value: "asdfghjkl",
					},
				},
			},
			expectErr: false,
		},
	}

	for _, test := range tests {
		enc, err := test.testSecret.EncodePEM()
		if test.expectErr {
			assert.EqualError(t, err, test.expectedError, test.testName)
			continue
		}
		dec, err := DecodePEM(enc)
		if test.expectErr {
			assert.EqualError(t, err, test.expectedError, test.testName)
			continue
		}
		assert.EqualValues(t, dec, test.testSecret, test.testName)
	}
}
