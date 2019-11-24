package padlfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Body represents the body of a Padlfile
type Body struct {
	Project     string            `json:"project_id" yaml:"project_id"`     // id of the project for this padlfile
	Variables   map[string]string `json:"variables" yaml:"variables"`       // map of ENV_VAR secret
	MemberKeys  []string          `json:"user_keys" yaml:"user_keys"`       // project member key ids
	ServiceKeys []string          `json:"service_keys" yaml:"service_keys"` // service account key ids
	SharedKey   string            `json:"shared_key" yaml:"shared_key"`     // shared project key id
}

// File represents the entire contents of a Padlfile
type File struct {
	Data Body `json:"data" yaml:"data"`
}

// ReadPadlfile reads a padlfile from the given path
func ReadPadlfile(path string) (*File, error) {
	if path == "" {
		padlfiles, err := filepath.Glob("./.padlfile.*")
		if err != nil {
			return nil, fmt.Errorf("error looking for padlfile: %s", err)

		}
		if len(padlfiles) == 0 {
			return nil, errors.New("no padlfile found in current directory")
		}
		path = padlfiles[0]
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read padlfile file %s: %s", path, err)
	}
	var f File
	if strings.HasSuffix(path, ".yaml") {
		if err := yaml.Unmarshal(dat, &f); err != nil {
			return nil, fmt.Errorf("could not unmarshal .yaml file: %s", err)
		}
	} else {
		if err := json.Unmarshal(dat, &f); err != nil {
			return nil, fmt.Errorf("could not unmarshal .json file: %s", err)
		}
	}
	return &f, nil
}

// Write writes the padlfile to a path
func (f *File) Write(path string) error {
	var fbyt []byte
	var err error
	// marshal onto encoding as per path
	if strings.HasSuffix(path, ".yaml") {
		if fbyt, err = yaml.Marshal(&f); err != nil {
			return fmt.Errorf("could not marshal padlfile to .yaml file: %s", err)
		}
	} else {
		if fbyt, err = json.Marshal(&f); err != nil {
			return fmt.Errorf("could not marshal padlfile to .json file: %s", err)
		}
	}
	// create and write padlfile
	fd, err := os.Create(fmt.Sprintf("%s", path))
	if err != nil {
		return fmt.Errorf("could not create padlfile: %s", err)
	}
	if _, err = fd.Write(fbyt); err != nil {
		return fmt.Errorf("could not write padlfile: %s", err)
	}
	return nil
}
