package store

import "github.com/adrianosela/padl/api/project"

// Database represents all database operations
// for the padl API
type Database interface {
	// CreateUser takes in an email and
	// an armoured PGP public key
	CreateUser(string, string) error

	CreateProject(project *project.Project) error

	getProject(projectID string) error
}
