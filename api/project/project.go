package project

type Project struct {
	ID         string            // project id
	Owners     []string          // usernames (emails)
	Editors   []string          // usernames (emails)
	Readers   []string          // usernames (emails)
	TeamKeys    map[string]string            // master key URI e.g. padl:{id}, aws:{arn}, gcp:{uri}
	DeployKeys map[string]string // read-only deploy keys - PEM enoded RSA keys
	Secrets    map[string]string          // secret ids 
	Settings   Rules             // mfa rules, etc
	Audit      string            // id of audit object for this project
}

type Rules struct {
	RequireOwnerMFA  bool
	RequireMemberMFA bool
	RequireTeamKey   bool
}


