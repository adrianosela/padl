package service

import (
	"encoding/json"
	"net/http"

	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/store"
	"github.com/gorilla/mux"
)

// Service holds the service configuration
// necessary for endpoints to respond to requests
type Service struct {
	Config   *config.Config
	Router   *mux.Router
	Database store.Database
}

// NewPadlService returns an HTTP router multiplexer with
// attached handler functions
func NewPadlService(c *config.Config) *Service {
	svc := &Service{
		Config:   c,
		Router:   mux.NewRouter(),
		Database: store.NewMockDatabase(),
	}

	// healthcheck endpoint
	svc.Router.Methods(http.MethodGet).Path("/healthcheck").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			byt, _ := json.Marshal(&payloads.HealthcheckResponse{
				Version:    c.Version,
				DeployTime: c.DeployTime.String(),
			})
			w.Write(byt)
			return
		},
	)

	svc.addAuthEndpoints()

	// TODO: add endpoints here

	return svc
}
