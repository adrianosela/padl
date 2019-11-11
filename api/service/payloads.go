package service

import (
	"errors"
)

type registrationRequest struct {
	Email  string `json:"email"`
	PubKey string `json:"pub"`
}

type newProjectRequest struct {
	// Token could be sent in the header. For now sent as payload param
	Token           string `json:"token"`
	Name            string `json:"name"`
	CreateDeployKey bool   `json:"createDeployKey"`
	RequireMFA      bool   `json:"requireMFA"`
	RequireTeamKey  bool   `json:"requireTeamKey"`
}

type newProjectResponse struct {
	ID string `json:"id"`
}

type addOwnerRequest struct {
	Email     string `json:"email"`
	ProjectID string `json:"projectId"`
}

func (p *newProjectRequest) validate() error {
	if p.Name == "" {
		return errors.New("No project name defined")
	}

	if p.Token == "" {
		return errors.New("No token")
	}

	if !p.CreateDeployKey {
		return errors.New("Create deploy key rule not set")
	}

	if !p.RequireMFA {
		return errors.New("Require MFA rule not set")
	}

	if !p.RequireTeamKey {
		return errors.New("Require TeamKey rule not set")
	}
	return nil
}
