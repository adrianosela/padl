package secret

import (
	"crypto/rsa"
)

// Secret represents an object containing a set of encrypted secrets
type Secret struct {
	KeyID string `json:"kid"`
	// TODO
}

// Decrypt decrypts all secrets in a Secret object with a given private key
// and returns all of the VAR_NAME=value pairs corresponding to the secrets
func (s *Secret) Decrypt(k *rsa.PrivateKey) ([]string, error) {
	// TODO
	return []string{"MOCK_PW=mockpassword", "MOCK_PW_2=mockpassword2"}, nil
}
