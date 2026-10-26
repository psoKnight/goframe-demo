package service

import (
	"context"
	"time"

	"goframe-demo/generalgateway/grpc_client"
	"goframe-demo/pbdata/greeterpb"
)

const greeterCallTimeout = 3 * time.Second

// SayHello calls the Greeter gRPC service through its GoFrame client.
func SayHello(ctx context.Context, name string) (*greeterpb.HelloReply, error) {
	callCtx, cancel := context.WithTimeout(ctx, greeterCallTimeout)
	defer cancel()

	return grpc_client.GreeterClient.SayHello(callCtx, &greeterpb.HelloRequest{Name: name})
}
