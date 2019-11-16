package payloads

import (
	"errors"
)

// NewProjectRequest TODO
type NewProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

// Validate TODO
func (p *NewProjectRequest) Validate() error {
	if p.Name == "" {
		return errors.New("no project name provided")
	}
	if p.Description == "" {
		return errors.New("no project description provided")
	}
	return nil
}
