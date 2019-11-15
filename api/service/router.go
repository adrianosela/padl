package service

import (
	"log"

	"github.com/adrianosela/padl/api/auth"
	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/keystore"
	"github.com/adrianosela/padl/api/store"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/gorilla/mux"
)

// Service holds the service configuration
// necessary for endpoints to respond to requests
type Service struct {
	Router *mux.Router

	config        *config.Config
	database      store.Database
	keystore      keystore.Keystore
	authenticator *auth.Authenticator
}

// NewPadlService returns an HTTP router multiplexer with
// attached handler functions
func NewPadlService(c *config.Config) *Service {
	// initialize mongodb
	// db, err := store.NewMongoDB(
	// 	c.Database.ConnectionString,
	// 	c.Database.Name,
	// 	c.Database.UsersCollectionName,
	// )
	// if err != nil {
	// 	log.Fatalf("could not initialize mongodb: %s", err)
	// }
	db := store.NewMockDatabase()
	ks := keystore.NewMockKeystore()

	priv, err := keys.DecodePrivKeyPEM([]byte(c.Auth.SigningKey))
	if err != nil {
		log.Fatalf("could not materialize jwt signing key: %s", err)
	}

	svc := &Service{
		Router:        mux.NewRouter(),
		config:        c,
		database:      db,
		keystore:      ks,
		authenticator: auth.NewAuthenticator(db, priv, "api.padl.com", "api"),
	}

	svc.addDebugEndpoints()
	svc.addAuthEndpoints()
	svc.addProjectEndpoints()
	svc.addKeyEndpoints()

	return svc
}
