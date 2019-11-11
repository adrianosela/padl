package service

import (
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/user"
)

func (s *Service) addAuthEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/register").HandlerFunc(s.registrationHandler)
	s.Router.Methods(http.MethodPost).Path("/login").HandlerFunc(s.loginHandler)
}

func (s *Service) registrationHandler(w http.ResponseWriter, r *http.Request) {
	var regPl *payloads.RegistrationRequest
	if err := unmarshalRequestBody(r, &regPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshal request body"))
		return
	}
	// validate payload
	if err := regPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// create new user object
	usr, err := user.NewUser(regPl.Email, regPl.Password, regPl.PubKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create user: %s", err)))
		return
	}
	// save new user in db
	if err := s.Database.PutUser(usr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new user: %s", err)))
		return
	}
	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successful registration of %s", regPl.Email)))
	return
}

func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginPl *payloads.LoginRequest
	if err := unmarshalRequestBody(r, &loginPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshal request body"))
		return
	}
	// validate payload
	if err := loginPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// get user object from db
	usr, err := s.Database.GetUser(loginPl.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid username or password")) // dont expose not-found
		return
	}
	// check passwords match
	if err := usr.CheckPassword(loginPl.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid username or password")) // dont expose bad password
	}
	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", "MOCK_TOKEN")))
	return
}
