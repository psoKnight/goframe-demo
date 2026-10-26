package resilience

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ClientBreaker wraps one downstream gRPC service with a circuit breaker.
type ClientBreaker struct {
	name string
	cb   *gobreaker.CircuitBreaker[any]
}

// NewClientBreaker creates a GoFrame-configured circuit breaker.
func NewClientBreaker(ctx context.Context, configKey, name string) (*ClientBreaker, error) {
	conf, err := LoadCircuitBreakerConfig(ctx, configKey)
	if err != nil {
		return nil, err
	}
	if !conf.Enabled {
		return nil, nil
	}
	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return nil, err
	}
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return nil, err
	}

	return &ClientBreaker{
		name: name,
		cb: gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
			Name:        name,
			MaxRequests: conf.HalfOpenMaxRequests,
			Interval:    interval,
			Timeout:     timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= conf.FailureThreshold
			},
			IsExcluded:    isBusinessError,
			OnStateChange: logStateChange,
		}),
	}, nil
}

// UnaryClientInterceptor executes the downstream RPC through the breaker.
func (b *ClientBreaker) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		_, err := b.cb.Execute(func() (any, error) {
			return nil, invoker(ctx, method, req, reply, cc, opts...)
		})
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return gerror.NewCode(
				gcode.New(int(codes.Unavailable), "", nil),
				fmt.Sprintf("%s service is temporarily unavailable", b.name),
			)
		}
		return err
	}
}

func isBusinessError(err error) bool {
	if err == nil {
		return false
	}
	switch grpcCode(err) {
	case codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.OutOfRange,
		codes.FailedPrecondition,
		codes.ResourceExhausted:
		return true
	default:
		return false
	}
}

func grpcCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if code := gerror.Code(err); code.Code() != gcode.CodeNil.Code() {
		return codes.Code(code.Code())
	}
	if grpcStatus, ok := status.FromError(err); ok {
		return grpcStatus.Code()
	}
	return codes.Unknown
}

func logStateChange(name string, from, to gobreaker.State) {
	g.Log().Warningf(gctx.GetInitCtx(), "circuit breaker state changed: %s %s -> %s", name, from, to)
}
