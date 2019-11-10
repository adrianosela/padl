package main

import (
	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/service"
	"log"
	"net/http"
	"fmt"
)

var filePath = "./config.yaml"

func main() {
	c := config.GetConfig(filePath)
	fmt.Println(c)

	svc := service.NewPadlService(&c)

	if err := http.ListenAndServe(":80", svc.Router); err != nil {
		log.Fatal(err)
	}
}
