package config

import (
	"io/ioutil"
	"log"
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
}

// BuildConfig returns a populated config struct from a yaml file
func BuildConfig(filePath, version string) *Config {
	config := configFromYaml(filePath)

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
