package store

// Database represents all database operations
// for the padl API
type Database interface {
	// CreateUser takes in an email, password,
	// and an ASCII armoured RSA public key
	CreateUser(string, string, string) error
}
