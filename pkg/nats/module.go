package nats

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(NatsApplication{})
}

// CaddyModule returns the Caddy module information.
// It implements the caddy.Module interface.
func (NatsApplication) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "nats",
		New: func() caddy.Module { return new(NatsApplication) },
	}
}
