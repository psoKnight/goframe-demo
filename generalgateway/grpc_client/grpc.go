package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/grpc"

	"goframe-demo/common/resilience"
	"goframe-demo/generalgateway/config"
	"goframe-demo/pbdata/calculatorpb"
	"goframe-demo/pbdata/greeterpb"
)

var (
	// GreeterClient and CalculatorClient are GoFrame service-discovery aware clients.
	GreeterClient    greeterpb.GreeterClient
	CalculatorClient calculatorpb.CalculatorClient
	greeterConn      *grpc.ClientConn
	calculatorConn   *grpc.ClientConn
)

// Init creates downstream gRPC clients from service names or direct addresses.
func Init(ctx context.Context, conf config.RPCClientConfig) error {
	greeterTarget, err := resolveTarget(conf.GreeterService, conf.GreeterAddress, "greeter")
	if err != nil {
		return err
	}
	calculatorTarget, err := resolveTarget(conf.CalculatorService, conf.CalculatorAddress, "calculator")
	if err != nil {
		return err
	}

	greeterBreaker, err := resilience.NewClientBreaker(ctx, "resilience.circuitBreaker.greeter", "greeter")
	if err != nil {
		return fmt.Errorf("init greeter circuit breaker: %w", err)
	}
	calculatorBreaker, err := resilience.NewClientBreaker(ctx, "resilience.circuitBreaker.calculator", "calculator")
	if err != nil {
		return fmt.Errorf("init calculator circuit breaker: %w", err)
	}

	greeterOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 20)),
	}
	if greeterBreaker != nil {
		greeterOptions = append(greeterOptions, grpcx.Client.ChainUnary(greeterBreaker.UnaryClientInterceptor()))
	}
	greeterConn, err = grpcx.Client.NewGrpcClientConn(
		greeterTarget,
		greeterOptions...,
	)
	if err != nil {
		return fmt.Errorf("create greeter connection: %w", err)
	}
	GreeterClient = greeterpb.NewGreeterClient(greeterConn)

	calculatorOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 20)),
	}
	if calculatorBreaker != nil {
		calculatorOptions = append(calculatorOptions, grpcx.Client.ChainUnary(calculatorBreaker.UnaryClientInterceptor()))
	}
	calculatorConn, err = grpcx.Client.NewGrpcClientConn(
		calculatorTarget,
		calculatorOptions...,
	)
	if err != nil {
		_ = greeterConn.Close()
		greeterConn = nil
		GreeterClient = nil
		return fmt.Errorf("create calculator connection: %w", err)
	}
	CalculatorClient = calculatorpb.NewCalculatorClient(calculatorConn)
	return nil
}

func resolveTarget(service, address, name string) (string, error) {
	if service = strings.TrimSpace(service); service != "" {
		return service, nil
	}
	if address = strings.TrimSpace(address); address != "" {
		return address, nil
	}
	return "", fmt.Errorf("%s service name or address is required", name)
}

// Close releases both downstream client connections.
func Close() error {
	var errs []error
	if greeterConn != nil {
		if err := greeterConn.Close(); err != nil {
			errs = append(errs, err)
		}
		greeterConn = nil
		GreeterClient = nil
	}
	if calculatorConn != nil {
		if err := calculatorConn.Close(); err != nil {
			errs = append(errs, err)
		}
		calculatorConn = nil
		CalculatorClient = nil
	}
	return errors.Join(errs...)
}
