package secret

import (
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

const (
	pemBlockType       = "PADL ENCRYPTED SECRET"
	simpleFmtSeparator = "-"

	// ErrMsgInvalidSimpleFmt is returned when trying to decode
	// a simple-format secret that is not simple-encoded
	ErrMsgInvalidSimpleFmt = "bad format"

	// ErrMsgCouldNotDecodePEM is returned when trying to decode
	// a PEM-format secret that is not PEM encoded
	ErrMsgCouldNotDecodePEM = "could not decode pem block"
)

// Secret represents an encrypted secret
type Secret struct {
	Shards []*EncryptedShard
}

// EncodePEM returns an encrypted secret in a PEM block
func (s *Secret) EncodePEM() (string, error) {
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  pemBlockType,
		Bytes: []byte(s.EncodeSimple()),
	})
	return string(pemBytes), nil
}

// DecodePEM returns an encrypted secret from a pem block
func DecodePEM(s string) (*Secret, error) {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, errors.New(ErrMsgCouldNotDecodePEM)
	}
	sec, err := DecodeSimpleSecret(string(block.Bytes))
	if err != nil {
		return nil, err
	}
	return sec, nil
}

// EncodeSimple returns a simple string representation of the encrypted secret.
// This format is KEY_ID(VALUE)
func (s *Secret) EncodeSimple() string {
	ret := ""
	for i, sh := range s.Shards {
		ret = strings.Join([]string{ret, fmt.Sprintf("%s(%s)", sh.KeyID, sh.Value)}, "")
		if i != len(s.Shards)-1 {
			ret = strings.Join([]string{ret, simpleFmtSeparator}, "")
		}
	}
	return ret
}

// DecodeSimpleSecret returns a sharded representation of the encrypted secret
func DecodeSimpleSecret(s string) (*Secret, error) {
	parts := strings.Split(s, simpleFmtSeparator)
	sec := &Secret{Shards: []*EncryptedShard{}}
	for _, p := range parts {
		es, err := decodeSimpleShard(p)
		if err != nil {
			return nil, err
		}
		sec.Shards = append(sec.Shards, es)
	}
	return sec, nil
}

func decodeSimpleShard(simpleEncodedShard string) (*EncryptedShard, error) {
	p1 := strings.Split(simpleEncodedShard, "(")
	if len(p1) < 2 {
		return nil, errors.New(ErrMsgInvalidSimpleFmt)
	}
	p2 := strings.Split(p1[1], ")")
	if len(p2) < 2 {
		return nil, errors.New(ErrMsgInvalidSimpleFmt)
	}
	return &EncryptedShard{
		KeyID: p1[0],
		Value: p2[0],
	}, nil
}
