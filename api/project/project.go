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

func (p *Project) AddOwner(email string) {
	p.Owners = append(p.Owners, email)
}

func (p *Project) AddReader(email string) {
	p.Readers = append(p.Readers, email)
}

func (p *Project) AddEditor(email string) {
	p.Readers = append(p.Readers, email)
}

func (p *Project) AddDeployKey(key string) {
	id := uuid.Must(uuid.NewRandom()).String()
	p.DeployKeys[id] = key
}

func (p *Project) RemoveDeployKey(keyId string) {
	delete(p.DeployKeys, keyId)
}

func (p *Project) RemoveOwner(email string) {
	for i, cur_email := range p.Owners {
		if cur_email == email {
			p.Owners = append(p.Owners[:i], p.Owners[i+1:]...)
			break
		}
	}
}

func (p *Project) RemoveReader(email string) {
	for i, cur_email := range p.Readers {
		if cur_email == email {
			p.Readers = append(p.Readers[:i], p.Readers[i+1:]...)
			break
		}
	}
}

func (p *Project) RemoveEditor(email string) {
	for i, cur_email := range p.Editors {
		if cur_email == email {
			p.Readers = append(p.Editors[:i], p.Editors[i+1:]...)
			break
		}
	}
}

func (p *Project) AddSecret(secret string) string {
	id := uuid.Must(uuid.NewRandom()).String()
	p.Secrets[id] = secret
	return id
}

func (p *Project) RemoveSecret(secretId string) {
	delete(p.Secrets, secretId)
}

func (p *Project) RequireMFA(setting bool) {
	p.Settings.RequireMFA = setting
}

func (p *Project) RequireTeamKey(setting bool) {
	p.Settings.RequireTeamKey = setting
}

//TODO
func (p *Project) GenerateDeployKey() string {
	return "key"
}

// TODO: Use token to get user's email. For now just return email
//		 Maybe should be part of another package
func getEmail(token string) string {
	return "email"
}
