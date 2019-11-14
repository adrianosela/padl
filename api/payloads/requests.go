package payloads

import "errors"

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

// NewProjectRequest TODO
type NewProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RegistrationRequest contains input for user registration
type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PubKey   string `json:"public_key"`
}

// LoginRequest contains input for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
