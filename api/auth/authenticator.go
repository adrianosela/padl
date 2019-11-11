package auth

import (
	"crypto/rsa"

	"github.com/adrianosela/padl/api/store"
)

const (
	defaultPadlIssuer   = "api.padl.com"
	defaultPadlAudience = "api"
)

// Authenticator is the module in charge of authentication
type Authenticator struct {
	db     store.Database
	signer *rsa.PrivateKey
	iss    string
	aud    string
}

// NewAuthenticator is the Authenticator constructor
func NewAuthenticator(db store.Database, key *rsa.PrivateKey, iss, aud string) *Authenticator {
	a := &Authenticator{
		db:     db,
		signer: key,
		aud:    aud,
		iss:    iss,
	}
	if a.iss == "" {
		a.iss = defaultPadlIssuer
	}
	if a.aud == "" {
		a.aud = defaultPadlAudience
	}
	return a
}