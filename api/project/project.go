package project

import "github.com/google/uuid"

type Project struct {
	ID         string            // project id
	Name       string            // project name
	Owners     []string          // usernames (emails)
	Editors    []string          // usernames (emails)
	Readers    []string          // usernames (emails)
	TeamKeys   map[string]string // master key URI e.g. padl:{id}, aws:{arn}, gcp:{uri}
	DeployKeys map[string]string // read-only deploy keys - PEM enoded RSA keys
	Secrets    map[string]string // secret ids
	Settings   Rules             // mfa rules, etc
	Audit      string            // id of audit object for this project
}

type Rules struct {
	RequireMFA     bool
	RequireTeamKey bool
}

type ProjectConfig struct {
}

func NewProject(token string, name string, rules Rules) *Project {
	// Get email using JWT
	email := getEmail(token)
	// Create new project
	owners := []string{email}

	// Create a random ID

	id := uuid.Must(uuid.NewRandom()).String()
	var project = Project{
		ID:       id,
		Name:     name,
		Owners:   owners,
		Settings: rules,
	}

	return &project
}

// TODO: Use token to get user's email. For now just return email
//		 Maybe should be part of another package
func getEmail(token string) string {
	return "email"
}
