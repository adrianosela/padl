package keymgr

import (
	"fmt"
	"io/ioutil"
	"os"
)

// FSManager is a File System key manager
type FSManager struct {
	basePath string
}

// NewFSManager is the FSManager constructor
func NewFSManager(path string) (*FSManager, error) {
	bPath := fmt.Sprintf("%s/keys", path)
	if _, err := os.Stat(bPath); os.IsNotExist(err) {
		if err := os.Mkdir(bPath, 755); err != nil {
			return nil, fmt.Errorf("could not create new directory %s: %s", path, err)
		}
		path = fmt.Sprintf("%s/privs", bPath)
		if err := os.Mkdir(path, 755); err != nil {
			return nil, fmt.Errorf("could not create new directory %s: %s", path, err)
		}
		path = fmt.Sprintf("%s/pubs", bPath)
		if err := os.Mkdir(path, 755); err != nil {
			return nil, fmt.Errorf("could not create new directory %s: %s", path, err)
		}
	}
	return &FSManager{
		basePath: bPath,
	}, nil
}

// PutPriv saves a private key by id
func (m *FSManager) PutPriv(id string, blob string) error {
	if err := ioutil.WriteFile(fmt.Sprintf("%s/privs/%s.priv", m.basePath, id), []byte(blob), 0644); err != nil {
		return fmt.Errorf("could not write private key to file system: %s", err)
	}
	return nil
}

// GetPriv gets a private key by id
func (m *FSManager) GetPriv(id string) (string, error) {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/privs/%s.priv", m.basePath, id))
	if err != nil {
		return "", fmt.Errorf("could not read private key from file system: %s", err)
	}
	return string(dat), nil
}

// PutPub saves a public key by id
func (m *FSManager) PutPub(id string, blob string) error {
	if err := ioutil.WriteFile(fmt.Sprintf("%s/pubs/%s.pub", m.basePath, id), []byte(blob), 0644); err != nil {
		return fmt.Errorf("could not write public key to file system: %s", err)
	}
	return nil
}

// GetPub gets a public key by id
func (m *FSManager) GetPub(id string) (string, error) {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/pubs/%s.pub", m.basePath, id))
	if err != nil {
		return "", fmt.Errorf("could not read public key from file system: %s", err)
	}
	return string(dat), nil
}
