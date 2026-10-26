package config

import (
	"context"
	"testing"

	"goframe-demo/common/resilience"
)

func TestConfigFiles(t *testing.T) {
	for _, file := range []string{"config_local.yml", "config_test.yml", "config.yml"} {
		t.Run(file, func(t *testing.T) {
			if err := Init(file); err != nil {
				t.Fatalf("Init(%q) error = %v", file, err)
			}
			conf, err := resilience.LoadRateLimitConfig(context.Background(), "resilience.grpcRateLimit")
			if err != nil {
				t.Fatalf("LoadRateLimitConfig() error = %v", err)
			}
			if !conf.Enabled {
				t.Fatal("grpc rate limit must be enabled")
			}
		})
	}
}
