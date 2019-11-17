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
	Projects   []string
}

// NewUser takes in user email, password, and public key id
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
		Projects:   []string{},
	}, nil
}

// CheckPassword verifies that a password matches the hash on the user
func (u *User) CheckPassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPass), []byte(pw))
}

// AddProject adds to the user
func (u *User) AddProject(name string) {
	if !setContains(u.Projects, name) {
		u.Projects = append(u.Projects, name)
	}
}

// RemoveProject removes a project from the user
func (u *User) RemoveProject(name string) {
	for i, e := range u.Projects {
		if name == e {
			u.Projects[i] = u.Projects[len(u.Projects)-1]
			u.Projects = u.Projects[:len(u.Projects)-1]
			return
		}
	}
}

func setContains(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}
