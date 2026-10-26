package config

import (
	"context"
	"testing"
)

func TestConfigFiles(t *testing.T) {
	for _, file := range []string{"config_local.yml", "config_test.yml", "config.yml"} {
		t.Run(file, func(t *testing.T) {
			if err := Init(file); err != nil {
				t.Fatalf("Init(%q) error = %v", file, err)
			}
			conf, err := RPCClient(context.Background())
			if err != nil {
				t.Fatalf("RPCClient() error = %v", err)
			}
			if conf.GreeterService == "" || conf.CalculatorService == "" {
				t.Fatalf("RPCClient() = %+v, want service discovery names", conf)
			}
		})
	}
}
