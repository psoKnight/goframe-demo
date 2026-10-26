package service

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"goframe-demo/pbdata/calculatorpb"
)

type fakeCalculatorClient struct {
	addCalled      bool
	multiplyCalled bool
	request        *calculatorpb.BinaryOperationRequest
}

func (f *fakeCalculatorClient) Add(
	_ context.Context,
	req *calculatorpb.BinaryOperationRequest,
	_ ...grpc.CallOption,
) (*calculatorpb.OperationReply, error) {
	f.addCalled = true
	f.request = req
	return &calculatorpb.OperationReply{Result: req.GetLeft() + req.GetRight()}, nil
}

func (f *fakeCalculatorClient) Multiply(
	_ context.Context,
	req *calculatorpb.BinaryOperationRequest,
	_ ...grpc.CallOption,
) (*calculatorpb.OperationReply, error) {
	f.multiplyCalled = true
	f.request = req
	return &calculatorpb.OperationReply{Result: req.GetLeft() * req.GetRight()}, nil
}

func TestCalculateDispatchesOperation(t *testing.T) {
	tests := []struct {
		name         string
		operation    string
		left         int64
		right        int64
		want         int64
		wantAdd      bool
		wantMultiply bool
	}{
		{name: "add", operation: OperationAdd, left: 21, right: 21, want: 42, wantAdd: true},
		{name: "multiply", operation: OperationMultiply, left: 21, right: 2, want: 42, wantMultiply: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &fakeCalculatorClient{}
			reply, err := calculate(context.Background(), client, test.operation, test.left, test.right)
			if err != nil {
				t.Fatalf("calculate() error = %v", err)
			}
			if reply.GetResult() != test.want {
				t.Fatalf("calculate() = %d, want %d", reply.GetResult(), test.want)
			}
			if client.addCalled != test.wantAdd || client.multiplyCalled != test.wantMultiply {
				t.Fatalf("dispatched add=%v multiply=%v", client.addCalled, client.multiplyCalled)
			}
			if got := client.request; got.GetLeft() != test.left || got.GetRight() != test.right {
				t.Fatalf("request = %+v, want left=%d right=%d", got, test.left, test.right)
			}
		})
	}
}

func TestCalculateRejectsUnknownOperation(t *testing.T) {
	client := &fakeCalculatorClient{}
	_, err := calculate(context.Background(), client, "divide", 21, 2)
	if got := gerror.Code(err).Code(); got != int(codes.InvalidArgument) {
		t.Fatalf("calculate() code = %d, want %d", got, int(codes.InvalidArgument))
	}
	if client.addCalled || client.multiplyCalled {
		t.Fatal("unknown operation must not call downstream")
	}
}
