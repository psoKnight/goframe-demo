package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc/codes"

	"goframe-demo/pbdata/greeterpb"
)

// Server implements the Greeter gRPC business service.
type Server struct {
	greeterpb.UnimplementedGreeterServer
	Name string
}

// New creates a Greeter service with the supplied display name.
func New(name string) *Server {
	return &Server{Name: name}
}

// SayHello validates the request and returns a greeting.
func (s *Server) SayHello(_ context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, gerror.NewCode(
			gcode.New(int(codes.InvalidArgument), "", nil),
			"name is required",
		)
	}

	return &greeterpb.HelloReply{
		Message: fmt.Sprintf("Hello, %s! (%s)", name, s.Name),
	}, nil
}
