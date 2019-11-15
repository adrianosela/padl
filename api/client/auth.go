package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/aidoraide/wrdup/api/auth"
)

// Register registers a new user with email and a public PGP key
func (p *Padl) Register(email, password, pubKey string) error {
	pl := &payloads.RegistrationRequest{
		Email:    email,
		Password: password,
		PubKey:   pubKey,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return fmt.Errorf("could not marshall payload: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/register", p.HostURL),
		bytes.NewBuffer(plBytes))
	if err != nil {
		return fmt.Errorf("could not build http request: %s", err)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send http request: %s", err)
	}

	respByt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("could not read http response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", string(respByt))
	}

	return nil
}

// Login logs an existing user into a padl server
func (p *Padl) Login(email, password string) (string, error) {
	pl := &payloads.LoginRequest{
		Email:    email,
		Password: password,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return "", fmt.Errorf("could not marshall payload: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/login", p.HostURL),
		bytes.NewBuffer(plBytes))
	if err != nil {
		return "", fmt.Errorf("could not build http request: %s", err)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send http request: %s", err)
	}

	respByt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("could not read http response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non 200 status code received: %d", resp.StatusCode)
	}

	var lr payloads.LoginResponse
	if err := json.Unmarshal(respByt, &lr); err != nil {
		return "", fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return lr.Token, nil
}

// Valid checks whether a client has a valid token or not and returns the
// claims represented by the token body
func (p *Padl) Valid() (*auth.CustomClaims, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/valid", p.HostURL), nil)
	if err != nil {
		return nil, fmt.Errorf("could not build http request: %s", err)
	}
	p.setAuth(req)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send http request: %s", err)
	}
	respByt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not read http response body: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 status code received: %d", resp.StatusCode)
	}
	var cc auth.CustomClaims
	if err := json.Unmarshal(respByt, &cc); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &cc, nil
}
