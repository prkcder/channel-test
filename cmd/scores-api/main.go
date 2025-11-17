package main

import (
	"channel-test/internal/api"
	"channel-test/internal/consumer"
	"channel-test/internal/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort     = "8000"
	sseURL          = "http://live-test-scores.herokuapp.com/scores"
	shutdownTimeout = 10 * time.Second
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize store
	dataStore := store.NewMemoryStore()
	log.Println("Initialized in-memory store")

	// Initialize SSE consumer
	sseConsumer := consumer.NewSSEConsumer(sseURL, dataStore)

	// Start SSE consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println("Starting SSE consumer...")
		if err := sseConsumer.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("SSE consumer error: %v", err)
		}
	}()

	// Initialize HTTP handler and router
	handler := api.NewHandler(dataStore)
	router := api.NewRouter(handler)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in background
	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cancel SSE consumer
	cancel()

	// Gracefully shut down HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
