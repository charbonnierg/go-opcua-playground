# Go Cheat Sheet

## Build

| Command | Description |
| --- | --- |
| `go build ./cmd/client` | Compile client executable |
| `go build ./cmd/server` | Compile server executable |

## Build & Run

| Command | Description |
| --- | --- |
| `go run ./cmd/client` | Compile and run client in one go 😊 |
| `go run ./cmd/server` | Compile and run server in one go 😊 |

## Format

| Command | Description |
| --- | --- |
| `go fmt ./...` | Format all files in the project |

## Test

| Command | Description |
| --- | --- |
| `go test ./...` | Run all tests in the project |

## Local Documentation

| Command | Description |
| --- | --- |
| `go install -v golang.org/x/tools/cmd/godoc@latest` | Install godoc tool |
| `godoc -http=:6060` | Run local documentation server |