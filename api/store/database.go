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
	ProjectExists(string) bool
	GetProject(string) (*project.Project, error)
	GetProjectByName(string) (*project.Project, error)
	UpdateProject(*project.Project) error
	ListProjects([]string) ([]*project.Project, error)
}
