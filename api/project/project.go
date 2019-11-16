package project

import (
	"errors"

	"github.com/adrianosela/padl/api/privilege"
	"github.com/google/uuid"
)

// Project represents a project in Padl
type Project struct {
	ID           string
	Name         string
	Description  string
	Members      map[string]privilege.Level
	DeployKeys   map[string]string
	PadlfileHash string
}

// Summary TODO
type Summary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewProject is the project object constructor
func NewProject(name, description, creator string) *Project {
	return &Project{
		ID:          uuid.Must(uuid.NewRandom()).String(),
		Name:        name,
		Description: description,
		Members: map[string]privilege.Level{
			creator: privilege.PrivilegeLvlOwner,
		},
		DeployKeys: make(map[string]string),
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

// SetDeployKey sets a deploy key on a project
func (p *Project) SetDeployKey(name, value string) error {
	if _, ok := p.DeployKeys[name]; ok {
		return errors.New("a deploy key with this name exists")
	}
	p.DeployKeys[name] = value
	return nil
}

// RemoveDeployKey removes a deploy key from a project
func (p *Project) RemoveDeployKey(name string) {
	if _, ok := p.DeployKeys[name]; ok {
		delete(p.DeployKeys, name)
	}
}
