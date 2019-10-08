package padl

type Project struct {
	ID         string            // project id
	Owners     []string          // usernames (emails)
	Members    []string          // usernames (emails)
	TeamKey    string            // master key URI e.g. padl:{id}, aws:{arn}, gcp:{uri}
	DeployKeys map[string]string // read-only deploy keys - PEM enoded RSA keys
	Secrets    []string          // secret ids
	Settings   Rules             // mfa rules, etc
	Audit      string            // id of audit object for this project
}

type Rules struct {
	RequireOwnerMFA  bool
	RequireMemberMFA bool
	RequireTeamKey   bool
}

type User struct {
	ID       string // user id
	Key      string // PEM encoded PGP key
	mfaToken string // e.g. duo mfa client to ping
}
