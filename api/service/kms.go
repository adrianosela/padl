package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
)

func (s *Service) addKeyEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/key").Handler(s.Auth(s.createKeyHandler))
	s.Router.Methods(http.MethodGet).Path("/key/{kid}").Handler(s.Auth(s.getKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/key/{kid}").Handler(s.Auth(s.deleteKeyHandler))

	s.Router.Methods(http.MethodPost).Path("/key/{kid}/user").Handler(s.Auth(s.addUserToKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/key/{kid}/user").Handler(s.Auth(s.removeUserFromKeyHandler))
}

func (s *Service) createKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get payload
	var newKeyPl *payloads.CreateKeyRequest
	if err := unmarshalRequestBody(r, &newKeyPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := newKeyPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate new key request: %s", err)))
		return
	}
	// create key object and save it
	key, err := kms.NewKey(newKeyPl.Bits, claims.Subject, newKeyPl.Name, newKeyPl.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not create new key: %s", err)))
		return
	}
	// save the new key
	if err := s.Database.PutKey(key); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not save new project: %s", err)))
		return
	}
	// return success
	keybyt, err := json.Marshal(&key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could marshal response: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(keybyt)
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
