package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
)

// RPCClientConfig configures downstream service discovery or direct addressing.
type RPCClientConfig struct {
	GreeterService    string `json:"greeterService"`
	GreeterAddress    string `json:"greeterAddress"`
	CalculatorService string `json:"calculatorService"`
	CalculatorAddress string `json:"calculatorAddress"`
}

// Init loads the supplied configuration as the process-wide GoFrame config.
func Init(fileName string) error {
	if !gfile.Exists(fileName) {
		return fmt.Errorf("config file does not exist: %s", fileName)
	}

	adapter, err := gcfg.NewAdapterFile(fileName)
	if err != nil {
		return fmt.Errorf("create config adapter: %w", err)
	}
	g.Cfg().SetAdapter(adapter)

	ctx := gctx.GetInitCtx()
	serverAddress, err := g.Cfg().Get(ctx, "server.address")
	if err != nil {
		return fmt.Errorf("read config server.address: %w", err)
	}
	if strings.TrimSpace(serverAddress.String()) == "" {
		return fmt.Errorf("server.address is required")
	}

	rpcConfig, err := RPCClient(ctx)
	if err != nil {
		return err
	}
	if err = validateTarget("greeter", rpcConfig.GreeterService, rpcConfig.GreeterAddress); err != nil {
		return err
	}
	if err = validateTarget("calculator", rpcConfig.CalculatorService, rpcConfig.CalculatorAddress); err != nil {
		return err
	}
	return nil
}

func validateTarget(name, service, address string) error {
	hasService := strings.TrimSpace(service) != ""
	hasAddress := strings.TrimSpace(address) != ""
	if hasService == hasAddress {
		return fmt.Errorf(
			"exactly one of rpcclient.%sService or rpcclient.%sAddress is required",
			name,
			name,
		)
	}
	return nil
}

// RPCClient returns the gRPC client settings from the current GoFrame config.
func RPCClient(ctx context.Context) (RPCClientConfig, error) {
	var conf RPCClientConfig
	value, err := g.Cfg().Get(ctx, "rpcclient")
	if err != nil {
		return conf, fmt.Errorf("read config rpcclient: %w", err)
	}
	if err = gconv.Struct(value.Map(), &conf); err != nil {
		return conf, fmt.Errorf("parse config rpcclient: %w", err)
	}
	return conf, nil
}
