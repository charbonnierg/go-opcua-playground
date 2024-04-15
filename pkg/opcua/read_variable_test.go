package opcua_test

import (
	"context"
	"fmt"
	"playground/pkg/opcua"
	"testing"
)

type StubOpcuaReader struct {
	value   any
	err     error
	nodeids []string
}

func (r *StubOpcuaReader) Read(ctx context.Context, nodeid string) (*opcua.ReadVariableResponse, error) {
	r.nodeids = append(r.nodeids, nodeid)
	if r.value == nil && r.err == nil {
		return nil, fmt.Errorf("stub value not set")
	}
	return &opcua.ReadVariableResponse{Value: r.value}, r.err
}

// SetValue sets the value that will be returned by the Read method
func (r *StubOpcuaReader) SetValue(value any) {
	r.value = value
}

// SetError sets the error that will be returned by the Read method
func (r *StubOpcuaReader) SetError(err error) {
	r.err = err
}

// Return the nodeids that were read
func (r *StubOpcuaReader) NodeIDs() []string {
	return r.nodeids
}

var (
	_ opcua.OpcuaReader = (*StubOpcuaReader)(nil)
)

func TestReadVariableUsecase(t *testing.T) {
	reader := &StubOpcuaReader{}
	usecase := opcua.NewReadVariableUsecase(reader)

	reader.SetValue(42)
	request := &opcua.ReadVariableRequest{NodeID: "ns=2;i=1234"}
	resp, err := usecase.Execute(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Value != 42 {
		t.Fatalf("unexpected value: %v", resp.Value)
	}
	if len(reader.NodeIDs()) != 1 {
		t.Fatalf("unexpected number of nodeids read: %v", len(reader.NodeIDs()))
	}
	if reader.NodeIDs()[0] != "ns=2;i=1234" {
		t.Fatalf("unexpected nodeid read: %v", reader.NodeIDs()[0])
	}
}
