package resilience

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const defaultRateLimitKeyTTL = 10 * time.Minute

// RateLimitConfig configures a local token-bucket limiter.
type RateLimitConfig struct {
	Enabled bool    `json:"enabled"`
	Rate    float64 `json:"rate"`
	Burst   int     `json:"burst"`
	KeyTTL  string  `json:"keyTTL"`
}

// CircuitBreakerConfig configures an outbound gRPC circuit breaker.
type CircuitBreakerConfig struct {
	Enabled             bool   `json:"enabled"`
	FailureThreshold    uint32 `json:"failureThreshold"`
	Interval            string `json:"interval"`
	Timeout             string `json:"timeout"`
	HalfOpenMaxRequests uint32 `json:"halfOpenMaxRequests"`
}

// LoadRateLimitConfig loads and validates rate limiting settings.
func LoadRateLimitConfig(ctx context.Context, key string) (RateLimitConfig, error) {
	var conf RateLimitConfig
	value, err := g.Cfg().Get(ctx, key)
	if err != nil {
		return conf, fmt.Errorf("read config %s: %w", key, err)
	}
	if err = gconv.Struct(value.Map(), &conf); err != nil {
		return conf, fmt.Errorf("parse config %s: %w", key, err)
	}
	if !conf.Enabled {
		return conf, nil
	}
	if conf.Rate <= 0 {
		return conf, fmt.Errorf("%s.rate must be greater than 0", key)
	}
	if conf.Burst <= 0 {
		return conf, fmt.Errorf("%s.burst must be greater than 0", key)
	}
	if strings.TrimSpace(conf.KeyTTL) == "" {
		conf.KeyTTL = defaultRateLimitKeyTTL.String()
	}
	if _, err = time.ParseDuration(conf.KeyTTL); err != nil {
		return conf, fmt.Errorf("parse %s.keyTTL: %w", key, err)
	}
	return conf, nil
}

// LoadCircuitBreakerConfig loads and validates circuit breaker settings.
func LoadCircuitBreakerConfig(ctx context.Context, key string) (CircuitBreakerConfig, error) {
	var conf CircuitBreakerConfig
	value, err := g.Cfg().Get(ctx, key)
	if err != nil {
		return conf, fmt.Errorf("read config %s: %w", key, err)
	}
	if err = gconv.Struct(value.Map(), &conf); err != nil {
		return conf, fmt.Errorf("parse config %s: %w", key, err)
	}
	if !conf.Enabled {
		return conf, nil
	}
	if conf.FailureThreshold == 0 {
		return conf, fmt.Errorf("%s.failureThreshold must be greater than 0", key)
	}
	if conf.HalfOpenMaxRequests == 0 {
		conf.HalfOpenMaxRequests = 1
	}
	if _, err = time.ParseDuration(conf.Interval); err != nil {
		return conf, fmt.Errorf("parse %s.interval: %w", key, err)
	}
	if _, err = time.ParseDuration(conf.Timeout); err != nil {
		return conf, fmt.Errorf("parse %s.timeout: %w", key, err)
	}
	return conf, nil
}
