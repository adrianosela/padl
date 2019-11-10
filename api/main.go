package main

import (
	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/service"
	"log"
	"net/http"
)

const filePath = "./config/config.yaml"

func main() {
	c := config.GetConfig(filePath)

	svc := service.NewPadlService(c)

	if err := http.ListenAndServe(c.Port, svc.Router); err != nil {
		log.Fatal(err)
	}
}
