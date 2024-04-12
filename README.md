# OPC-UA PLaygroud

## Objectives

1. Learn how to use the [Go language](https://go.dev/)
2. Learn how to use the [gopcua](https://github.com/gopcua/opcua) library
3. Learn how to use [nats.go](https://github.com/nats-io/nats.go) library

## Steps

1. Run a basic OPC-UA server

2. Run a basic OPC-UA client that will read a value from the server

3. Run a basic OPC-UA client that will write a value to the server

## Install

1. [Install go](https://askubuntu.com/a/1222190/863359):

```bash
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```

2. Build the server:

```bash
go build ./cmd/server
```

3. Build the client:

```bash
go build ./cmd/client
```

## Run

1. Run the server:

```bash
./server
```

2. Run the client:

```bash
./client
```
