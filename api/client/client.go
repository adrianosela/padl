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
func (p *Padl) Register(email, pubKey string) error {
	pl := &payloads.RegistrationRequest{
		Email:  email,
		PubKey: pubKey,
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
