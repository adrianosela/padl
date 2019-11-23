package project

import (
	"errors"

	"github.com/adrianosela/padl/api/privilege"
)

// Project represents a project in Padl
type Project struct {
	Name                string
	Description         string
	Members             map[string]privilege.Level
	ProjectKey          string
	PadlServiceAccounts map[string]string
}

// Summary TODO
type Summary struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewProject is the project object constructor
func NewProject(name, description, creator, projectKey string) *Project {
	return &Project{
		Name:        name,
		Description: description,
		ProjectKey:  projectKey,
		Members: map[string]privilege.Level{
			creator: privilege.PrivilegeLvlOwner,
		},
		PadlServiceAccounts: make(map[string]string),
	}
}

// AddUser adds a user to the project with the specified priv level
func (p *Project) AddUser(email string, priv privilege.Level) error {
	if p.HasUser(email) {
		return errors.New("user already in project")
	}
	p.Members[email] = priv
	return nil
}

// ChangeUserPrivilege changes a user's level of privilege on the project
func (p *Project) ChangeUserPrivilege(email string, priv privilege.Level) error {
	if _, ok := p.Members[email]; !ok {
		return errors.New("user not in project")
	}
	p.Members[email] = priv
	return nil
}

// HasUser checks whether a project has a user as a member
func (p *Project) HasUser(email string) bool {
	_, ok := p.Members[email]
	return ok
}

// RemoveUser removes a user from the project
func (p *Project) RemoveUser(email string) {
	if _, ok := p.Members[email]; ok {
		delete(p.Members, email)
	}
}

// SetPadlServiceAccount sets a padl service account for a project
func (p *Project) SetPadlServiceAccount(name string, tokenID string) error {
	if _, ok := p.PadlServiceAccounts[name]; ok {
		return errors.New("a service account with this name exists")
	}
	p.PadlServiceAccounts[name] = tokenID
	return nil
}

// RemovePadlServiceAccount removes a Padl service account
func (p *Project) RemovePadlServiceAccount(name string) {
	if _, ok := p.PadlServiceAccounts[name]; ok {
		delete(p.PadlServiceAccounts, name)
	}
}

// RotateProjectKey rotates the private key for a project
func (p *Project) RotateProjectKey() error {
	// TODO
	return nil
}
