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

// Adds a Project to the user
func (u *User) AddProject(projectID string) {
	if !setContains(u.Projects, projectID) {
		u.Projects = append(u.Projects, projectID)
	}
}

//Removes project from the user
func (u *User) RemoveProject(projectID string) {
	for i, e := range u.Projects {
		if projectID == e {
			u.Projects[i] = u.Projects[len(u.Projects[i])-1]
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
