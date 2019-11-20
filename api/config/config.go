package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config holds the service configuration
// necessary for endpoints to respond to requests
type Config struct {
	Version    string // server git hash
	DeployTime time.Time
	Env        string `yaml:"env"`
	Port       string `yaml:"port"`
	Debug      bool   `yaml:"debug"`

	Database struct {
		ConnectionString       string `yaml:"connectionString"`
		Name                   string `yaml:"name"`
		UsersCollectionName    string `yaml:"usersCollectionName"`
		ProjectsCollectionName string `yaml:"projectsCollectionName"`
	} `yaml:"database"`

	Auth struct {
		SigningKey string `yaml:"signingKey"`
	} `yaml:"auth"`
}

// BuildConfig returns a populated config struct from a yaml file
func BuildConfig(filePath, version string) *Config {
	config := configFromYaml(filePath)

	// When running on Google App Engine, the PORT env
	// variable is set by the runtime. If set, we will
	// serve on the port specified there.
	if port := os.Getenv("PORT"); port != "" {
		config.Port = fmt.Sprintf(":%s", port)
	}

	config.DeployTime = time.Now()
	config.Version = version

	return config
}

func configFromYaml(filePath string) *Config {
	config := &Config{}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(yamlFile, config); err != nil {
		log.Fatal(err)
	}

	return config
}
