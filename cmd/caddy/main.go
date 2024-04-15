// SPDX-License-Identifier: Apache-2.0

package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	// plug in Caddy modules here
	_ "playground/pkg/nats"
)

func main() {
	caddycmd.Main()
}
