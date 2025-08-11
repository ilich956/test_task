package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forte/internal/config"
	"forte/internal/handlers"

	"go.temporal.io/sdk/client"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	c, err := client.Dial(client.Options{HostPort: cfg.TemporalHost})
	if err != nil {
		log.Fatalf("failed to create temporal client %v", err)
	}
	defer c.Close()

	mux := setupRoutes(c)

	server := &http.Server{
		Addr:    cfg.ApiHost,
		Handler: mux,
	}

	go func() {
		log.Printf("Launcing server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func setupRoutes(c client.Client) *http.ServeMux {
	mux := http.NewServeMux()

	orderHandlers := handlers.NewOrderHandlers(c)
	workflowHandlers := handlers.NewWorkflowHandlers(c)
	healthHandlers := handlers.NewHealthHandlers(c)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))

	mux.HandleFunc("/api/orders/create", orderHandlers.CreateOrder)
	mux.HandleFunc("/api/orders/", orderHandlers.HandleOrder)
	mux.HandleFunc("/api/workflows/", workflowHandlers.GetWorkflowStatus)

	mux.HandleFunc("/health", healthHandlers.HealthCheck)

	return mux
}
