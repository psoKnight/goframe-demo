package resilience

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCircuitBreakerOpensAndRejects(t *testing.T) {
	breaker := testBreaker(1)
	calls := 0
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "downstream unavailable")
	}
	interceptor := breaker.UnaryClientInterceptor()

	_ = interceptor(context.Background(), "/test", nil, nil, nil, invoker)
	err := interceptor(context.Background(), "/test", nil, nil, nil, invoker)
	if got := gerror.Code(err).Code(); got != int(codes.Unavailable) {
		t.Fatalf("second call code = %d, want %d", got, int(codes.Unavailable))
	}
	if calls != 1 {
		t.Fatalf("downstream calls = %d, want 1", calls)
	}
}

func TestCircuitBreakerExcludesBusinessErrors(t *testing.T) {
	breaker := testBreaker(1)
	calls := 0
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.InvalidArgument, "invalid request")
	}
	interceptor := breaker.UnaryClientInterceptor()

	for range 3 {
		if err := interceptor(context.Background(), "/test", nil, nil, nil, invoker); err == nil {
			t.Fatal("expected business error")
		}
	}
	if calls != 3 {
		t.Fatalf("downstream calls = %d, want 3", calls)
	}
}

func testBreaker(threshold uint32) *ClientBreaker {
	return &ClientBreaker{
		name: "test",
		cb: gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
			Name: "test",
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= threshold
			},
			Timeout:    time.Second,
			IsExcluded: isBusinessError,
		}),
	}
}
