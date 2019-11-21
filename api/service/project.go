package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/auth"
	"github.com/adrianosela/padl/api/kms"
	"github.com/adrianosela/padl/api/payloads"
	"github.com/adrianosela/padl/api/privilege"
	"github.com/adrianosela/padl/api/project"
	"github.com/adrianosela/padl/lib/padlfile"
	"github.com/gorilla/mux"
)

func (s *Service) addProjectEndpoints() {
	s.Router.Methods(http.MethodPost).Path("/project").Handler(s.Auth(s.createProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/project/{name}").Handler(s.Auth(s.getProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/project/{name}/keys").Handler(s.Auth(s.getProjectKeysHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{name}").Handler(s.Auth(s.deleteProjectHandler))
	s.Router.Methods(http.MethodGet).Path("/projects").Handler(s.Auth(s.listProjectsHandler))

	s.Router.Methods(http.MethodPost).Path("/project/{name}/user").Handler(s.Auth(s.addUserHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{name}/user").Handler(s.Auth(s.removeUserHandler))

	s.Router.Methods(http.MethodPost).Path("/project/{name}/deploy_key").Handler(s.Auth(s.createDeployKeyHandler))
	s.Router.Methods(http.MethodDelete).Path("/project/{name}/deploy_key").Handler(s.Auth(s.removeDeployKeyHandler))
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

	// create padlfile
	body := &padlfile.Body{
		Project:    project.Name,
		Variables:  make(map[string]string),
		MemberKeys: []string{user.KeyID},
		SharedKey:  pKey.ID,
	}
	pf, err := body.HashAndSign([]byte("Some crazy secret"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to sign padlfile: %s", err)))
		return
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

	keyIDs := []string{}
	for member := range p.Members {
		user, err := s.database.GetUser(member)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("could not get a member user from db: %s", err)))
			return
		}
		keyIDs = append(keyIDs, user.KeyID)
	}

	getProjectKeysResp := payloads.GetProjectKeysReponse{
		Name:       p.Name,
		MemberKeys: keyIDs,
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

	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("only owners can delete a project")))
		return
	}
	// delete project's public key
	err = s.keystore.DeletePubKey(p.ProjectKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to delete project's public key: %s", err)))
	}
	// delete project's private key
	err = s.keystore.DeletePrivKey(p.ProjectKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to delete project's private key: %s", err)))
	}
	// delete project
	err = s.database.DeleteProject(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to delete project: %s", err)))
	}

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

	fmt.Println(claims.Subject) // REMOVE
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

func (s *Service) createDeployKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims := GetClaims(r)
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
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
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	if _, ok := p.Members[claims.Subject]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("User not in requested Project: %s", err)))
		return
	}

	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Only Owners can create deploy Keys")))
		return
	}

	keyEmail := dkeyPl.DeployKeyName + "." + p.Name + "@padl.adrianosela.com"

	keyToken, keyID, err := s.authenticator.GenerateJWT(keyEmail, auth.PadlDeployKeyAudience)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not generate keyToken: %s", err)))
		return
	}

	if err = p.SetDeployKey(dkeyPl.DeployKeyName, keyID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not set new deploy key: %s", err)))
		return
	}
	// update project
	if err := s.database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	// marshall response
	byt, err := json.Marshal(&payloads.CreateDeployKeyResponse{
		Token: keyToken,
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
	var name string
	if name = mux.Vars(r)["name"]; name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no project Name in request URL"))
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
	p, err := s.database.GetProject(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not find project: %s", err)))
		return
	}

	if _, ok := p.Members[claims.Subject]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("User not in requested Project: %s", err)))
		return
	}

	if p.Members[claims.Subject] < privilege.PrivilegeLvlOwner {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Only Owners can remove deploy Keys")))
		return
	}

	// update project
	p.RemoveDeployKey(deleteKeyPl.DeployKeyName)
	if err := s.database.UpdateProject(p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not update project: %s", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("successfully removed deploy key from project %s", p.Name)))
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
