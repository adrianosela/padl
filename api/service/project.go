package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/project"
)

func (s *Service) addProjectEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/newProject").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/addOwner").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/addEditor").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/addReader").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/removeOwner").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/removeEditor").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/removeReader").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/addSecret").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/removeSecret").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/createDeployKey").HandlerFunc(s.createNewProjectHandler)

}

// Change tules endpoint
// Add owner, editor, reader endpoint
// Remove owner, editor, reader endpoint
// Add secret ednpoint
// Remove secret endpoint
// Create deploy key
//

func (s *Service) createNewProjectHandler(w http.ResponseWriter, r *http.Request) {
	var proj *payloads.NewProjRequest
	if err := unmarshalRequestBody(r, &proj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	if err := proj.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new project: %s", err)))
		return
	}

	rules := project.Rules{
		RequireMFA:     proj.RequireMFA,
		RequireTeamKey: proj.RequireTeamKey,
	}
	project := project.NewProject(proj.Token, proj.Name, rules)

	if err := s.Database.CreateProject(project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new project: %s", err)))
		return
	}

	res := newProjectResponse{
		ID: project.ID,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

func (s *Service) addOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddOwnerRequest
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not add owner to project: %s", err)))
		return
	}
	// Find projects

	if p, err := s.Database.GetProject(req.ProjectID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as owner", req.Email, req.ProjectID)))
}

func (s *Service) removeOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var 
}