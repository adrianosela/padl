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

// CreatePadlServiceAccount Creates a new padl service account
func (p *Padl) CreatePadlServiceAccount(projectName string, accountName string) (*payloads.CreatePadlServiceAccountResponse, error) {
	pl := &payloads.CreatePadlServiceAccountRequest{
		PadlServiceAccountName: accountName,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return nil, fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/project/%s/service_account", p.HostURL, projectName),
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

	var createKeyResp payloads.CreatePadlServiceAccountResponse

	if err := json.Unmarshal(respByt, &createKeyResp); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}

	return &createKeyResp, nil
}

// RemovePadlServiceAccount Removes padl service account
func (p *Padl) RemovePadlServiceAccount(projectName string, keyName string) error {
	pl := &payloads.DeletePadlServiceAccountRequest{
		PadlServiceAccountName: keyName,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/project/%s/service_account", p.HostURL, projectName),
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

// GetProjectKeys gets a project's list of keys by project name, the request will
// fail if the authorization token in the client is not for a user in the project
func (p *Padl) GetProjectKeys(name string) (*payloads.GetProjectKeysReponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/project/%s/keys", p.HostURL, name), nil)
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
	var keysResp payloads.GetProjectKeysReponse
	if err := json.Unmarshal(respByt, &keysResp); err != nil {
		return nil, fmt.Errorf("could not unmarshal http response body: %s", err)
	}
	return &keysResp, nil
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
func (p *Padl) AddUserToProject(projectName string, email string, privilegeLvl int) error {
	pl := &payloads.AddUserToProjectRequest{
		Email:        email,
		PrivilegeLvl: privilegeLvl,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
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

// RemoveUserFromProject TODO
func (p *Padl) RemoveUserFromProject(projectName string, email string) error {
	pl := &payloads.RemoveUserFromProjectRequest{
		Email: email,
	}
	plBytes, err := json.Marshal(&pl)
	if err != nil {
		return fmt.Errorf("could not marshall payload: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodDelete,
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
func (p *Padl) DeleteProject(projectName string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/project/%s", p.HostURL, projectName), nil)
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
