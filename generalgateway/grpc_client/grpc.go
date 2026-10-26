package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/grpc"

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

	greeterConn, err = grpcx.Client.NewGrpcClientConn(
		greeterTarget,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*20)),
	)
	if err != nil {
		return fmt.Errorf("create greeter connection: %w", err)
	}
	GreeterClient = greeterpb.NewGreeterClient(greeterConn)

	calculatorConn, err = grpcx.Client.NewGrpcClientConn(
		calculatorTarget,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*20)),
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
