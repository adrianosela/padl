package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
)

// CreateKey creates a new padl managed private key
func (p *Padl) CreateKey(name, description string, bits int) (*kms.PrivateKey, error) {
	pl := &payloads.CreateKeyRequest{
		Name:        name,
		Description: description,
		Bits:        bits,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return nil, fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/key", p.HostURL),
		bytes.NewBuffer(plBytes))
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

	var pk kms.PrivateKey
	if err := json.Unmarshal(respByt, &pk); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &pk, nil
}

// GetPrivateKey gets a padl-managed private key from the server
// TODO: Remove this
func (p *Padl) GetPrivateKey(kid string) (*kms.PrivateKey, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/key/%s", p.HostURL, kid), nil)
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
	var pk kms.PrivateKey
	if err := json.Unmarshal(respByt, &pk); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &pk, nil
}

// GetPublicKey gets a padl-managed public key from the server
func (p *Padl) GetPublicKey(kid string) (*kms.PublicKey, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/key/public/%s", p.HostURL, kid), nil)
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
		return nil, fmt.Errorf("non 200 status code received: %d - %s", resp.StatusCode, string(respByt))
	}
	var pub kms.PublicKey
	if err := json.Unmarshal(respByt, &pub); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &pub, nil
}

// DeletePrivateKey deletes a padl-managed private key from the server
func (p *Padl) DeletePrivateKey(kid string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/key/%s", p.HostURL, kid), nil)
	if err != nil {
		return fmt.Errorf("could not build http request: %s", err)
	}
	p.setAuth(req)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send http request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non 200 status code received: %d", resp.StatusCode)
	}
	return nil
}

// AddUserToPrivateKey adds a user to the private key's access constrol list
func (p *Padl) AddUserToPrivateKey(kid, email string, priv int) (*kms.PrivateKey, error) {
	pl := &payloads.AddUserToKeyRequest{Email: email, PrivilegeLvl: priv}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return nil, fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/key/%s/user", p.HostURL, kid),
		bytes.NewBuffer(plBytes))
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
	var pk kms.PrivateKey
	if err := json.Unmarshal(respByt, &pk); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &pk, nil
}

// RemoveUserFromPrivateKey removes a user to the private key's access constrol list
func (p *Padl) RemoveUserFromPrivateKey(kid, email string) (*kms.PrivateKey, error) {
	pl := &payloads.RegistrationRequest{Email: email}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return nil, fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/key/%s/user", p.HostURL, kid),
		bytes.NewBuffer(plBytes))
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
	var pk kms.PrivateKey
	if err := json.Unmarshal(respByt, &pk); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &pk, nil
}

// DecryptSecret TODO
func (p *Padl) DecryptSecret(secret, kid string) (string, error) {
	pl := &payloads.DecryptSecretRequest{Secret: secret}

	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return "", fmt.Errorf("could not marshall payload: %s", err)
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/key/%s/decrypt", p.HostURL, kid),
		bytes.NewBuffer(plBytes))
	if err != nil {
		return "", fmt.Errorf("could not build http request: %s", err)
	}

	p.setAuth(req)

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

	var res payloads.DecryptSecretResponse
	if err := json.Unmarshal(respByt, &res); err != nil {
		return "", fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return res.Message, nil
}
