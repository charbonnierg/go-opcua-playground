package main

import (
	"context"
	"os"
	"os/signal"
	"playground/pkg/nats"

	"go.uber.org/zap"
)

// Very minimalist lifecycle of an NATS app in go
func main() {
	logger := zap.NewExample()
	ctx := context.Background()
	// Create the app
	app := nats.NewApplication("nats://localhost:4222")
	// Provision the app
	if err := app.Provision(ctx); err != nil {
		logger.Fatal("Failed to provision the app", zap.Error(err))
	}
	// Start the app
	if err := app.Start(); err != nil {
		logger.Fatal("Failed to start the app", zap.Error(err))
	}
	// Stop the app on exit
	defer func() {
		if err := app.Stop(); err != nil {
			logger.Fatal("Failed to stop the app", zap.Error(err))
		}
	}()
	// catch ctrl-c and gracefully shutdown the server.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	defer signal.Stop(sigch)
	logger.Info("Press Ctrl+C to exit")

	<-sigch
	logger.Info("Shutting down the app...")
}
