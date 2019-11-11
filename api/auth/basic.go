package auth

import (
	"fmt"
)

// Basic tests whether a pair of basic credentials are valid
func (a *Authenticator) Basic(uname, password string) error {
	usr, err := a.db.GetUser(uname)
	if err != nil {
		return fmt.Errorf("could not get user from db: %s", err)
	}
	// check passwords match
	if err := usr.CheckPassword(password); err != nil {
		return fmt.Errorf("could not verify password: %s", err)
	}
	return nil
}
