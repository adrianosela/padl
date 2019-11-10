package main

import (
	"github.com/adrianosela/padl/api/service"
	"log"
	"net/http"
)

func main() {
	c := getConfig()

	svc := service.NewPadlService(c)

	if err := http.ListenAndServe(":80", svc); err != nil {
		log.Fatal(err)
	}
}
