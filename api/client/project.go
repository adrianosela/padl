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
func (p *Padl) CreateProject(name, description string, bits int) (*padlfile.File, error) {
	pl := &payloads.NewProjectRequest{
		Name:        name,
		Description: description,
		KeyBits:     bits,
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
		return nil, fmt.Errorf("error: %s", string(respByt))
	}

	var pf padlfile.File
	if err := json.Unmarshal(respByt, &pf); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return &pf, nil
}

//CreateDeployKey Creates a new DeployKey
func (p *Padl) CreateDeployKey(projectName string, keyName string) (*payloads.CreateDeployKeyResponse, error) {
	pl := &payloads.CreateDeployKeyRequest{
		DeployKeyName: keyName,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return nil, fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/project/%s/deploy_key", p.HostURL, projectName),
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
		return nil, fmt.Errorf("error: %s", string(respByt))
	}

	var createKeyResp payloads.CreateDeployKeyResponse

	if err := json.Unmarshal(respByt, &createKeyResp); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return &createKeyResp, nil
}

//RemoveDeployKey Creates a new DeployKey
func (p *Padl) RemoveDeployKey(projectName string, keyName string) error {
	pl := &payloads.DeleteDeployKeyRequest{
		DeployKeyName: keyName,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/project/%s/deploy_key", p.HostURL, projectName),
		bytes.NewBuffer(plBytes))
	if err != nil {
		return fmt.Errorf("could not build http requests: %s", err)
	}
	p.setAuth(req)

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

// GetProject gets a project by name if the requesting user has access to it
func (p *Padl) GetProject(name string) (*project.Project, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/project/%s", p.HostURL, name), nil)
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
		return nil, fmt.Errorf("error: %s", string(respByt))
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
		return nil, fmt.Errorf("error: %s", string(respByt))
	}

	var listProjResp payloads.ListProjectsResponse
	if err := json.Unmarshal(respByt, &listProjResp); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return &listProjResp, nil
}

// AddUserToProject TODO
func (p *Padl) AddUserToProject(projectName string, email string, privilegeLvl int) (string, error) {
	pl := &payloads.AddUserToProjectRequest{
		Email:        email,
		PrivilegeLvl: privilegeLvl,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return "", fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/project/%s/user", p.HostURL, projectName),
		bytes.NewBuffer(plBytes))
	if err != nil {
		return "", fmt.Errorf("could not build http requests: %s", err)
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
	return string(respByt), nil
}

// RemoveUserFromProject TODO
func (p *Padl) RemoveUserFromProject(projectName string, email string) error {
	pl := &payloads.RemoveUserFromProjectRequest{
		Email: email,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/project/%s/user", p.HostURL, projectName),
		bytes.NewBuffer(plBytes))

	if err != nil {
		return fmt.Errorf("could not build http requests: %s", err)
	}

	p.setAuth(req)

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

// DeleteProject TODO
func (p *Padl) DeleteProject(projectName string) (string, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/project/%s", p.HostURL, projectName), nil)
	if err != nil {
		return "", fmt.Errorf("could not build http requests: %s", err)
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
	return string(respByt), nil
}
