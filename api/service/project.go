package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/project"
	"github.com/gorilla/mux"
)

func (s *Service) addProjectEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/project").Handler(s.Auth(s.createProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/project/{pid}").Handler(s.Auth(s.getProjectHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{pid}").Handler(s.Auth(s.deleteProjectHandler))

	s.Router.Methods(http.MethodPost).Path("/project/{pid}/user").Handler(s.Auth(s.addUserHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{pid}/user").Handler(s.Auth(s.removeUserHandler))

	s.Router.Methods(http.MethodPost).Path("/project/{pid}/deploy_key").Handler(s.Auth(s.createDeployKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{pid}/deploy_key").Handler(s.Auth(s.removeDeployKeyHandler))
}

func (s *Service) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get payload
	var projPl *payloads.NewProjectRequest
	if err := unmarshalRequestBody(r, &projPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload
	if err := projPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate new project request: %s", err)))
		return
	}
	// create project object and save it
	project := project.NewProject(projPl.Name, claims.Subject)
	if err := s.Database.PutProject(project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not save new project: %s", err)))
		return
	}
	// marshal response
	byt, err := json.Marshal(&payloads.NewProjectResponse{ID: project.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal project json: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}

func (s *Service) getProjectHandler(w http.ResponseWriter, r *http.Request) {
	var id string
	if id = mux.Vars(r)["pid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project ID in request URL"))
		return
	}
	p, err := s.Database.GetProject(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get project: %s", err)))
		return
	}
	/*
	 TODO: check user is in project or else return 401
	*/
	byt, err := json.Marshal(&p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal project json: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}

func (s *Service) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (s *Service) addUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var id string
	if id = mux.Vars(r)["pid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project ID in request URL"))
		return
	}
	// read request body
	var addUserPl *payloads.AddUserToProjectRequest
	if err := unmarshalRequestBody(r, &addUserPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload data
	if err := addUserPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}
	/*
	   TODO: check user has privs for project or else return 403
	*/
	fmt.Println(claims.Subject) // REMOVE
	if err = p.AddUser(addUserPl.Email, addUserPl.PrivilegeLvl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could add user to project: %s", err)))
		return
	}
	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("user %s added to project %s successfully!", addUserPl.Email, p.Name)))
	return
}

func (s *Service) removeUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var id string
	if id = mux.Vars(r)["pid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project ID in request URL"))
		return
	}
	// read request body
	var rmUserPl *payloads.RemoveUserFromProjectRequest
	if err := unmarshalRequestBody(r, &rmUserPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload data
	if err := rmUserPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}
	/*
	   TODO: check user has privs for project or else return 403
	*/
	if rmUserPl.Email == claims.Subject {
		if p.HasUser(claims.Subject, project.PrivilegeLvlOwner) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("owners cannot remove themselves from projects"))
			return
		}
	}
	// update project
	p.RemoveUser(rmUserPl.Email)
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("user %s removed from project %s successfully!", rmUserPl.Email, p.Name)))
	return
}

func (s *Service) createDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var id string
	if id = mux.Vars(r)["pid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project ID in request URL"))
		return
	}
	// get payload data
	var dkeyPl *payloads.CreateDeployKeyRequest
	if err := unmarshalRequestBody(r, &dkeyPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
		return
	}
	// validate payload data
	if err := dkeyPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}
	/*
	 TODO: check user is in project or else return 401
	*/
	fmt.Println(claims.Subject) // REMOVE
	newDeployKey := "FIXME:MOCKDEPLOYKEY"
	if err = p.SetDeployKey("mock", newDeployKey); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not set new deploy key: %s", err)))
		return
	}
	// update project
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	// marshall response
	byt, err := json.Marshal(&payloads.CreateDeployKeyResponse{
		DeployKey: newDeployKey,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal project json: %s", err)))
		return
	}
	// send resp
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}

func (s *Service) removeDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var id string
	if id = mux.Vars(r)["pid"]; id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project ID in request URL"))
		return
	}
	// get payload data
	var deleteKeyPl payloads.DeleteDeployKeyRequest
	if err := unmarshalRequestBody(r, &deleteKeyPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not unmarshall request body: %s", err)))
		return
	}
	// validate payload data
	if err := deleteKeyPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate payload: %s", err)))
		return
	}
	// fetch project from database
	p, err := s.Database.GetProject(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}
	/*
	 TODO: check user is in project or else return 401
	*/
	fmt.Println(claims.Subject) // REMOVE
	// update project
	p.RemoveDeployKey(deleteKeyPl.DeployKeyName)
	if err := s.Database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed deploy key from project %s", p.Name)))
	return
}
