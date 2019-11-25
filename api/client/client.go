package client

import (
	"errors"
	"fmt"
	"net/http"
)

// Padl represents a padl API client
type Padl struct {
	HostURL    string
	AuthToken  string
	HTTPClient *http.Client
}

// NewPadlClient is the constructor for the Client object
func NewPadlClient(hostURL, token string, httpClient *http.Client) (*Padl, error) {
	if hostURL == "" {
		return nil, errors.New("host cannot be empty")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Padl{
		HostURL:    hostURL,
		AuthToken:  token,
		HTTPClient: httpClient,
	}, nil
}

// setAuth sets the client's token in the authorization header of an http request
func (p *Padl) setAuth(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.AuthToken))
}
