package main

import (
	"log"
	"net/http"

	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/service"
)

const filePath = "./config/config.yaml"

var (
	version string // injected at build-time
)

func main() {
	c := config.BuildConfig(filePath, version)

	svc := service.NewPadlService(c)

	if err := http.ListenAndServe(c.Port, svc.Router); err != nil {
		log.Fatal(err)
	}
}
