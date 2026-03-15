package main

import (
	"log"

	"github.com/souvik03-136/Go-Store/internal/config"
	"github.com/souvik03-136/Go-Store/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	s, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}