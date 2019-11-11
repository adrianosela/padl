package payloads

import "errors"

type AddOwnerRequest struct {
	Email     string `json:"email"`
	ProjectID string `json:"projectId"`
}

type RemoveOwnerRequest struct {
	Email     string `json:"email"`
	ProjectID string `json:"projectId"`
}

type AddEditorRequest struct {
	ProjectID string `json:"projectId"`
}

type RemoveEditorRequest struct {
	ProjectID string `json:"projectId"`
}

type AddReaderRequest struct {
	ProjectID string `json:"projectId"`
}

type RemoveReaderRequest struct {
	ProjectID string `json:"projectId"`
}

type AddSecretRequest struct {
	ProjectID string `json:"projectId"`
}

type RemoveSecretRequest struct {
	ProjectID string `json:"projectId"`
}

type UpdateRulesRequest struct {
	ProjectID string `json:"projectId"`
}

type CreateDeployKey struct {
	ProjectID string `json:"projectId"`
}

type NewProjRequest struct {
	// Token could be sent in the header. For now sent as payload param
	Token           string `json:"token"`
	Name            string `json:"name"`
	CreateDeployKey bool   `json:"createDeployKey"`
	RequireMFA      bool   `json:"requireMFA"`
	RequireTeamKey  bool   `json:"requireTeamKey"`
}

// LoginRequest contains input for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegistrationRequest contains input for user registration
type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PubKey   string `json:"public_key"`
}

// Validate validates a login request payload
func (l *LoginRequest) Validate() error {
	if l.Email == "" {
		return errors.New("no email provided")
	}
	if l.Password == "" {
		return errors.New("no password provided")
	}
	return nil
}

// Validate validates a registration request payload
func (r *RegistrationRequest) Validate() error {
	if r.Email == "" {
		return errors.New("no email provided")
	}
	if r.Password == "" {
		return errors.New("no password provided")
	}
	// TODO: check PW complex enough
	if r.PubKey == "" {
		return errors.New("no public key provided")
	}
	// TODO: check pub key is valid RSA and at least 2048 bit
	return nil
}

func (p *NewProjRequest) Validate() error {
	if p.Name == "" {
		return errors.New("No project name defined")
	}

	if p.Token == "" {
		return errors.New("No token")
	}

	if !p.CreateDeployKey {
		return errors.New("Create deploy key rule not set")
	}

	if !p.RequireMFA {
		return errors.New("Require MFA rule not set")
	}

	if !p.RequireTeamKey {
		return errors.New("Require TeamKey rule not set")
	}
	return nil
}

func (a *AddOwnerRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	return nil
}
