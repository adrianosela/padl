package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adrianosela/padl/api/project"

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

	var pf padlfile.File
	if err := json.Unmarshal(respByt, &pf); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return &pf, nil
}

// GetProject TODO
func (p *Padl) GetProject(pid string) (*project.Project, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/project/%s", p.HostURL, pid), nil)
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
	var project project.Project
	if err := json.Unmarshal(respByt, &project); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &project, nil
}

// GetProjectByName TODO
func (p *Padl) GetProjectByName(name string) (*project.Project, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/project/%s/%s", p.HostURL, "find", name), nil)
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
	var project project.Project
	if err := json.Unmarshal(respByt, &project); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &project, nil
}

// ListProjects TODO
func (p *Padl) ListProjects() (*payloads.ListProjectsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/projects", p.HostURL), nil)
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

	var listProjResp payloads.ListProjectsResponse
	if err := json.Unmarshal(respByt, &listProjResp); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return &listProjResp, nil
}
