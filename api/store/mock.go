package store

import "errors"

// MockDatabase is an in-memory database mock
type MockDatabase struct {
	users map[string]string
}

// NewMockDatabase is the constructor for MockDatabase
func NewMockDatabase() *MockDatabase {
	mdb := &MockDatabase{
		users: make(map[string]string),
	}
	return mdb
}

// CreateUser adds a new user to the database
func (db *MockDatabase) CreateUser(email, pub string) error {
	if _, ok := db.users[email]; ok {
		return errors.New("a padl account is already associated with that user")
	}
	db.users[email] = pub
	return nil
}
