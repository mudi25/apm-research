package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"research-apm/services/alert/cmd/config"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := config.NewApp(ctx)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return
	}
	defer app.Shutdown()
	fmt.Println("[INFO] Run Alert Worker")

	// Wait for termination signal (Ctrl+C / Docker stop / etc.)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("[INFO] Shutdown Alert Worker")
}
