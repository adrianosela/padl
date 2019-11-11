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
	}
	return &FSManager{
		basePath: bPath,
	}, nil
}

// PutKey saves a key by id
func (m *FSManager) PutKey(id string, blob string) error {
	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", m.basePath, id), []byte(blob), 0644); err != nil {
		return fmt.Errorf("could not write key to file system: %s", err)
	}
	return nil
}

// GetKey gets a key by id
func (m *FSManager) GetKey(id string) (string, error) {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", m.basePath, id))
	if err != nil {
		return "", fmt.Errorf("could not read key from file system: %s", err)
	}
	return string(dat), nil
}
