package project

import (
	"errors"

	"github.com/adrianosela/padl/api/privilege"
)

// Project represents a project in Padl
type Project struct {
	Name            string
	Description     string
	Members         map[string]privilege.Level
	ProjectKey      string
	ServiceAccounts map[string]string
}

// Summary is a name-description representation of a Project
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
		ServiceAccounts: make(map[string]string),
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

// HasServiceAccount checks whether a project has given service account
func (p *Project) HasServiceAccount(name string) bool {
	_, ok := p.ServiceAccounts[name]
	return ok
}

// RemoveUser removes a user from the project
func (p *Project) RemoveUser(email string) {
	if _, ok := p.Members[email]; ok {
		delete(p.Members, email)
	}
}

// SetServiceAccount sets a service account for a project
func (p *Project) SetServiceAccount(name string, keyID string) {
	p.ServiceAccounts[name] = keyID
}

// RemoveServiceAccount removes a service account
func (p *Project) RemoveServiceAccount(name string) {
	if _, ok := p.ServiceAccounts[name]; ok {
		delete(p.ServiceAccounts, name)
	}
}
