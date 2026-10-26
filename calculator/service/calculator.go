package service

import (
	"context"
	"math"
	"math/big"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc/codes"

	"goframe-demo/pbdata/calculatorpb"
)

// Server implements the Calculator gRPC business service.
type Server struct {
	calculatorpb.UnimplementedCalculatorServer
	Name string
}

// New creates a Calculator service with the supplied display name.
func New(name string) *Server {
	return &Server{Name: name}
}

// Add returns left + right and rejects int64 overflow.
func (s *Server) Add(_ context.Context, req *calculatorpb.BinaryOperationRequest) (*calculatorpb.OperationReply, error) {
	left, right := req.GetLeft(), req.GetRight()
	if (right > 0 && left > math.MaxInt64-right) ||
		(right < 0 && left < math.MinInt64-right) {
		return nil, overflowError()
	}

	return &calculatorpb.OperationReply{Result: left + right}, nil
}

// Multiply returns left * right and rejects int64 overflow.
func (s *Server) Multiply(_ context.Context, req *calculatorpb.BinaryOperationRequest) (*calculatorpb.OperationReply, error) {
	result := new(big.Int).Mul(big.NewInt(req.GetLeft()), big.NewInt(req.GetRight()))
	if !result.IsInt64() {
		return nil, overflowError()
	}
	return &calculatorpb.OperationReply{Result: result.Int64()}, nil
}

func overflowError() error {
	return gerror.NewCode(
		gcode.New(int(codes.OutOfRange), "", nil),
		"integer overflow",
	)
}
