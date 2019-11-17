package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/gorilla/mux"
)

func (s *Service) addKeyEndpoints() {
	// private key operations
	s.Router.Methods(http.MethodPost).Path("/key").Handler(s.Auth(s.createKeyHandler))
	s.Router.Methods(http.MethodGet).Path("/key/{kid}").Handler(s.Auth(s.getKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/key/{kid}").Handler(s.Auth(s.deleteKeyHandler))
	s.Router.Methods(http.MethodPost).Path("/key/{kid}/user").Handler(s.Auth(s.addUserToKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/key/{kid}/user").Handler(s.Auth(s.removeUserFromKeyHandler))
	// public key operations
	s.Router.Methods(http.MethodGet).Path("/key/public/{kid}").Handler(s.Auth(s.getPubKeyHandler))
	// decryption operations
	s.Router.Methods(http.MethodGet).Path("/decrypt/{kid}").Handler(s.Auth(s.decryptSecretHandler))
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
	key, err := kms.NewPrivateKey(newKeyPl.Bits, claims.Subject, newKeyPl.Name, newKeyPl.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not create new key: %s", err)))
		return
	}
	// save the new priv and pub keys
	if err := s.keystore.PutPrivKey(key); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not save new private key: %s", err)))
		return
	}
	pub, err := key.Pub()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not generate public half of key: %s", err)))
		return
	}
	if err = s.keystore.PutPubKey(pub); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not save new public key: %s", err)))
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
	claims := GetClaims(r)
	// get key id from request URL
	var id string
	if id = mux.Vars(r)["kid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no key id in request URL"))
		return
	}
	// get key from store
	key, err := s.keystore.GetPrivKey(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get key: %s", err)))
		return
	}
	// treat not having visibility of a key the same as the key not existing
	if key == nil || !key.IsVisibleTo(privilege.PrivilegeLvlReader, claims.Subject) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
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

func (s *Service) deleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get key id from request URL
	var id string
	if id = mux.Vars(r)["kid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no key id in request URL"))
		return
	}
	// get key from store
	key, err := s.keystore.GetPrivKey(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to get key: %s", err)))
		return
	}
	// treat not having visibility of a key the same as the key not existing
	if key == nil || !key.IsVisibleTo(privilege.PrivilegeLvlOwner, claims.Subject) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
		return
	}
	// delete key from store
	if err = s.keystore.DeletePrivKey(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to delete key: %s", err)))
		return
	}
	// send success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("deleted key %s successfully!", id)))
	return
}

func (s *Service) addUserToKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get key id from request URL
	var id string
	if id = mux.Vars(r)["kid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no key id in request URL"))
		return
	}
	// get payload
	var addUserPl *payloads.AddUserToKeyRequest
	if err := unmarshalRequestBody(r, &addUserPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := addUserPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate user addition request: %s", err)))
		return
	}
	// get key from store
	key, err := s.keystore.GetPrivKey(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to get key: %s", err)))
		return
	}
	// treat not having visibility of a key the same as the key not existing
	if key == nil || !key.IsVisibleTo(privilege.PrivilegeLvlOwner, claims.Subject) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
		return
	}
	// add user to key's users and save
	key.AddUser(addUserPl.Email, privilege.Level(addUserPl.PrivilegeLvl))
	if err = s.keystore.UpdatePrivKey(key); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to modify key: %s", err)))
		return
	}
	key.HideSecret()
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

func (s *Service) removeUserFromKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get key id from request URL
	var id string
	if id = mux.Vars(r)["kid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no key id in request URL"))
		return
	}
	// get payload
	var rmUserPl *payloads.RemoveUserFromKeyRequest
	if err := unmarshalRequestBody(r, &rmUserPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := rmUserPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate new key request: %s", err)))
		return
	}
	// get key from store
	key, err := s.keystore.GetPrivKey(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to get key: %s", err)))
		return
	}
	// treat not having visibility of a key the same as the key not existing
	if key == nil || !key.IsVisibleTo(privilege.PrivilegeLvlOwner, claims.Subject) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
		return
	}
	// rm user from key's users and save
	key.RemoveUser(rmUserPl.Email)
	if err = s.keystore.UpdatePrivKey(key); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error attempting to modify key: %s", err)))
		return
	}
	key.HideSecret()
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
	var decrpytPl *payloads.DecryptSecretRequest
	if err := unmarshalRequestBody(r, &decrpytPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := decrpytPl.Validate(); err != nil {
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

	// treat not having visibility of a key the same as the key not existing
	if key == nil || !key.IsVisibleTo(privilege.PrivilegeLvlReader, claims.Subject) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found"))
		return
	}
	// decode pem
	pkey, err := keys.DecodePrivKeyPEM([]byte(key.PEM))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("could not decode pem"))
		return
	}
	// decrypt secret
	message, err := keys.DecryptMessage([]byte(decrpytPl.Secret), pkey)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("could not decrypt secret"))
		return
	}

	res := payloads.DecryptSecretResponse{
		Message: string(message),
	}

	mbyt, err := json.Marshal(&res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could marshal response: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(mbyt)
	return
}
