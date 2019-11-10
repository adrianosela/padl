package service

import (
	"github.com/adrianosela/padl/api/store"
	"github.com/gorilla/mux"
	"time"
)

// Config holds the service configuration
// necessary for endpoints to respond to requests
type Config struct {
	Version    string // server git hash
	DeployTime time.Time
	Database   store.Database
}

// NewPadlService returns an HTTP router multiplexer with
// attached handler functions
func NewPadlService(c Config) *mux.Router {
	router := mux.NewRouter()
	// TODO: add endpoints here
	return router
}
