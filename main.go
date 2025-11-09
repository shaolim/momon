package main

import (
	"context"
	"github/shaolim/momon/internal/messaging"
	"github/shaolim/momon/pkg/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Start the server
	port := "8080"
	log.Printf("Starting HTTP server on port %s...", port)
	s, err := server.New(port)
	if err != nil {
		log.Fatal("failed to initiate the server:", err)
	}

	m := messaging.New(nil)

	if err := s.ServeHTTPHandler(context.Background(), m.Routes()); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
