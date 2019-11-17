package store

import (
	"errors"

	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/api/user"
)

// MockDatabase is an in-memory database mock
type MockDatabase struct {
	users    map[string]*user.User
	projects map[string]*project.Project
}

// NewMockDatabase is the constructor for MockDatabase
func NewMockDatabase() *MockDatabase {
	mdb := &MockDatabase{
		users:    make(map[string]*user.User),
		projects: make(map[string]*project.Project),
	}
	return mdb
}

// UpdateUser updates a user in the database
func (db *MockDatabase) UpdateUser(usr *user.User) error {
	if _, ok := db.users[usr.Email]; !ok {
		return errors.New("User doesn't exist")
	}
	db.users[usr.Email] = usr
	return nil
}

// PutUser adds a new user to the database
func (db *MockDatabase) PutUser(usr *user.User) error {
	if _, ok := db.users[usr.Email]; ok {
		return errors.New("a padl account is already associated with that email")
	}
	db.users[usr.Email] = usr
	return nil
}

// GetUser gets a user from the database
func (db *MockDatabase) GetUser(email string) (*user.User, error) {
	u, ok := db.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

// PutProject puts a project in the database
func (db *MockDatabase) PutProject(p *project.Project) error {
	if _, ok := db.projects[p.Name]; ok {
		return errors.New("project already exists in the DB")
	}
	db.projects[p.Name] = p
	return nil
}

// GetProject gets a project from the database
func (db *MockDatabase) GetProject(name string) (*project.Project, error) {
	if p, ok := db.projects[name]; ok {
		return p, nil
	}
	return nil, errors.New("project not found")
}

// ProjectNameExists returns true if a name exists in the
// padl global namespace for projects
func (db *MockDatabase) ProjectNameExists(name string) bool {
	if _, ok := db.projects[name]; ok {
		return true
	}
	return false
}

// UpdateProject updates a project in the database
func (db *MockDatabase) UpdateProject(p *project.Project) error {
	if _, ok := db.projects[p.Name]; !ok {
		return errors.New("project not found")
	}
	db.projects[p.Name] = p
	return nil
}

// ListProjects returns a list of requested (by name) projects
func (db *MockDatabase) ListProjects(names []string) ([]*project.Project, []string, error) {
	prjs := []*project.Project{}
	notFound := []string{}
	for _, n := range names {
		if p, ok := db.projects[n]; ok {
			prjs = append(prjs, p)
		} else {
			notFound = append(notFound, n)
		}
	}
	return prjs, notFound, nil
}
