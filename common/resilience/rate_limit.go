package resilience

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const rateLimitedResponse = `{"error_no":1429,"error_msg":"请求过于频繁，请稍后重试"}`

// RateLimiter provides per-key token buckets backed by GoFrame gcache.
type RateLimiter struct {
	name  string
	limit rate.Limit
	burst int
	ttl   time.Duration
	cache *gcache.Cache
}

// NewRateLimiter creates a limiter from a GoFrame configuration key.
func NewRateLimiter(ctx context.Context, configKey, name string) (*RateLimiter, error) {
	conf, err := LoadRateLimitConfig(ctx, configKey)
	if err != nil {
		return nil, err
	}
	if !conf.Enabled {
		return nil, nil
	}
	ttl, err := time.ParseDuration(conf.KeyTTL)
	if err != nil {
		return nil, fmt.Errorf("parse rate limiter key TTL: %w", err)
	}
	return &RateLimiter{
		name:  name,
		limit: rate.Limit(conf.Rate),
		burst: conf.Burst,
		ttl:   ttl,
		cache: gcache.New(),
	}, nil
}

// HTTPMiddleware limits requests by client IP, HTTP method and route.
func (l *RateLimiter) HTTPMiddleware(r *ghttp.Request) {
	key := fmt.Sprintf("http|%s|%s|%s", r.GetClientIp(), r.Method, r.URL.Path)
	if l.allow(r.Context(), key) {
		r.Middleware.Next()
		return
	}

	g.Log().Warningf(r.Context(), "rate limit rejected: limiter=%s key=%s", l.name, key)
	r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.Response.Header().Set("Retry-After", "1")
	r.Response.WriteStatus(http.StatusTooManyRequests, rateLimitedResponse)
	r.Exit()
}

// UnaryServerInterceptor applies a per-RPC-method limit on each service instance.
func (l *RateLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		key := "grpc|" + info.FullMethod
		if !l.allow(ctx, key) {
			g.Log().Warningf(ctx, "rate limit rejected: limiter=%s key=%s", l.name, key)
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(ctx, req)
	}
}

func (l *RateLimiter) allow(ctx context.Context, key string) bool {
	value, err := l.cache.GetOrSetFuncLock(ctx, key, func(context.Context) (any, error) {
		return rate.NewLimiter(l.limit, l.burst), nil
	}, l.ttl)
	if err != nil {
		g.Log().Errorf(ctx, "rate limiter cache failed, allowing request: %v", err)
		return true
	}
	limiter, ok := value.Val().(*rate.Limiter)
	if !ok {
		g.Log().Errorf(ctx, "rate limiter cache returned unexpected type: %T", value.Val())
		return true
	}
	return limiter.Allow()
}
