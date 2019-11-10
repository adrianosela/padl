package main

import (
	"github.com/adrianosela/padl/api/service"
	"github.com/adrianosela/padl/api/store"
	"time"
)

var (
	// injected at build-time
	version string
)

func getConfig() service.Config {
	config := service.Config{
		Version:    version,
		DeployTime: time.Now(),
		Database:   store.NewMockDatabase(),
	}
	return config
}
