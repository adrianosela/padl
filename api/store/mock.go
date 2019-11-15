package store

import (
	"errors"

	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/api/user"
)

// MockDatabase is an in-memory database mock
type MockDatabase struct {
	users    map[string]*user.User
	projects map[string]*project.Project
	keys     map[string]*kms.Key
}

// NewMockDatabase is the constructor for MockDatabase
func NewMockDatabase() *MockDatabase {
	mdb := &MockDatabase{
		users:    make(map[string]*user.User),
		projects: make(map[string]*project.Project),
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
	u, ok := db.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (db *MockDatabase) PutProject(p *project.Project) error {
	if _, ok := db.projects[p.ID]; ok {
		return errors.New("project already exists in the DB")
	}
	db.projects[p.ID] = p
	return nil
}

func (db *MockDatabase) GetProject(projectId string) (*project.Project, error) {
	if p, ok := db.projects[projectId]; ok {
		return p, nil
	}
	return nil, errors.New("project not found")
}

func (db *MockDatabase) UpdateProject(p *project.Project) error {
	if _, ok := db.projects[p.ID]; !ok {
		return errors.New("project not found")
	}
	db.projects[p.ID] = p
	return nil
}

func (db *MockDatabase) PutKey(k *kms.Key) error {
	if _, ok := db.keys[k.ID]; ok {
		return errors.New("key already in store")
	}
	db.keys[k.ID] = k
	return nil
}

func (db *MockDatabase) GetKey(id string) (*kms.Key, error) {
	if k, ok := db.keys[id]; ok {
		return k, nil
	}
	return nil, errors.New("key not found")
}

func (db *MockDatabase) UpdateKey(k *kms.Key) error {
	if _, ok := db.keys[k.ID]; !ok {
		return errors.New("key not found")
	}
	db.keys[k.ID] = k
	return nil
}
