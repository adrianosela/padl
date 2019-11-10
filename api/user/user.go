package user
 
type User struct {
	ID       string // user id
	Key      string // PEM encoded PGP key
	mfaToken string // e.g. duo mfa client to ping
}