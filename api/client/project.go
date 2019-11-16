package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/lib/padlfile"
)

// CreateProject creates a new project and returns a signed padlfile
func (p *Padl) CreateProject(name, description string) (*padlfile.File, error) {

	pl := &payloads.NewProjectRequest{
		Name:        name,
		Description: description,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return nil, fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/project", p.HostURL),
		bytes.NewBuffer(plBytes))
	if err != nil {
		return nil, fmt.Errorf("could not build http requests: %s", err)
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

	var r payloads.NewProjectResponse
	if err := json.Unmarshal(respByt, &r); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return r.Padfile, nil
}
