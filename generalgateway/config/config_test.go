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
			conf, err := RPCClient(context.Background())
			if err != nil {
				t.Fatalf("RPCClient() error = %v", err)
			}
			if conf.GreeterService == "" || conf.CalculatorService == "" {
				t.Fatalf("RPCClient() = %+v, want service discovery names", conf)
			}
			rateLimit, err := resilience.LoadRateLimitConfig(context.Background(), "resilience.httpRateLimit")
			if err != nil {
				t.Fatalf("LoadRateLimitConfig() error = %v", err)
			}
			if !rateLimit.Enabled {
				t.Fatal("HTTP rate limit must be enabled")
			}
			for _, target := range []string{"greeter", "calculator"} {
				breaker, err := resilience.LoadCircuitBreakerConfig(
					context.Background(),
					"resilience.circuitBreaker."+target,
				)
				if err != nil {
					t.Fatalf("LoadCircuitBreakerConfig(%s) error = %v", target, err)
				}
				if !breaker.Enabled {
					t.Fatalf("%s circuit breaker must be enabled", target)
				}
			}
		})
	}
}
