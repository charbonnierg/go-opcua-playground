package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"go.uber.org/zap"
)

func NewApplication(serverUrl string) *NatsApplication {
	return &NatsApplication{
		NatsServerUrl: serverUrl,
	}
}

type NatsApplication struct {
	logger *zap.Logger
	ctx    context.Context
	conn   *nats.Conn

	NatsServerUrl     string        `json:"url,omitempty"`
	ConnectionTimeout time.Duration `json:"connect_timeout,omitempty"`
}

// Provision sets up the application.
func (a *NatsApplication) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger()
	a.ctx = ctx
	return nil
}

// Start starts the application.
func (a *NatsApplication) Start() error {
	a.logger.Info("Connecting to NATS server", zap.String("url", a.NatsServerUrl))
	conn, err := a.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}
	a.conn = conn
	if err := a.subscribe(); err != nil {
		return fmt.Errorf("failed to subscribe to NATS server: %w", err)
	}
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

// connect is a private function used to connect the NATS client
func (a *NatsApplication) connect() (*nats.Conn, error) {
	timeout := a.ConnectionTimeout
	// Default timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	// Connect to the NATS server
	options := nats.GetDefaultOptions()
	options.Url = a.NatsServerUrl
	options.Timeout = timeout
	return options.Connect()
}

// subscribe is a private function used to subscribe to NATS
func (a *NatsApplication) subscribe() error {
	// Define the callback function
	cb := func(msg micro.Request) {
		a.logger.Info("Received a message", zap.String("subject", msg.Subject()), zap.String("data", string(msg.Data())))
		msg.Respond([]byte("Hello!"))
	}
	// Add the service
	service, err := micro.AddService(a.conn, micro.Config{
		Name:        "demo",
		Version:     "0.1.0",
		Description: "A demo micro-service",
	})
	if err != nil {
		return fmt.Errorf("failed to add service: %w", err)
	}
	// Add the endpoint
	if err := service.AddEndpoint("foo", HandlerFunc(cb)); err != nil {
		return fmt.Errorf("failed to add endpoint: %w", err)
	}
	return nil
}

type HandlerFunc func(micro.Request)

func (h HandlerFunc) Handle(msg micro.Request) {
	h(msg)
}

// Typeguards
var (
	_ caddy.App         = &NatsApplication{}
	_ caddy.Provisioner = &NatsApplication{}
)
