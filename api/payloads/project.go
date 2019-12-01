package payloads

import (
	"errors"
	"strings"

	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/lib/keys"
)

// NewProjectRequest is the expected payload for the
// of the project creation endpooint
type NewProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	KeyBits     int    `json:"bits"`
}

// GetProjectKeysReponse returns all the public
// key ids associated with a project
type GetProjectKeysReponse struct {
	Name       string   `json:"name"`
	MemberKeys []string `json:"member_keys"`
	ProjectKey string   `json:"project_key"`
	DeployKeys []string `json:"deploy_keys"`
}

// AddUserToProjectRequest is the expected payload
// for the user addition to project endpoint
type AddUserToProjectRequest struct {
	Email        string `json:"email"`
	PrivilegeLvl int    `json:"privilege"`
}

// RemoveUserFromProjectRequest is the expected payload
// for the user removal from project endpoint
type RemoveUserFromProjectRequest struct {
	Email string `json:"email"`
}

// CreateServiceAccountRequest is the expected payload
// for the service-account creation endpoint
type CreateServiceAccountRequest struct {
	ServiceAccountName string `json:"name"`
	PubKey             string `json:"public_key"`
}

// DeleteServiceAccountRequest is the expected payload
// for the service-account deletion endpoint
type DeleteServiceAccountRequest struct {
	ServiceAccountName string `json:"serviceAccountName"`
}

// CreateServiceAccountResponse is the response
// type of the service-account creation endpoint
type CreateServiceAccountResponse struct {
	Token string `json:"token"`
}

// ListProjectsResponse is the response of the project list endpoint
type ListProjectsResponse struct {
	Projects []*project.Summary `json:"projects"`
}

// Validate validates a user addition to project request
func (a *AddUserToProjectRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	if a.PrivilegeLvl < 0 || a.PrivilegeLvl > 3 {
		return errors.New("invalid privilege level provided")
	}
	return nil
}

// Validate validates a user removal from project request
func (a *RemoveUserFromProjectRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	return nil
}

// Validate validates a service-account creation request
func (r *CreateServiceAccountRequest) Validate() error {
	if r.ServiceAccountName == "" {
		return errors.New("No name provided")
	}
	if strings.Contains(r.ServiceAccountName, " ") {
		return errors.New("no space characters allowed for service account name")
	}
	if strings.Contains(r.ServiceAccountName, ".") {
		return errors.New("no . characters allowed for service account name")
	}
	if r.PubKey == "" {
		return errors.New("no public key provided")
	}
	if _, err := keys.DecodePubKeyPEM([]byte(r.PubKey)); err != nil {
		return errors.New("provided key is not a PEM encoded RSA public key")
	}
	return nil
}

// Validate validates a service-account deletion request
func (r *DeleteServiceAccountRequest) Validate() error {
	if r.ServiceAccountName == "" {
		return errors.New("No name provided")
	}
	return nil
}

// Validate validates a project creation request
func (p *NewProjectRequest) Validate() error {
	if p.Name == "" {
		return errors.New("no project name provided")
	}
	if strings.Contains(p.Name, " ") {
		return errors.New("no space characters allowed for project name")
	}
	if strings.Contains(p.Name, ".") {
		return errors.New("no . characters allowed for project name")
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
