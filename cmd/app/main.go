package main

import (
	"context"
	"playground/pkg/nats"
)

// Very minimalist lifecycle of an NATS app in go
func main() {
	// Create the app
	app := nats.NatsApplication{}
	// Provision the app
	app.Provision(context.Background())
	// Start the app
	app.Start()
	// Stop the app
	app.Stop()
}
