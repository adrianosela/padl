package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adrianosela/padl/api/auth"
	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/lib/padlfile"
	"github.com/gorilla/mux"
)

const (
	defaultSvcAccountEmailDomain = "@padl.adrianosela.com"
)

func (s *Service) addProjectEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/project").Handler(s.Auth(s.createProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/project/{name}").Handler(s.Auth(s.getProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/project/{name}/keys").Handler(s.Auth(s.getProjectKeysHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{name}").Handler(s.Auth(s.deleteProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/projects").Handler(s.Auth(s.listProjectsHandler))

	s.Router.Methods(http.MethodPost).Path("/project/{name}/user").Handler(s.Auth(s.addUserHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{name}/user").Handler(s.Auth(s.removeUserHandler))

	s.Router.Methods(http.MethodPost).Path("/project/{name}/service_account").Handler(s.Auth(s.createServiceAccountHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{name}/service_account").Handler(s.Auth(s.removeServiceAccountHandler))
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
	// check name is unique before doing anything
	exists, err := s.database.ProjectExists(projPl.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not check if project exists %s", err)))
		return
	}
	if exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("provided project name is taken")))
		return
	}
	// create shared team key for project and save it
	pKey, err := kms.NewPrivateKey(projPl.KeyBits, projPl.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not create project key: %s", err)))
		return
	}
	if err = s.keystore.PutPrivKey(pKey); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not save project key: %s", err)))
		return
	}
	pub, err := pKey.Pub()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not extract public key: %s", err)))
		return
	}
	if err = s.keystore.PutPubKey(pub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not save project pub key: %s", err)))
		return
	}
	// create project object and save it
	project := project.NewProject(projPl.Name, projPl.Description, claims.Subject, pKey.ID)
	if err := s.database.PutProject(project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not save new project: %s", err)))
		return
	}
	// add project to user claims
	user, err := s.database.GetUser(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(fmt.Sprintf("unable to get user from the database: %s", err)))
		return
	}
	user.AddProject(project.Name)
	if err := s.database.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not update user: %s", err)))
		return
	}

	// create and send padlFile
	pf := &padlfile.File{
		Data: padlfile.Body{
			Project:     project.Name,
			Variables:   make(map[string]string),
			MemberKeys:  []string{user.KeyID},
			ServiceKeys: []string{},
			SharedKey:   pKey.ID,
		},
	}

	byt, err := json.Marshal(&pf)
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
	claims := GetClaims(r)
	// project name from request URL
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
		return
	}
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get project: %s", err)))
		return
	}

	if _, ok := p.Members[claims.Subject]; !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("requesting user not in project: %s", err)))
		return
	}

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

func (s *Service) getProjectKeysHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)

	// project name from request URL
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
		return
	}

	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get project: %s", err)))
		return
	}
	if _, ok := p.Members[claims.Subject]; !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("requesting user not in project: %s", err)))
		return
	}

	memKeyIDs := []string{}
	for member := range p.Members {
		user, err := s.database.GetUser(member)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("could not get a member user from db: %s", err)))
			return
		}
		memKeyIDs = append(memKeyIDs, user.KeyID)
	}

	svcKeyIDs := []string{}
	for _, svcKeyID := range p.ServiceAccounts {
		svcKeyIDs = append(svcKeyIDs, svcKeyID)
	}

	getProjectKeysResp := payloads.GetProjectKeysReponse{
		Name:       p.Name,
		MemberKeys: memKeyIDs,
		DeployKeys: svcKeyIDs,
		ProjectKey: p.ProjectKey,
	}

	byt, err := json.Marshal(&getProjectKeysResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal project keys json: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}

func (s *Service) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// name from GET params
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
		return
	}
	// get project from db
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get project: %s", err)))
		return
	}
	// check caller is owner, else reject request
	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can delete a project")))
		return
	}
	// delete project's private key
	if err = s.keystore.DeletePrivKey(p.ProjectKey); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to delete project's private key: %s", err)))
	}
	// delete project's public key
	if err = s.keystore.DeletePubKey(p.ProjectKey); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to delete project's public key: %s", err)))
	}
	// delete all service account public keys
	for _, keyID := range p.ServiceAccounts {
		if err = s.keystore.DeletePubKey(keyID); err != nil {
			// fail open, just log
			log.Printf("unable to delete project service account key %s: %s", keyID, err)
		}
	}
	// remove project from all users
	for member := range p.Members {
		if u, err := s.database.GetUser(member); err == nil {
			u.RemoveProject(p.Name)
			s.database.UpdateUser(u) // note the ignored error here
		}
	}
	// delete project
	if err = s.database.DeleteProject(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to delete project: %s", err)))
	}
	// send success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("project %s deleted successfully!", name)))
}

func (s *Service) addUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
		return
	}
	// read request body
	var addUserPl *payloads.AddUserToProjectRequest
	if err := unmarshalRequestBody(r, &addUserPl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not unmarshal request body: %s", err)))
		return
	}
	// validate payload data
	if err := addUserPl.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not validate request: %s", err)))
		return
	}
	exists, err := s.database.UserExists(addUserPl.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("problem getting users from db: %s", err)))
		return
	}
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("user %s does not exist", addUserPl.Email)))
		return
	}
	// fetch project from database
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	user, err := s.database.GetUser(addUserPl.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to get user from the database: %s", err)))
		return
	}

	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can add users to a project")))
		return
	}

	user.AddProject(p.Name)
	if err := s.database.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not update user: %s", err)))
		return
	}

	if err = p.AddUser(addUserPl.Email, privilege.Level(addUserPl.PrivilegeLvl)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could add user to project: %s", err)))
		return
	}
	// update project
	if err := s.database.UpdateProject(p); err != nil {
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
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
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
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	user, err := s.database.GetUser(rmUserPl.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to get user from the database: %s", err)))
		return
	}

	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can remove users from a project")))
		return
	}

	user.RemoveProject(p.Name)
	if err := s.database.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not update user: %s", err)))
		return
	}

	if rmUserPl.Email == claims.Subject {
		if p.HasUser(claims.Subject) && p.Members[claims.Subject] == privilege.PrivilegeLvlOwner {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("owners cannot remove themselves from projects"))
			return
		}
	}
	// update project
	p.RemoveUser(rmUserPl.Email)
	if err := s.database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("user %s removed from project %s successfully!", rmUserPl.Email, p.Name)))
	return
}

func (s *Service) createServiceAccountHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	// get project name from request URL
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
		return
	}
	// get payload data (svc acct name + pub key)
	var dkeyPl *payloads.CreateServiceAccountRequest
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
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}
	// check if caller is authorized to create a service account for project
	if _, ok := p.Members[claims.Subject]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("User not in requested Project: %s", err)))
		return
	}
	if p.Members[claims.Subject] < privilege.PrivilegeLvlEditor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Only owners and editors can create service accounts")))
		return
	}
	// create padl pub key object for svc account and store it publicly
	pub, err := kms.NewPublicKey(dkeyPl.PubKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := s.keystore.PutPubKey(pub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not store service account's public key: %s", err)))
		return
	}
	// build a service account email of the form {name}.{project_name}@{padl_hostname}
	// i.e. cicd.my-first-project@api.padl.com
	svcEmail := dkeyPl.ServiceAccountName + "." + p.Name + defaultSvcAccountEmailDomain
	keyToken, err := s.authenticator.GenerateJWT(svcEmail, auth.ServiceAccountAudience)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not generate keyToken: %s", err)))
		return
	}
	// add service account to project object and save it
	p.SetServiceAccount(dkeyPl.ServiceAccountName, pub.ID)
	if err := s.database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	// marshall response
	byt, err := json.Marshal(&payloads.CreateServiceAccountResponse{Token: keyToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal project json: %s", err)))
		return
	}
	// send success
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
	return
}

func (s *Service) removeServiceAccountHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
		return
	}
	// get payload data
	var deleteKeyPl payloads.DeleteServiceAccountRequest
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
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	if _, ok := p.Members[claims.Subject]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("User not in requested Project: %s", p.Name)))
		return
	}

	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Only Owners can remove service accounts")))
		return
	}

	// update project
	p.RemoveServiceAccount(deleteKeyPl.ServiceAccountName)
	if err := s.database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed service account from project %s", p.Name)))
	return
}

func (s *Service) listProjectsHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)

	user, err := s.database.GetUser(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not get user from db: %s", err)))
		return
	}

	projects, err := s.database.ListProjects(user.Projects)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not get list of projects: %s", err)))
		return
	}

	ps := []*project.Summary{}
	for _, p := range projects {
		ps = append(ps,
			&project.Summary{
				Name:        p.Name,
				Description: p.Description,
			},
		)
	}

	listProjResp := payloads.ListProjectsResponse{
		Projects: ps,
	}
	// marshall response
	byt, err := json.Marshal(&listProjResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal projects summary json: %s", err)))
		return
	}
	// send resp
	w.WriteHeader(http.StatusOK)
	w.Write(byt)
}
