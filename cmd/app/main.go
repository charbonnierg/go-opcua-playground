package main

import (
	"context"
	"log"
	"playground/pkg/nats"
)

// Very minimalist lifecycle of an NATS app in go
func main() {
	ctx := context.Background()
	// Create the app
	app := nats.NewApplication("nats://localhost:4222")
	// Provision the app
	if err := app.Provision(ctx); err != nil {
		log.Fatalf("Failed to provision the app: %v", err)
	}
	// Start the app
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start the app: %v", err)
	}
	// Stop the app
	if err := app.Stop(); err != nil {
		log.Fatalf("Failed to stop the app: %v", err)
	}
}
