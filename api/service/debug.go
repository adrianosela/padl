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

	s.Router.Methods(http.MethodPost).Path("/getProject").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req *payloads.GetProjectRequest

			// get payload data
			if err := unmarshalRequestBody(r, &req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			p, err := s.Database.GetProject(req.ProjectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("no project found"))
				return
			}
			w.WriteHeader(http.StatusOK)
			byt, _ := json.Marshal(&p)
			w.Write(byt)
			return
		},
	)
}
