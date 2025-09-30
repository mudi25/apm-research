package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"research-apm/services/api/cmd/config"
	"syscall"

	_ "go.elastic.co/apm/v2"
)

func main() {
	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize application (includes HTTP server, DB, NATS, New Relic)
	app, err := config.NewApp(ctx)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return
	}
	defer app.Shutdown(ctx)

	// Start HTTP server in a separate goroutine
	fmt.Println("[INFO] Run User API On", app.Server.Addr)
	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("[ERROR] Failed to run server:", err.Error())
			app.Shutdown(ctx)
			os.Exit(1)
		}
	}()

	// Wait for termination signal (Ctrl+C / Docker stop / etc.)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("[INFO] Shutdown User API")
}
