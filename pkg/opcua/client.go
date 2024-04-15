package opcua

import (
	"context"
	"fmt"
)

/*
ReadVariableRequest is a struct used to request a variable read.
*/
type ReadVariableRequest struct {
	NodeID string
}

/*
ReadVariableResponse is a struct returned by the ReadVariable usecase.
*/
type ReadVariableResponse struct {
	Value any
}

/*
IReadVariableUseCase is an interface for the ReadVariable usecase.
*/
type IReadVariableUseCase interface {
	Execute(ctx context.Context, request *ReadVariableRequest) (*ReadVariableResponse, error)
}

func NewReadVariableUsecase(reader OpcuaReader) *ReadVariableUseCase {
	return &ReadVariableUseCase{
		client: reader,
	}
}

func NewReadVariableRequest(nodeid string) *ReadVariableRequest {
	return &ReadVariableRequest{
		NodeID: nodeid,
	}
}

/*
ReadVariableUseCase is a usecase for reading a variable from an OPC UA server.
*/
type ReadVariableUseCase struct {
	client OpcuaReader
}

/*
ReadVariable reads a variable from an OPC UA server.
*/
func (r *ReadVariableUseCase) Execute(ctx context.Context, request *ReadVariableRequest) (*ReadVariableResponse, error) {

	response, err := r.client.Read(ctx, request.NodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable: %w", err)
	}
	return response, nil
}

// Typeguards
var _ IReadVariableUseCase = &ReadVariableUseCase{}
