package main

import (
	"log"

	"github.com/AVVKavvk/LMS/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")

	}
}
func main() {
	server.Server()
}
