package service

import (
	"net/http"
)

func (s *Service) addKeyEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/key").Handler(s.Auth(s.createKeyHandler))
	s.Router.Methods(http.MethodGet).Path("/key/{kid}").Handler(s.Auth(s.getKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/key/{kid}").Handler(s.Auth(s.deleteKeyHandler))

	s.Router.Methods(http.MethodPost).Path("/key/{kid}/user").Handler(s.Auth(s.addUserToKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/key/{kid}/user").Handler(s.Auth(s.removeUserFromKeyHandler))
}

func (s *Service) createKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (s *Service) getKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (s *Service) deleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (s *Service) addUserToKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func (s *Service) removeUserFromKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}
