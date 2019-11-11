package store

import (
	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/api/user"
)

// Database represents all database operations
// for the padl API
type Database interface {
	// PutUser stores a new user in the db
	PutUser(*user.User) error
	// GetUser gets a user from the db by email
	GetUser(string) (*user.User, error)

	CreateProject(project *project.Project) error

	GetProject(projectID string) (*project.Project, error)
}
