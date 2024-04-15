package opcua

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

/*
OpcuaReader is an interface for an OPC UA reader.
*/
type OpcuaReader interface {
	Read(ctx context.Context, nodeID string) (*ReadVariableResponse, error)
}

/*
NewGopcuaReader creates a new GopcuaReader instance
*/
func NewGopcuaReader(endpoint string) (*GopcuaReader, error) {
	client, err := opcua.NewClient(endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &GopcuaReader{
		client: client,
	}, nil
}

/*
GopcuaReader is an implementation of OpcuaReader interface using gopcua library.
*/
type GopcuaReader struct {
	client *opcua.Client
}

func (g *GopcuaReader) Connect(ctx context.Context) error {
	return g.client.Connect(ctx)
}

func (g *GopcuaReader) Close(ctx context.Context) error {
	return g.client.Close(ctx)
}

func (g *GopcuaReader) Read(ctx context.Context, nodeID string) (*ReadVariableResponse, error) {
	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node id: %w", err)
	}
	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}
	var resp *ua.ReadResponse
	for {
		resp, err = g.client.Read(ctx, req)
		if err == nil {
			break
		}

		// Following switch contains known errors that can be retried by the user.
		// Best practice is to do it on read operations.
		switch {
		case err == io.EOF && g.client.State() != opcua.Closed:
			// has to be retried unless user closed the connection
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionIDInvalid):
			// Session is not activated has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionNotActivated):
			// Session is invalid has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSecureChannelIDInvalid):
			// secure channel will be recreated internally.
			time.After(1 * time.Second)
			continue

		default:
			return nil, fmt.Errorf("read failed: %w", err)
		}
	}
	status := resp.Results[0].Status
	if status != ua.StatusOK {
		return nil, fmt.Errorf("status not OK: %s", resp.Results[0].Status.Error())
	}
	value := resp.Results[0].Value.Value()
	return &ReadVariableResponse{Value: value}, nil
}

// Typeguard
var (
	_ OpcuaReader = &GopcuaReader{}
)
