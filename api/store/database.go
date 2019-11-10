package store

// Database represents all database operations
// for the padl API
type Database interface {
	// CreateUser takes in an email and
	// an armoured PGP public key
	CreateUser(string, string) error
}
