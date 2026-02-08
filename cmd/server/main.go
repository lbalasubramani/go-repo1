package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lbalasubramani/go-repo1/internal/config"
	"github.com/lbalasubramani/go-repo1/internal/handlers"
	"github.com/lbalasubramani/go-repo1/pkg/logger"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger := logger.New(cfg.LogLevel)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(appLogger)

	// Setup router
	mux := http.NewServeMux()
	
	// Register routes
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/readiness", handlers.ReadinessCheck)
	mux.HandleFunc("/items", itemHandler.HandleItems)
	mux.HandleFunc("/items/", itemHandler.HandleItemByID)

	// Apply middleware
	handler := handlers.LoggingMiddleware(appLogger)(
		handlers.RecoveryMiddleware(appLogger)(mux),
	)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info(fmt.Sprintf("Starting server on port %d", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error(fmt.Sprintf("Server error: %v", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
		os.Exit(1)
	}

	appLogger.Info("Server exited gracefully")
}
