package opcua

import (
	"context"
	"fmt"
)

type ReadVariableResponse struct {
	Value any
}

type IReadVariableUseCase interface {
	ReadVariable(ctx context.Context, nodeID string) (*ReadVariableResponse, error)
}

type IClient interface {
	Read(ctx context.Context, nodeID string) (*ReadVariableResponse, error)
}

type ReadVariableUseCase struct {
	client IClient
}

func (r *ReadVariableUseCase) ReadVariable(ctx context.Context, nodeID string) (*ReadVariableResponse, error) {

	response, err := r.client.Read(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable: %w", err)
	}
	return response, nil
}

var _ IReadVariableUseCase = &ReadVariableUseCase{}
