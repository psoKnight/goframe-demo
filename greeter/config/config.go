package config

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

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
	for _, key := range []string{"grpc.name", "grpc.address", "healthserver.address"} {
		value, err := g.Cfg().Get(ctx, key)
		if err != nil {
			return fmt.Errorf("read config %s: %w", key, err)
		}
		if strings.TrimSpace(value.String()) == "" {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}
