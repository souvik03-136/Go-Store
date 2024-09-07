package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

// NewServer creates and returns a new HTTP server with routes initialized.
func NewServer() *Server {
	router := gin.Default()
	InitRoutes(router)

	httpServer := &http.Server{
		Addr:         ":8080",          // Port the server will run on
		Handler:      router,           // Our router
		ReadTimeout:  10 * time.Second, // Prevents slow client attacks
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return &Server{httpServer: httpServer}
}

// Start begins the server and handles graceful shutdown.
func (s *Server) Start() error {
	go func() {
		log.Println("Starting server on port 8080...")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 8080: %v\n", err)
		}
	}()

	// Graceful shutdown handling
	return waitForShutdown(s.httpServer)
}

// waitForShutdown gracefully shuts down the server.
func waitForShutdown(server *http.Server) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
