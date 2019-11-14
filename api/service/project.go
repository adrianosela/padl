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
	// [ Owners ]
	s.Router.Methods(http.MethodPost).Path("/addOwner").Handler(s.Auth(s.addOwnerHandler))
	s.Router.Methods(http.MethodPost).Path("/removeOwner").Handler(s.Auth(s.removeOwnerHandler))
	// [ Readers ]
	s.Router.Methods(http.MethodPost).Path("/addReader").Handler(s.Auth(s.addReaderHandler))
	s.Router.Methods(http.MethodPost).Path("/removeReader").Handler(s.Auth(s.removeReaderHandler))
	// [ Editor ]
	s.Router.Methods(http.MethodPost).Path("/addEditor").Handler(s.Auth(s.addEditorHandler))
	s.Router.Methods(http.MethodPost).Path("/removeEditor").Handler(s.Auth(s.removeEditorHandler))
	// [ Secret ]
	s.Router.Methods(http.MethodPost).Path("/addSecret").Handler(s.Auth(s.addSecretHandler))
	s.Router.Methods(http.MethodPost).Path("/removeSecret").Handler(s.Auth(s.removeSecretHandler))
	// [ Deploy Keys ]
	s.Router.Methods(http.MethodPost).Path("/createDeployKey").Handler(s.Auth(s.createDeployKeyHandler))
	s.Router.Methods(http.MethodPost).Path("/removeDeployKey").Handler(s.Auth(s.removeDeployKeyHandler))

	// TODO

	// s.Router.Methods(http.MethodPost).Path("/updateRules").HandlerFunc(s.createNewProjectHandler)
	// s.Router.Methods(http.MethodPost).Path("/removeProject").HandlerFunc(s.createNewProjectHandler)
	// s.Router.Methods(http.MethodPost).Path("/updatePriviledge").HandlerFunc(s.createNewProjectHandler)

}

const Reader = 0
const Editor = 1
const Owner = 2

// [ PROJECT ]
func (s *Service) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var proj *payloads.NewProjRequest
	// get claims
	claims := GetClaims(r)
	_, err := json.Marshal(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal claims"))
		return
	}
	// get request data
	if err := unmarshalRequestBody(r, &proj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate request data
	if err := proj.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new project: %s", err)))
		return
	}
	// create project rules
	rules := project.Rules{
		RequireMFA:     proj.RequireMFA,
		RequireTeamKey: proj.RequireTeamKey,
	}
	// create project
	project := project.NewProject(claims.Subject, proj.Name, rules)
	// save oproject
	if err := s.Database.PutProject(project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new project: %s", err)))
		return
	}

	res := payloads.NewProjResponse{
		ID: project.ID,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

// [ OWNERS ]

func (s *Service) addOwnerHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if client user has permission to remove an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// add owner to project
	key, err := p.AddTeamMember(req.Email, Owner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to add owner: %s", err)))
		return
	}

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	res := payloads.AddOwnerResponse{
		TeamKey: key,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

func (s *Service) removeOwnerHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to remove owner: %s", err)))
		return
	}

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed owner %s from project %s", req.Email, req.ProjectID)))

}

// [ READERS ]

func (s *Service) addReaderHandler(w http.ResponseWriter, r *http.Request) {
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

	// find project
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to add a reader
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can add readers: %s", err)))
		return
	}

	// add reader to project
	key, err := p.AddTeamMember(req.Email, Reader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to add reader: %s", err)))
		return
	}

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	res := payloads.AddReaderResponse{
		TeamKey: key,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

func (s *Service) removeReaderHandler(w http.ResponseWriter, r *http.Request) {
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
	// check if client user has permission to remove readers
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove readers: %s", err)))
		return
	}

	// Remove reader from project
	err = p.RemoveTeamMember(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to remove reader: %s", err)))
		return
	}

	// Update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed reader %s from project %s", req.Email, req.ProjectID)))
}

// [ Editors ]

func (s *Service) addEditorHandler(w http.ResponseWriter, r *http.Request) {
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

	// fetch project from database
	p, err := s.Database.GetProject(req.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	// check if client user has permission to add an editor
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can add editors: %s", err)))
		return
	}

	// add editor to project
	key, err := p.AddTeamMember(req.Email, Editor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to add editor: %s", err)))
		return
	}

	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	res := payloads.AddEditorResponse{
		TeamKey: key,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

func (s *Service) removeEditorHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Owner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove editors: %s", err)))
		return
	}

	// remove editor from project
	err = p.RemoveTeamMember(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to remove editor: %s", err)))
		return
	}

	// update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed editor %s from project %s", req.Email, req.ProjectID)))
}

// [ SECRETS ]

func (s *Service) addSecretHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("readers cannot add secrets: %s", err)))
		return
	}

	// add secret to project
	id := p.AddSecret(req.Secret)

	// update projects
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}

	res := payloads.AddSecretResponse{
		ID: id,
	}

	bytesJSON, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(bytesJSON))
}

func (s *Service) removeSecretHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if client user has permission to remove secrets
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("readers cannot remove secrets: %s", err)))
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

func (s *Service) createDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if client user has permission to add deploy keys
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err)))
		return
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("readers cannot create deploy keys: %s", err)))
		return
	}

	// create deploy key
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

func (s *Service) removeDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if client user has permission to remove deploy keys
	t, err := p.GetUserType(claims.Subject)
	if err != nil {
	}
	if t < Editor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("readers cannot remove deploy keys: %s", err)))
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
