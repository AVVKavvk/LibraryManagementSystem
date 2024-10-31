package main

import (
	"log"
	"net/http"

	"github.com/AVVKavvk/LMS/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")

	}
}
func Handler(w http.ResponseWriter, r *http.Request) {
	handler := server.Server()
	handler.ServeHTTP(w, r)
}
