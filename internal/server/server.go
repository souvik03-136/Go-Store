// internal/server/server.go

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/config"
	internaldb "github.com/souvik03-136/Go-Store/internal/database"
	"github.com/souvik03-136/Go-Store/internal/storage"
)

// Server owns the HTTP server, the DB pool, and the storage backend.
type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

// NewServer builds the full application: connects to the DB, initialises the
// storage backend, wires all dependencies, and registers all routes.
func NewServer(cfg *config.Config) (*Server, error) {
	// ── Database ──────────────────────────────────────────────────────────────
	db, err := internaldb.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	log.Printf("Connected to %s database on %s:%s",
		cfg.Database.Driver, cfg.Database.Host, cfg.Database.Port)

	// ── Storage backend ───────────────────────────────────────────────────────
	store, err := buildStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialising storage backend: %w", err)
	}
	log.Printf("Storage provider: %s", cfg.StorageProvider)

	// ── Router ────────────────────────────────────────────────────────────────
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery()) // always recover from panics

	initRoutes(router, db, store, cfg.JWT.SecretKey)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		cfg: cfg,
	}, nil
}

// Start begins serving requests and blocks until an OS interrupt signal is
// received, then performs a graceful 10-second shutdown.
func (s *Server) Start() error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("Server listening on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Printf("Received signal %s — shutting down gracefully", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("Server stopped cleanly")
	return nil
}

// buildStorage constructs the appropriate storage.Storage implementation from cfg.
func buildStorage(cfg *config.Config) (storage.Storage, error) {
	switch cfg.StorageProvider {
	case "s3":
		return storage.NewS3Storage(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			cfg.AWS.Region,
			cfg.AWS.BucketName,
		)
	case "gcs":
		ctx := context.Background()
		return storage.NewGCSStorage(
			ctx,
			cfg.GoogleCloud.CredentialsFile,
			cfg.GoogleCloud.BucketName,
		)
	default:
		return nil, fmt.Errorf("unsupported storage provider %q — use \"s3\" or \"gcs\"", cfg.StorageProvider)
	}
}
