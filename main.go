package main

import (
	"context"
	"fmt"
	"github/shaolim/momon/internal/receipt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	openAIClient := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_APIKEY")))
	receipt := receipt.New(&openAIClient)

	// Check if file path is provided as command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path-to-receipt-image>")
	}

	filePath := os.Args[1]

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", filePath)
	}

	res, err := receipt.ReadReceipt(context.Background(), filePath)
	if err != nil {
		log.Fatalf("Error read receipt: %v", err)
	}
	fmt.Println(res.String())
}
