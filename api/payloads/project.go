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

// GetProjectKeysReponse returns all the public
// key ids associated with a project
type GetProjectKeysReponse struct {
	Name       string
	MemberKeys []string
	ProjectKey string
	//	DeployKeys  []string (TODO)
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

// CreateServiceAccountRequest TODO
type CreateServiceAccountRequest struct {
	ServiceAccountName string `json:"name"`
}

// DeleteServiceAccountRequest TODO
type DeleteServiceAccountRequest struct {
	ServiceAccountName string `json:"serviceAccountName"`
}

// CreateServiceAccountResponse TODO
type CreateServiceAccountResponse struct {
	Token string `json:"token"`
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
func (r *CreateServiceAccountRequest) Validate() error {
	if r.ServiceAccountName == "" {
		return errors.New("No name provided")
	}
	return nil
}

// Validate TODO
func (r *DeleteServiceAccountRequest) Validate() error {
	if r.ServiceAccountName == "" {
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
		return errors.New("no space characters allowed for project name")
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
