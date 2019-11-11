package service

import (
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
)

func (s *Service) addAuthEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/register").HandlerFunc(s.registrationHandler)
}

func (s *Service) registrationHandler(w http.ResponseWriter, r *http.Request) {
	var regPl *payloads.RegistrationRequest
	if err := unmarshalRequestBody(r, &regPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := regPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// save new user in db
	if err := s.Database.CreateUser(regPl.Email, regPl.Password, regPl.PubKey); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new user: %s", err)))
		return
	}
	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successful registration of %s", regPl.Email)))
}
