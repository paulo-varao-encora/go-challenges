package main

import (
	"example/challenges/internal/api"
	"log"
	"net/http"
)

func main() {
	server, err := api.NewTaskServer()

	if err != nil {
		log.Fatalf("failed to create server %v", err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
