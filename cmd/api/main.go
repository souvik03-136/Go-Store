package main

import (
	"log"

	"github.com/souvik03-136/Go-Store/internal/config"
	"github.com/souvik03-136/Go-Store/internal/server"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize the server
	s := server.NewServer()

	// Start the server
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
