package main

import (
	"log"
	"net/http"

	"example.com/user-management/internal/server"
)

func main() {
	log.Println("Starting server on port: 8080")
	err := http.ListenAndServe(":8080", server.New())
	if err != nil {
		log.Fatal(err)
	}
}
