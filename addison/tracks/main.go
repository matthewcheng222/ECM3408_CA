package main

import (
	"fmt"
	"log"
	"net/http"
	"tracks/repository"
	"tracks/resources"
)

func main() {
	repository.Init()
	repository.Clear()
	repository.Create()

	fmt.Println("Tracks microservices started: listening on port 3000.")
	log.Fatal(http.ListenAndServe(":3000", resources.Router()))
}
