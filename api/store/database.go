package store

import (
	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/api/user"
)

// Database represents all database operations
// for the padl API
type Database interface {
	PutUser(*user.User) error
	GetUser(string) (*user.User, error)

	PutProject(*project.Project) error
	GetProject(string) (*project.Project, error)
	UpdateProject(*project.Project) error

	PutKey(*kms.Key) error
	GetKey(string) (*kms.Key, error)
	UpdateKey(*kms.Key) error
}
