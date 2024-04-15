package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type App interface {
	Provision(ctx context.Context) error
	Start() error
	Stop() error
}

func NewApplication(serverUrl string) *NatsApplication {
	return &NatsApplication{
		NatsServerUrl: serverUrl,
	}
}

type NatsApplication struct {
	logger            *zap.Logger
	ctx               context.Context
	conn              *nats.Conn
	NatsServerUrl     string
	ConnectionTimeout time.Duration
}

// Provision sets up the application.
func (a *NatsApplication) Provision(ctx context.Context) error {
	a.logger = zap.NewExample()
	a.ctx = ctx
	return nil
}

// Start starts the application.
func (a *NatsApplication) Start() error {
	timeout := a.ConnectionTimeout
	// Default timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	// Connect to the NATS server
	options := nats.GetDefaultOptions()
	options.Url = a.NatsServerUrl
	options.Timeout = timeout
	a.logger.Info("Connecting to NATS server", zap.String("url", a.NatsServerUrl))
	conn, err := options.Connect()
	if err != nil {
		return err
	}
	a.conn = conn
	a.logger.Info("Connected to NATS server")
	return nil
}

// Stop stops the application.
func (a *NatsApplication) Stop() error {
	if a.conn != nil {
		a.logger.Info("Disconnecting from NATS server")
		a.conn.Close()
	}
	a.logger.Info("Connection to NATS server closed")
	return nil
}

// Typeguards
var _ App = &NatsApplication{}
