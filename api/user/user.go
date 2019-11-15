package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents a padl user
type User struct {
	Email      string
	HashedPass string
	KeyID      string
}

// NewUser takes in user email, password, and public key
// and returns a populated User with the hashed password
func NewUser(email, pass, keyID string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %s", err)
	}
	return &User{
		Email:      email,
		HashedPass: string(hash),
		KeyID:      keyID,
	}, nil
}

// CheckPassword verifies that a password matches the hash on the user
func (u *User) CheckPassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPass), []byte(pw))
}
