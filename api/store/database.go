package store

import (
	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/api/user"
)

// Database represents all database operations
// for the padl API
type Database interface {
	PutUser(*user.User) error
	GetUser(string) (*user.User, error)
	UpdateUser(*user.User) error

	PutProject(*project.Project) error
	ProjectNameExists(string) bool
	GetProject(string) (*project.Project, error)
	UpdateProject(*project.Project) error
	ListProjects([]string) ([]*project.Project, []string, error)
}
