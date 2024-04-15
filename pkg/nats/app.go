package nats

import (
	"context"

	"go.uber.org/zap"
)

type App interface {
	Provision(ctx context.Context) error
	Start() error
	Stop() error
}

type NatsApplication struct {
	logger *zap.Logger
	ctx    context.Context
}

// Provision sets up the application.
func (a *NatsApplication) Provision(ctx context.Context) error {
	a.logger = zap.NewExample()
	a.ctx = ctx
	return nil
}

// Start starts the application.
func (a *NatsApplication) Start() error {
	a.logger.Info("NatsApplication started")
	return nil
}

// Stop stops the application.
func (a *NatsApplication) Stop() error {
	a.logger.Info("NatsApplication stopped")
	return nil
}
