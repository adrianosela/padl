package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
)

// GetPublicKey gets a padl-managed public key from the server
func (p *Padl) GetPublicKey(kid string) (*kms.PublicKey, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/key/%s", p.HostURL, kid), nil)
	if err != nil {
		return nil, fmt.Errorf("could not build http request: %s", err)
	}
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
		return nil, fmt.Errorf("error: %s", string(respByt))
	}
	var pub kms.PublicKey
	if err := json.Unmarshal(respByt, &pub); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &pub, nil
}

// DecryptSecret decrypts and returns a secret in plaintext
func (p *Padl) DecryptSecret(secret, kid string) (string, error) {
	plBytes, err := json.Marshal(&payloads.DecryptSecretRequest{Secret: secret})
	if err != nil {
		return "", fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
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
		return "", fmt.Errorf("error: %s", string(respByt))
	}
	var res payloads.DecryptSecretResponse
	if err := json.Unmarshal(respByt, &res); err != nil {
		return "", fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(res.Message)
	if err != nil {
		return "", fmt.Errorf("could not decode decrypted message: %s", err)
	}
	return string(decoded), nil
}
