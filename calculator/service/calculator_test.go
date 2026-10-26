package service

import (
	"context"
	"math"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc/codes"

	"goframe-demo/pbdata/calculatorpb"
)

func TestServerAdd(t *testing.T) {
	server := New("Calculator")

	reply, err := server.Add(context.Background(), &calculatorpb.BinaryOperationRequest{Left: 21, Right: 21})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if got, want := reply.GetResult(), int64(42); got != want {
		t.Fatalf("Add() = %d, want %d", got, want)
	}
}

func TestServerMultiply(t *testing.T) {
	server := New("Calculator")

	reply, err := server.Multiply(context.Background(), &calculatorpb.BinaryOperationRequest{Left: 6, Right: 7})
	if err != nil {
		t.Fatalf("Multiply() error = %v", err)
	}
	if got, want := reply.GetResult(), int64(42); got != want {
		t.Fatalf("Multiply() = %d, want %d", got, want)
	}
}

func TestServerRejectsOverflow(t *testing.T) {
	server := New("Calculator")
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "add positive",
			call: func() error {
				_, err := server.Add(context.Background(), &calculatorpb.BinaryOperationRequest{Left: math.MaxInt64, Right: 1})
				return err
			},
		},
		{
			name: "add negative",
			call: func() error {
				_, err := server.Add(context.Background(), &calculatorpb.BinaryOperationRequest{Left: math.MinInt64, Right: -1})
				return err
			},
		},
		{
			name: "multiply",
			call: func() error {
				_, err := server.Multiply(context.Background(), &calculatorpb.BinaryOperationRequest{Left: math.MaxInt64, Right: 2})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := gerror.Code(test.call()).Code(); got != int(codes.OutOfRange) {
				t.Fatalf("error code = %d, want %d", got, int(codes.OutOfRange))
			}
		})
	}
}
