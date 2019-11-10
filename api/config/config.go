package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

var (
	// injected at build-time
	version string
)

// Config holds the service configuration
// necessary for endpoints to respond to requests
type Config struct {
	Version    string // server git hash
	DeployTime time.Time
	Env        string `yaml:"env"`
}

// GetConfig returns a populated config struct from a yaml file
func GetConfig(filePath string) Config {
	config := configFromYaml(filePath)

	config.DeployTime = time.Now()
	config.Version = version

	return config
}

func configFromYaml(filePath string) Config {
	config := Config{}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
