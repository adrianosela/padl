package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/padl/api/project"
	"github.com/gorilla/mux"
)

func (c *Config) addProjectEndpoints(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/newProject").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/addOwner").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/addEditor").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/addReader").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/removeOwner").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/removeEditor").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/removeReader").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/addSecret").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/removeSecret").HandlerFunc(c.createNewProjectHandler)
	r.Methods(http.MethodPost).Path("/createDeployKey").HandlerFunc(c.createNewProjectHandler)

}

// Change tules endpoint
// Add owner, editor, reader endpoint
// Remove owner, editor, reader endpoint
// Add secret ednpoint
// Remove secret endpoint
// Create deploy key
//

func (c *Config) createNewProjectHandler(w http.ResponseWriter, r *http.Request) {
	var proj *newProjectRequest
	if err := unmarshalRequestBody(r, &proj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	if err := proj.validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new project: %s", err)))
		return
	}

	rules := project.Rules{
		RequireMFA:     proj.RequireMFA,
		RequireTeamKey: proj.RequireTeamKey,
	}
	project := project.NewProject(proj.Token, proj.Name, rules)

	if err := c.Database.CreateProject(project); err != nil {
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

func (c *Config) addOwner(w http.ResponseWriter, r *http.Request) {
	var req *addOwnerRequest
	if err := unmarshalRequestBody(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshall request body"))
	}
	if err := req.validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not create new project: %s", err)))
		return
	}

	rules := project.Rules{
		RequireMFA:     proj.RequireMFA,
		RequireTeamKey: proj.RequireTeamKey,
	}
	project := project.NewProject(proj.Token, proj.Name, rules)

	if err := c.Database.CreateProject(project); err != nil {
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
