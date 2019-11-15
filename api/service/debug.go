package service

import (
	"encoding/json"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
)

func (s *Service) addDebugEndpoints() {
	// we conditionally add debug endpoints
	if !s.config.Debug {
		return
	}

	// healthcheck endpoint
	s.Router.Methods(http.MethodGet).Path("/healthcheck").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			byt, _ := json.Marshal(&payloads.HealthcheckResponse{
				Version:    s.config.Version,
				DeployTime: s.config.DeployTime.String(),
			})
			w.Write(byt)
			return
		},
	)
}
