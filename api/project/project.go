package project

import (
	"errors"

	"github.com/google/uuid"
)

type Project struct {
	ID          string            // project id
	Name        string            // project name
	TeamMembers map[string]int    // maps a team memeber to a privalge level {0: reader, 1: editor, 2: owner}
	TeamKeys    map[string]string // master key URI e.g. padl:{id}, aws:{arn}, gcp:{uri}
	DeployKeys  []string          // read-only deploy keys - PEM enoded RSA keys
	Secrets     map[string]string // secret ids
	Settings    Rules             // mfa rules, etc
}

type Rules struct {
	RequireMFA     bool
	RequireTeamKey bool
}

func NewProject(email string, name string, rules Rules) *Project {
	// Create new project
	deployKeys := []string{}
	teamMembers := make(map[string]int)
	teamKeys := make(map[string]string)
	secrets := make(map[string]string)
	// Create a random ID
	id := uuid.Must(uuid.NewRandom()).String()

	var project = Project{
		ID:          id,
		Name:        name,
		TeamMembers: teamMembers,
		TeamKeys:    teamKeys,
		DeployKeys:  deployKeys,
		Secrets:     secrets,
		Settings:    rules,
	}
	project.TeamKeys[email] = project.generateTeamKey()
	project.TeamMembers[email] = Owner

	return &project
}

const Reader = 0
const Editor = 1
const Owner = 2

/*
NOTES
	remove keys when removing users

*/

func (p *Project) AddTeamMember(email string, userType int) (string, error) {
	if _, ok := p.TeamMembers[email]; ok {
		return "", errors.New("Team Member already exists")
	} else {
		p.TeamMembers[email] = userType
		p.TeamKeys[email] = p.generateTeamKey()
		return p.TeamKeys[email], nil

	}
}

func (p *Project) UpdateUserType(email string, userType int) error {
	if _, ok := p.TeamMembers[email]; ok {
		p.TeamMembers[email] = userType
		return nil
	} else {
		return errors.New("Team memeber does not exist")
	}
}

func (p *Project) GetUserType(email string) (int, error) {
	userType, ok := p.TeamMembers[email]
	if ok {
		return userType, nil
	} else {
		return -1, errors.New("Team memeber does not exist")
	}
}

func (p *Project) RemoveTeamMember(email string) error {
	if userType, ok := p.TeamMembers[email]; ok {
		//Check that they are not the only owner
		if userType == Owner {
			numOwners := 0
			for _, userType := range p.TeamMembers {
				if userType == Owner {
					numOwners++
				}
			}
			if numOwners > 1 {
				delete(p.TeamMembers, email)
				delete(p.TeamKeys, email)
				return nil
			} else {
				return errors.New("Only one owner exists, can not remove owner")
			}

		} else {
			delete(p.TeamMembers, email)
			delete(p.TeamKeys, email)
			return nil
		}
	}
	return nil
}

func (p *Project) CreateDeployKey() string {
	key := p.generateDeployKey()
	p.DeployKeys = append(p.DeployKeys, key)
	return key
}

func (p *Project) RemoveDeployKey(key string) error {
	i := p.indexOfSplice(key, p.DeployKeys)
	if i == -1 {
		return errors.New("Deploy key does not exist")
	}
	p.DeployKeys = append(p.DeployKeys[:i], p.DeployKeys[i+1:]...)
	return nil
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
func (p *Project) generateDeployKey() string {
	return "key"
}

func (p *Project) generateTeamKey() string {
	return "key"
}

func (p *Project) indexOfSplice(element string, list []string) int {
	for i, cur := range list {
		if cur == element {
			return i
		}
	}
	return -1
}
