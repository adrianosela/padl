package payloads

import (
	"errors"
	"strings"

	"github.com/adrianosela/padl/api/project"
)

// NewProjectRequest TODO
type NewProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	KeyBits     int    `json:"bits"`
}

// AddUserToProjectRequest TODO
type AddUserToProjectRequest struct {
	Email        string `json:"email"`
	PrivilegeLvl int    `json:"privilege"`
}

// RemoveUserFromProjectRequest TODO
type RemoveUserFromProjectRequest struct {
	Email string `json:"email"`
}

// CreateDeployKeyRequest TODO
type CreateDeployKeyRequest struct {
	DeployKeyName        string `json:"name"`
	DeployKeyDescription string `json:"description"`
}

// DeleteDeployKeyRequest TODO
type DeleteDeployKeyRequest struct {
	DeployKeyName string `json:"deployKey"`
}

// CreateDeployKeyResponse TODO
type CreateDeployKeyResponse struct {
	DeployKey string `json:"deployKey"`
}

// ListProjectsResponse TODO
type ListProjectsResponse struct {
	Projects []*project.Summary `json:"projects"`
}

// Validate TODO
func (a *AddUserToProjectRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	if a.PrivilegeLvl < 0 || a.PrivilegeLvl > 3 {
		return errors.New("invalid privilege level provided")
	}
	return nil
}

// Validate TODO
func (a *RemoveUserFromProjectRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	return nil
}

// Validate TODO
func (r *CreateDeployKeyRequest) Validate() error {
	if r.DeployKeyName == "" {
		return errors.New("No name provided")
	}
	if r.DeployKeyDescription == "" {
		return errors.New("No description provided")
	}
	return nil
}

// Validate TODO
func (r *DeleteDeployKeyRequest) Validate() error {
	if r.DeployKeyName == "" {
		return errors.New("No name provided")
	}
	return nil
}

// Validate validates a NewProjectRequest
func (p *NewProjectRequest) Validate() error {
	if p.Name == "" {
		return errors.New("no project name provided")
	}
	if strings.Contains(p.Name, " ") {
		return errors.New("no spaces characters allowed for project names")
	}
	if p.Description == "" {
		return errors.New("no project description provided")
	}
	if p.KeyBits != 512 &&
		p.KeyBits != 1024 &&
		p.KeyBits != 2048 &&
		p.KeyBits != 4096 {
		return errors.New("invalid bits, must be one of { 512, 1024, 2048, 4096 }")
	}
	return nil
}
