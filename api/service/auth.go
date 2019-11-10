package service

import (
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"

	"github.com/gorilla/mux"
)

func (c *Config) addAuthEndpoints(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/register").HandlerFunc(c.registrationHandler)
}

func (c *Config) registrationHandler(w http.ResponseWriter, r *http.Request) {
	// pick up email and public key from request body
	var reg *payloads.RegistrationRequest
	if err := unmarshalRequestBody(r, &reg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate email and pub key non-empty
	if reg.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no email provided"))
		return
	}
	if reg.PubKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no public key provided"))
		return
	}
	// save new user in db
	if err := c.Database.CreateUser(reg.Email, reg.PubKey); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new user: %s", err)))
		return
	}
	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successful registration of %s", reg.Email)))
}
