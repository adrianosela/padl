package project

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	// PrivilegeLvlReader TODO
	PrivilegeLvlReader = 0
	// PrivilegeLvlEditor TODO
	PrivilegeLvlEditor = 1
	// PrivilegeLvlOwner TODO
	PrivilegeLvlOwner = 2
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
func (p *Project) AddUser(email string, priv int) error {
	switch priv {
	case PrivilegeLvlOwner:
		p.Owners = append(p.Owners, email)
	case PrivilegeLvlEditor:
		p.Editors = append(p.Editors, email)
	case PrivilegeLvlReader:
		p.Readers = append(p.Readers, email)
	default:
		return fmt.Errorf("invalid privilege level %d", priv)
	}
	return nil
}

// RemoveUser removes a user from the project
func (p *Project) RemoveUser(email string) {
	p.Owners = removeFromSet(email, p.Owners)
	p.Editors = removeFromSet(email, p.Editors)
	p.Readers = removeFromSet(email, p.Readers)
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

func removeFromSet(elem string, slice []string) []string {
	for i, e := range slice {
		if e == elem {
			// move element to the back and pop it off
			slice[i] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}
