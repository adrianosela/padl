package service

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Config) addAuthEndpoints(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/register").HandlerFunc(c.registrationHandler)
}

func (c *Config) registrationHandler(w http.ResponseWriter, r *http.Request) {

}
