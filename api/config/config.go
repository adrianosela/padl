package config

import (
	"log"
	"time"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var (
	// injected at build-time
	version string
)

type External struct {
	Config struct {
		Test string `yaml:"test"`
	} `yaml:"conf"`
}

// Config holds the service configuration
// necessary for endpoints to respond to requests
type Config struct {
	Version    string // server git hash
	DeployTime time.Time
	External   External
}

func GetConfig(filePath string) Config {
	config := Config{
		Version:    version,
		DeployTime: time.Now(),
		External:   parseYaml(filePath),
	}
	return config
}

func parseYaml(filePath string) External {
	externalConfig := External{}

	yamlFile, err := ioutil.ReadFile(filePath)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }

	err = yaml.Unmarshal(yamlFile, &externalConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return externalConfig
}
