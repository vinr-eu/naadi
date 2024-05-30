package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/vinr-eu/naadi/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	configPath := "config/config.yaml"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to load config", "err", err)
	}
	slog.Info("Config loaded", "receivers", len(cfg.Receivers))

	idleConnectionsClosed := make(chan struct{})

	mux := http.NewServeMux()

	startHTTPServer(mux, idleConnectionsClosed)

	<-idleConnectionsClosed
}

func startHTTPServer(mux *http.ServeMux, idleConnectionsClosed chan struct{}) {
	// Create http server
	address := ""
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	serverAddr := fmt.Sprintf("%s:%s", address, serverPort)
	srv := &http.Server{Addr: serverAddr, Handler: mux}

	// Prepare http server for graceful shutdown.
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout.
			slog.Error("Server shutdown failed", "err", err)
		}
		slog.Info("Server shutdown")
		close(idleConnectionsClosed)
	}()

	// Start http server.
	go func() {
		slog.Info("Server started")
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener.
			slog.Error("Server startup failed", "err", err)
		}
	}()
}
