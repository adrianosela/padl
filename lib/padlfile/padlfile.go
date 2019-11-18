package padlfile

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
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
	Project    string            `json:"project_id" yaml:"project_id"` // id of the project for this padlfile
	Variables  map[string]string `json:"variables" yaml:"variables"`   // map of ENV_VAR to hash
	MemberKeys []string          `json:"keys" yaml:"keys"`             // project member key ids
	SharedKey  string            `json:"shared_key" yaml:"shared_key"` // shared project key id
}

// File represents the entire contents of a Padlfile
type File struct {
	Data Body   `json:"data" yaml:"data"`
	HMAC string `json:"HMAC" yaml:"HMAC"`
}

// HashAndSign returns the finalized padlfile contents
// including a signed hash (HMAC)
func (b *Body) HashAndSign(secret []byte) (*File, error) {
	jsonByt, err := json.Marshal(&b)
	if err != nil {
		return nil, fmt.Errorf("could not encode body contents: %s", err)
	}
	h := hmac.New(sha512.New, secret)
	h.Write(jsonByt)
	return &File{
		Data: *b,
		HMAC: hex.EncodeToString(h.Sum(nil)),
	}, nil
}

// VerifySignature verifies the hash and signature on the file
func (f *File) VerifySignature(secret []byte) (bool, error) {
	decoded, err := hex.DecodeString(f.HMAC)
	if err != nil {
		return false, fmt.Errorf("could not decode file's HMAC: %s", err)
	}
	jsonByt, err := json.Marshal(&f.Data)
	if err != nil {
		return false, fmt.Errorf("could not encode body contents: %s", err)
	}
	h := hmac.New(sha512.New, secret)
	h.Write(jsonByt)
	return hmac.Equal(h.Sum(nil), decoded), nil
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
