package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
)

// Padl represents a padl API client
type Padl struct {
	HostURL    string
	HTTPClient *http.Client
}

// NewPadlClient is the constructor for the Client object
func NewPadlClient(hostURL string, httpClient *http.Client) (*Padl, error) {
	if hostURL == "" {
		return nil, errors.New("host cannot be empty")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Padl{HostURL: hostURL, HTTPClient: httpClient}, nil
}

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

	resp, err := http.DefaultClient.Do(req)
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

	resp, err := http.DefaultClient.Do(req)
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
