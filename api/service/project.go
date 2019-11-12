package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/project"
)

func (s *Service) addProjectEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/newProject").Handler(s.Auth(s.createProjectHandler))

	s.Router.Methods(http.MethodPost).Path("/addOwner").HandlerFunc(s.addOwnerHandler)
	s.Router.Methods(http.MethodPost).Path("/removeOwner").HandlerFunc(s.removeOwnerHandler)

	s.Router.Methods(http.MethodPost).Path("/addReader").HandlerFunc(s.addReaderHandler)
	s.Router.Methods(http.MethodPost).Path("/removeReader").HandlerFunc(s.removeReaderHandler)

	s.Router.Methods(http.MethodPost).Path("/addEditor").HandlerFunc(s.addEditorHandler)
	s.Router.Methods(http.MethodPost).Path("/removeEditor").HandlerFunc(s.removeEditorHandler)

	s.Router.Methods(http.MethodPost).Path("/addSecret").HandlerFunc(s.addSecretHandler)
	s.Router.Methods(http.MethodPost).Path("/removeSecret").HandlerFunc(s.removeSecretHandler)

	s.Router.Methods(http.MethodPost).Path("/createDeployKey").HandlerFunc(s.createDeployKeyHandler)

	// TODO

	s.Router.Methods(http.MethodPost).Path("/removeDeployKey").HandlerFunc(s.createDeployKeyHandler)
	// s.Router.Methods(http.MethodPost).Path("/updateRules").HandlerFunc(s.createNewProjectHandler)
	// s.Router.Methods(http.MethodPost).Path("/removeProject").HandlerFunc(s.createNewProjectHandler)
	// s.Router.Methods(http.MethodPost).Path("/updatePriviledge").HandlerFunc(s.createNewProjectHandler)

}

const Reader = 0
const Editor = 1
const Owner = 2

// Change tules endpoint
// Add owner, editor, reader endpoint
// Remove owner, editor, reader endpoint
// Add secret ednpoint
// Remove secret endpoint
// Create deploy key
// Remove project
// Change priviledge level
// [ PROJECT ]
func (s Service) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var proj *payloads.NewProjRequest
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	if err := unmarshalRequestBody(r, &proj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
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
	project := project.NewProject(claims.Subject, proj.Name, rules)

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

func (s Service) addOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddOwnerRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to remove an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// add owner to project
	p.AddTeamMember(req.Email, Owner)

	// Update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as owner", req.Email, req.ProjectID)))
}

func (s Service) removeOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveOwnerRequest
	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to remove an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to remove owner: %s", err)))
		return
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// remove owner from project
	err = p.RemoveTeamMember(req.Email)
	if err != nil {

	}

	// Update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed %s from project %s owners", req.Email, req.ProjectID)))

}

// [ READERS ]

func (s Service) addReaderHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddReaderRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}

	// validate paylaod data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}

	// Find project
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to remove an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// Add reader to project
	p.AddTeamMember(req.Email, Reader)

	// Update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as reader", req.Email, req.ProjectID)))
}

func (s Service) removeReaderHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveReaderRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}
	// Check if client user has permission to remove an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		//TODO
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// Remove reader from project
	p.RemoveTeamMember(req.Email)

	// Update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed %s from project %s readers", req.Email, req.ProjectID)))
}

// [ Editors ]

func (s Service) addEditorHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddEditorRequest
	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}

	// Fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to add an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// Add editor to project
	p.AddTeamMember(req.Email, Editor)

	// Update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added %s to project %s as editor", req.Email, req.ProjectID)))
}

func (s Service) removeEditorHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveEditorRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// Check if client user has permission to remove an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		// TODO
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// remove editor from project
	p.RemoveTeamMember(req.Email)

	// update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed %s from project %s editors", req.Email, req.ProjectID)))
}

// [ SECRETS ]

func (s Service) addSecretHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.AddSecretRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to add secrets
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		// TODO
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("reader cannot add secrets: %s", err)))
		return
	}

	// add secret to project
	p.AddSecret(req.Secret)

	// update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully added secret to project %s", req.ProjectID)))
}

func (s Service) removeSecretHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveSecretRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to add secrets
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("reader cannot remove secrets: %s", err)))
		return
	}

	// remove secret from project
	p.RemoveSecret(req.Secret)

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed secret from project %s", req.ProjectID)))
}

// [ DEPLOY KEYS ]

func (s Service) createDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.CreateDeployKeyRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to add secrets
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		// TODO
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("reader cannot create deploy keys: %s", err)))
		return
	}

	// remove secret from project
	k := p.CreateDeployKey()

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	res := payloads.CreateDeployKeyResponse{
		DeployKey: k,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

func (s Service) removeDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req *payloads.RemoveDeployKeyRequest

	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}

	// get payload data
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	// validate payload data
	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to add secrets
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("reader cannot remove deploy keys: %s", err)))
		return
	}

	if err := p.RemoveDeployKey(req.DeployKey); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to remove key: %s", err)))
		return
	}

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed deploy key from project %s", req.ProjectID)))
}
