package main

import (
	"log"
	"net/http"

	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/service"
)

const filePath = "./config/config.yaml"

func main() {
	c := config.GetConfig(filePath)

	svc := service.NewPadlService(c)

	if err := http.ListenAndServe(svc.Config.Port, svc.Router); err != nil {
		log.Fatal(err)
	}
}
