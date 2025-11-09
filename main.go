package main

import (
	"context"
	"github/shaolim/momon/internal/messaging"
	"github/shaolim/momon/internal/serverenv"
	"github/shaolim/momon/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv"
	messagingapi "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
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

	messagingConfig := messaging.NewConfig()

	lineMessagingAPI, err := messagingapi.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal("failed to initiate line messaging API", err)
	}

	senv := serverenv.New(serverenv.WithLineMessagingAPI(lineMessagingAPI))

	m := messaging.New(messagingConfig, senv)

	if err := s.ServeHTTPHandler(context.Background(), m.Routes()); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
