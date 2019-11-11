package store

import (
	"github.com/adrianosela/padl/api/user"
)

// Database represents all database operations
// for the padl API
type Database interface {
	// PutUser stores a new user in the db
	PutUser(*user.User) error
}
