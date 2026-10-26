package resilience

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gcache"
)

func TestRateLimiterAllowsBurstThenRejects(t *testing.T) {
	limiter := &RateLimiter{
		name:  "test",
		limit: 0.001,
		burst: 1,
		ttl:   defaultRateLimitKeyTTL,
		cache: newTestCache(),
	}
	if !limiter.allow(context.Background(), "client") {
		t.Fatal("first request should be allowed")
	}
	if limiter.allow(context.Background(), "client") {
		t.Fatal("second request should be rejected")
	}
	if !limiter.allow(context.Background(), "another-client") {
		t.Fatal("different key should have an independent bucket")
	}
}

func newTestCache() *gcache.Cache {
	return gcache.New()
}
