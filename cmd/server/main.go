// Copyright 2024 Guillaume Charbonnier.
// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/server/attrs"
	"github.com/gopcua/opcua/ua"
)

/**
 * createServer creates a new test server with the given options.
 */
func createServer(
	listenAddress string,
	listenPort int,
	logWriter io.Writer,
) *server.Server {
	var opts []server.Option

	// Disable security for now.  This is the simplest way to get the server up and running.
	opts = append(opts,
		server.EnableSecurity("None", ua.MessageSecurityModeNone),
	)

	// Enable anonymous authentication.  This is the simplest way to get the server up and running.
	opts = append(opts,
		server.EnableAuthMode(ua.UserTokenTypeAnonymous),
	)

	// Set the server's hostname and port.  This is required for the server to know where to listen.
	opts = append(opts,
		server.EndPoint(listenAddress, listenPort),
	)

	// the server.SetLogger takes a server.Logger interface.  This interface is met by
	// the slog.Logger{}.  A simple wrapper could be made for other loggers if they don't already
	// meet the interface.
	logger := slog.New(slog.NewTextHandler(logWriter, nil))
	opts = append(opts,
		server.SetLogger(logger),
	)

	// Now that all the options are set, create the server.
	return server.New(opts...)
}

/**
 * initializeServer initializes the server with the default options.
 */
func initializeServer(ctx context.Context, instance *server.Server) (*server.Node, error) {
	// add the namespaces to the server, and add a reference to them if desired.
	// here we are choosing to add the namespaces to the root/object folder
	// to do this we first need to get the root namespace object folder so we
	// get the object node
	root_namespace, _ := instance.Namespace(0)
	root_node := root_namespace.Objects()
	// Start the server
	err := instance.Start(ctx)
	// check for errors
	if err != nil {
		return nil, err
	}
	// return the root node
	return root_node, nil
}

/**
 * initializeNamespace initializes the namespace with the default options.
 */
func initializeNamespace(name string, instance *server.Server, root *server.Node) (*server.NodeNameSpace, *server.Node) {
	// Now we'll add a node namespace.  This is a more traditional way to add nodes to the server
	// and is more in line with the opc ua node model, but may be more cumbersome for some use cases.
	// You can add namespaces before or after starting the server.
	node_namespace := server.NewNodeNameSpace(instance, name)

	// add the reference for this namespace's root object folder to the server's root object folder
	// @note: nns stands for "Node NameSpace"
	node_namespace_root := node_namespace.Objects()
	// @note: id.HasComponent is the type of reference that is being added.
	root.AddRef(node_namespace_root, id.HasComponent, true)
	// Return the namespace and the root node of the namespace
	return node_namespace, node_namespace_root
}

/**
 * addVariableNode adds a variable node to the namespace.
 */
func addVariableNode(ns *server.NodeNameSpace, ns_root *server.Node, name string, value any) *server.Node {
	// Create some nodes for it.  Here we are using the AddNewVariableNode utility function to create a new variable node
	variable := ns.AddNewVariableNode(name, value)
	// be sure to add the reference to the node to the node namespace, or clients won't be able to browse it.
	ns_root.AddRef(variable, id.HasComponent, true)
	return variable
}

/**
 * addVariableNodeAdvanced adds a variable node to the namespace with more control.
 */
func addVariableNodeAdvanced(
	ns *server.NodeNameSpace,
	ns_root *server.Node,
	nodeid uint32,
	name string,
	value any,
	class interface{},
) *server.Node {
	// Now we'll add a node from scratch:
	// build the node up with the correct attributes and references then reference it from
	// the parent node in the namespace if applicable.
	variable := server.NewNode(
		ua.NewNumericNodeID(ns.ID(), nodeid),
		map[ua.AttributeID]*ua.Variant{
			ua.AttributeIDBrowseName: ua.MustVariant(attrs.BrowseName(name)),
			ua.AttributeIDNodeClass:  ua.MustVariant(class),
		},
		nil,
		func() *ua.Variant { return ua.MustVariant(value) },
	)
	ns.AddNode(variable)
	ns_root.AddRef(variable, id.HasComponent, true)
	return variable
}

/**
 * main is the entry point for the server.
 */
func main() {
	endpoint := flag.String("endpoint", "0.0.0.0", "Address where OPC UA Endpoint should listen to")
	port := flag.Int("port", 4840, "Port which OPC UA Endpoint should listen to")
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	// Create a new server
	srv := createServer(*endpoint, *port, os.Stdout)

	// Start the server and get the root object node
	srv_root_node, err := initializeServer(context.Background(), srv)
	if err != nil {
		log.Fatalf("Error starting server, exiting: %s", err)
	}
	// Close the server on exit
	defer srv.Close()

	// Initialize the namespace
	ns, ns_root_node := initializeNamespace("MyNamespace", srv, srv_root_node)
	log.Printf("Namespace root node id is: %s", ns_root_node.ID())
	// First way to add a variable to the namespace
	var1 := addVariableNode(ns, ns_root_node, "Var1", 12.34)
	log.Printf("Variable 1 node id is: %s", var1.ID().String())

	// Another way to add a variable to the namespace
	var2 := addVariableNodeAdvanced(ns, ns_root_node, 12345, "Var2", 12.34, uint32(ua.NodeClassVariable))
	log.Printf("Variable 2 node id is: %s", var2.ID().String())

	// catch ctrl-c and gracefully shutdown the server.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	defer signal.Stop(sigch)
	log.Printf("Press CTRL-C to exit")

	<-sigch
	log.Printf("Shutting down the server...")
}
