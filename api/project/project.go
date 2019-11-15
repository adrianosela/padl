package project

import (
	"errors"
	"fmt"

	"github.com/adrianosela/padl/api/privilege"
	"github.com/google/uuid"
)

// Project represents a project in Padl
type Project struct {
	ID         string
	Name       string
	Owners     []string
	Editors    []string
	Readers    []string
	DeployKeys map[string]string
}

// NewProject is the project object constructor
func NewProject(name, creator string) *Project {
	return &Project{
		ID:         uuid.Must(uuid.NewRandom()).String(),
		Name:       name,
		Owners:     []string{creator},
		Editors:    []string{},
		Readers:    []string{},
		DeployKeys: make(map[string]string),
	}
}

// AddUser adds a user to the project with the specified priv level
func (p *Project) AddUser(email string, priv privilege.Level) error {
	switch priv {
	case privilege.PrivilegeLvlOwner:
		p.Owners = addToSet(p.Owners, email)
	case privilege.PrivilegeLvlEditor:
		p.Editors = addToSet(p.Editors, email)
	case privilege.PrivilegeLvlReader:
		p.Readers = addToSet(p.Readers, email)
	default:
		return fmt.Errorf("invalid privilege level %d", priv)
	}
	return nil
}

// HasUser checks whether a project has a user with a given priv level
func (p *Project) HasUser(email string, priv privilege.Level) bool {
	switch priv {
	case privilege.PrivilegeLvlOwner:
		return setContains(p.Owners, email)
	case privilege.PrivilegeLvlEditor:
		return setContains(p.Editors, email)
	case privilege.PrivilegeLvlReader:
		return setContains(p.Readers, email)
	}
	return false
}

// RemoveUser removes a user from the project
func (p *Project) RemoveUser(email string) {
	p.Owners = removeFromSet(p.Owners, email)
	p.Editors = removeFromSet(p.Editors, email)
	p.Readers = removeFromSet(p.Readers, email)
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

func setContains(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

func addToSet(slice []string, elem string) []string {
	if setContains(slice, elem) {
		return slice
	}
	return append(slice, elem)
}

func removeFromSet(slice []string, elem string) []string {
	for i, e := range slice {
		if e == elem {
			// move element to the back and pop it off
			slice[i] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}
