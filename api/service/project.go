package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/project"
)

func (s *Service) addProjectEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/newProject").HandlerFunc(s.createProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/addOwner").HandlerFunc(s.addOwnerHandler)
	s.Router.Methods(http.MethodPost).Path("/addReader").HandlerFunc(s.addReaderHandler)
	s.Router.Methods(http.MethodPost).Path("/removeOwner").HandlerFunc(s.removeOwnerHandler)
	s.Router.Methods(http.MethodPost).Path("/removeReader").HandlerFunc(s.removeReaderHandler)

	s.Router.Methods(http.MethodPost).Path("/addEditor").HandlerFunc(s.addEditorHandler)
	s.Router.Methods(http.MethodPost).Path("/removeEditor").HandlerFunc(s.removeEditorHandler)

	// TODO

	s.Router.Methods(http.MethodPost).Path("/addSecret").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/removeSecret").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/createDeployKey").HandlerFunc(s.createNewProjectHandler)
	s.Router.Methods(http.MethodPost).Path("/updateRules").HandlerFunc(s.createNewProjectHandler)

}

// Change tules endpoint
// Add owner, editor, reader endpoint
// Remove owner, editor, reader endpoint
// Add secret ednpoint
// Remove secret endpoint
// Create deploy key
//

// [ PROJECT ]
func (s *Service) createProjectHandler(w http.ResponseWriter, r *http.Request) {
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
	// TODO: Get email from request context JWT
	project := project.NewProject(proj.Token, proj.Name, rules)

	if err := s.Database.PutProject(project); err != nil {
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

// [ OWNERS ]

func (s *Service) addOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddOwnerRequest
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// Find projects

	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to add an owner

	p.AddOwner(req.Email)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as owner", req.Email, req.ProjectID)))
}

func (s *Service) removeOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveOwnerRequest

	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}

	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to remove an owner

	p.RemoveOwner(req.Email)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed %s from project %s owners", req.Email, req.ProjectID)))

}

// [ READERS ]

func (s *Service) addReaderHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddReaderRequest
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// Find projects

	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to add an owner

	p.AddReader(req.Email)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as reader", req.Email, req.ProjectID)))
}

func (s *Service) removeReaderHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveOwnerRequest

	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}

	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to remove an owner

	p.RemoveReader(req.Email)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed %s from project %s readers", req.Email, req.ProjectID)))
}

// [ Editors ]

func (s *Service) addEditorHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddEditorRequest
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// Find projects

	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to add an owner

	p.AddEditor(req.Email)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as editor", req.Email, req.ProjectID)))
}

func (s *Service) removeEditorHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveEditorRequest

	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}

	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to remove an owner

	p.RemoveEditor(req.Email)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed %s from project %s editors", req.Email, req.ProjectID)))
}

// [ SECRETS ]

func (s *Service) addSecretHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddSecretRequest

	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
}
func (s *Service) removeSecretHandler(w http.ResponseWriter, r *http.Request) {

}
