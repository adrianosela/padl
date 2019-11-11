package service

import (
	"encoding/json"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
)

func (s *Service) addDebugEndpoints() {
	// we conditionally add debug endpoints
	if !s.Config.Debug {
		return
	}

	// healthcheck endpoint
	s.Router.Methods(http.MethodGet).Path("/healthcheck").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			byt, _ := json.Marshal(&payloads.HealthcheckResponse{
				Version:    s.Config.Version,
				DeployTime: s.Config.DeployTime.String(),
			})
			w.Write(byt)
			return
		},
	)
}
