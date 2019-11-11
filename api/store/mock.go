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
		users: make(map[string]*user.User),
	}
	return mdb
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
	if u, ok := db.users[email]; ok {
		return u, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (db *MockDatabase) PutProject(p *project.Project) error {
	if _, ok := db.projects[p.ID]; ok {
		return errors.New("Project already exists in the DB")
	}
	db.projects[p.ID] = p
	return nil
}

func (db *MockDatabase) GetProject(projectId string) (*project.Project, error) {

	if p, ok := db.projects[projectId]; ok {
		return p, nil
	} else {
		return nil, errors.New("project not found")
	}
}

func (db *MockDatabase) UpdateProject(p *project.Project) error {

	oldP, ok := db.projects[p.ID]
	if !ok {
		return errors.New("project not found")
	}

	if oldP.ID != p.ID {
		delete(db.projects, oldP.ID)
	}

	db.projects[p.ID] = p

	return nil
}
