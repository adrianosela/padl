package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/auth"
	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/gorilla/mux"
)

func (s *Service) addKeyEndpoints() {
	s.Router.Methods(http.MethodGet).Path("/key/{kid}").HandlerFunc(s.getPubKeyHandler) // note no auth
	s.Router.Methods(http.MethodPost).Path("/key/{kid}/decrypt").Handler(
		s.Auth(s.createServiceAccountHandler, []string{auth.ServiceAccountAudience, auth.PadlAPIAudience}...))
}

func (s *Service) getPubKeyHandler(w http.ResponseWriter, r *http.Request) {
	// get key id from request URL
	var id string
	if id = mux.Vars(r)["kid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no key id in request URL"))
		return
	}
	// get key from store, no need to check privs, pub keys are public
	pub, err := s.keystore.GetPubKey(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get key: %s", err)))
		return
	}
	if pub == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
		return
	}
	// return success
	pubByt, err := json.Marshal(&pub)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could marshal response: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(pubByt)
	return
}

func (s *Service) decryptSecretHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get key id from request URL
	var id string
	if id = mux.Vars(r)["kid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no key id in request URL"))
		return
	}
	// get payload
	var decryptPl *payloads.DecryptSecretRequest
	if err := unmarshalRequestBody(r, &decryptPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := decryptPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate decrypt secret request: %s", err)))
		return
	}
	// get key from store
	key, err := s.keystore.GetPrivKey(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to get key: %s", err)))
		return
	}
	// get owning project for the key
	p, err := s.database.GetProject(key.Project)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could get project: %s", err)))
		return
	}
	// treat not having visibility of a key the same as the key not existing
	if key == nil || !p.HasUser(claims.Subject) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
		return
	}
	// decode pem
	pkey, err := keys.DecodePrivKeyPEM([]byte(key.PEM))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not decode pem"))
		return
	}
	// decode secret
	raw, err := base64.StdEncoding.DecodeString(decryptPl.Secret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("secret is not base64 encoded"))
		return
	}
	// decrypt secret
	message, err := keys.DecryptMessage(raw, pkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not decrypt secret: %s", err)))
		return
	}
	// send success
	mbyt, err := json.Marshal(&payloads.DecryptSecretResponse{Message: base64.StdEncoding.EncodeToString(message)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could marshal response: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(mbyt)
	return
}
