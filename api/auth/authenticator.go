package auth

import (
	"crypto/rsa"

	"github.com/adrianosela/padl/api/store"
)

const (
	defaultPadlIssuer = "padl.adrianosela.com"
	// PadlAPIAudience user audience for Padl API
	PadlAPIAudience = "api"
	// ServiceAccountAudience is the service account
	// audience for Padl API
	ServiceAccountAudience = "decrypt"
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
		a.aud = PadlAPIAudience
	}
	return a
}
