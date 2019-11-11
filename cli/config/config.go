package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

const (
	defaultConfigDirName = ".padl"
)

// Config is the padl cli configuration
type Config struct {
	HostURL string `json:"host_url"`
	User    string `json:"user,omitempty"`
	Token   string `json:"auth_token,omitempty"`
}

// GetDefaultPath returns the best place to save / look-for a config directory.
// Using this function for both saving and reading config or resources (with
// the same OS) guarantees that a directory will be found
func GetDefaultPath() string {
	usr, err := user.Current()
	if err != nil {
		// settle for rootdir
		return fmt.Sprintf("/%s", defaultConfigDirName)
	}
	return fmt.Sprintf("%s/%s", usr.HomeDir, defaultConfigDirName)
}

// SetConfig writes a configuration file to the given path
func SetConfig(c *Config, path string) error {
	if c == nil {
		return errors.New("config cannot be nil")
	}
	if c.HostURL == "" {
		return errors.New("url cannot be empty")
	}
	if path == "" {
		path = GetDefaultPath()
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 755); err != nil {
			return fmt.Errorf("could not create new directory %s: %s", path, err)
		}
	}
	f, err := os.Create(fmt.Sprintf("%s/config", path))
	if err != nil {
		return fmt.Errorf("could not create new file %s/config: %s", path, err)
	}
	byt, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not marshal configuration file: %s", err)
	}
	if _, err := f.Write(byt); err != nil {
		return fmt.Errorf("could not write configuration file: %s", err)
	}
	return nil
}

// GetConfig returns the configuration at a given path
func GetConfig(path string) (*Config, error) {
	if path == "" {
		path = GetDefaultPath()
	}
	return readFSConfig(path)
}

func readFSConfig(path string) (*Config, error) {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/config", path))
	if err != nil {
		return nil, fmt.Errorf("could not read configuration file %s: %s", path, err)
	}
	var c *Config
	if err = json.Unmarshal(dat, &c); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %s", err)
	}
	return c, nil
}
