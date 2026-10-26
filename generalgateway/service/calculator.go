package service

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc/codes"

	"goframe-demo/generalgateway/grpc_client"
	"goframe-demo/pbdata/calculatorpb"
)

const (
	OperationAdd      = "add"
	OperationMultiply = "multiply"
)

const calculatorCallTimeout = 3 * time.Second

// Calculate dispatches an add or multiply call to the Calculator service.
func Calculate(ctx context.Context, operation string, left, right int64) (*calculatorpb.OperationReply, error) {
	callCtx, cancel := context.WithTimeout(ctx, calculatorCallTimeout)
	defer cancel()

	return calculate(callCtx, grpc_client.CalculatorClient, operation, left, right)
}

func calculate(
	ctx context.Context,
	client calculatorpb.CalculatorClient,
	operation string,
	left, right int64,
) (*calculatorpb.OperationReply, error) {
	req := &calculatorpb.BinaryOperationRequest{Left: left, Right: right}
	switch operation {
	case OperationAdd:
		return client.Add(ctx, req)
	case OperationMultiply:
		return client.Multiply(ctx, req)
	default:
		return nil, gerror.NewCode(
			gcode.New(int(codes.InvalidArgument), "", nil),
			"unsupported operation",
		)
	}
}
