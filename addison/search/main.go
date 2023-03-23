package main

import (
	"fmt"
	"log"
	"net/http"

	"search/resources"
)

func main() {
	fmt.Println("Search microservices started: listening on port 3001.")
	log.Fatal(http.ListenAndServe(":3001", resources.Router()))
}
