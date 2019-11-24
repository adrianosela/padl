package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/auth"
	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/user"
)

func (s *Service) addAuthEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/register").HandlerFunc(s.registrationHandler)
	s.Router.Methods(http.MethodPost).Path("/login").HandlerFunc(s.loginHandler)
	s.Router.Methods(http.MethodPost).Path("/rotate").Handler(s.Auth(s.rotateKeyHandler))
	s.Router.Methods(http.MethodGet).Path("/valid").Handler(s.Auth(s.validHandler))
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
	// create padl pub key object and store it publicly
	pub, err := kms.NewPublicKey(regPl.PubKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := s.keystore.PutPubKey(pub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not store user's public key: %s", err)))
		return
	}
	// create new user object
	usr, err := user.NewUser(regPl.Email, regPl.Password, pub.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create user: %s", err)))
		return
	}
	// save new user in db
	if err := s.database.PutUser(usr); err != nil {
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
	if err := s.authenticator.Basic(loginPl.Email, loginPl.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid username or password")) // do not expose reason
		return
	}

	user, err := s.database.GetUser(loginPl.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to get user from the database: %s", err)))
		return
	}

	token, err := s.authenticator.GenerateJWT(user.Email, auth.PadlAPIAudience)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // fixme: if this happens we want to know
		return
	}

	lr := &payloads.LoginResponse{Token: token}
	byt, err := json.Marshal(&lr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // fixme: if this happens we want to know
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}

func (s *Service) rotateKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// unmarshal payload
	var rotatePl *payloads.RotateKeyRequest
	if err := unmarshalRequestBody(r, &rotatePl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshal request body"))
		return
	}
	// validate payload
	if err := rotatePl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// create padl pub key object and store it publicly
	pub, err := kms.NewPublicKey(rotatePl.PubKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := s.keystore.PutPubKey(pub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not store user's public key: %s", err)))
		return
	}
	// update user in db
	user, err := s.database.GetUser(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to get user from the database: %s", err)))
		return
	}
	user.KeyID = pub.ID
	if err := s.database.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to update user in the database: %s", err)))
		return
	}
	// send success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("rotated user's key successfully!"))
	return
}

func (s *Service) validHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	byt, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}
