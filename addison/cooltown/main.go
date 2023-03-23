package main

import (
	"cooltown/resources"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Cooltown microservices started: listening on port 3002.")
	log.Fatal(http.ListenAndServe(":3002", resources.Router()))
}
