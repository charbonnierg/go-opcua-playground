package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/gopcua/opcua/debug"

	"playground/pkg/opcua"
)

/*
This example shows how to read a variable value from an OPC-UA server
using its node ID.

For demonstration purpose, all operations are performed using a timeout context.
*/
func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://localhost:4840", "OPC UA Endpoint URL")
		nodeID   = flag.String("node", "ns=1;i=101", "NodeID to read")
	)
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	// Create root context for the application
	rootCtx := context.Background()

	// Create a new GopcuaReader instance
	reader, err := opcua.NewGopcuaReader(*endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Connect the GopcuaReader to the server
	connectCtx, cancel := context.WithTimeout(rootCtx, 5*time.Second)
	defer cancel()
	if err := reader.Connect(connectCtx); err != nil {
		log.Fatal(err)
	}

	// Close the GopcuaReader on exit
	defer func() {
		closeCtx, cancel := context.WithTimeout(rootCtx, 5*time.Second)
		defer cancel()
		reader.Close(closeCtx)
	}()

	// Sleep 10 seconds just to make sure that connect context is ONLY used for connection
	// and can be cancelled before reading
	<-time.After(10 * time.Second)

	// Create a new usecase
	usecase := opcua.NewReadVariableUsecase(reader)
	// Create a new context with a timeout
	usecaseCtx, cancel := context.WithTimeout(rootCtx, 1*time.Second)
	defer cancel()
	// Read with a timeout of 1 second
	request := opcua.NewReadVariableRequest(*nodeID)
	response, err := usecase.Execute(usecaseCtx, request)
	// Handle errors
	if err != nil {
		log.Fatal(err)
	}
	// Print the response
	log.Printf("%#v", response.Value)
}
