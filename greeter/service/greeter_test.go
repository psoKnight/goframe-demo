package service

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc/codes"

	"goframe-demo/pbdata/greeterpb"
)

func TestServerSayHello(t *testing.T) {
	server := New("Greeter")

	reply, err := server.SayHello(context.Background(), &greeterpb.HelloRequest{Name: "GoFrame"})
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}
	if got, want := reply.GetMessage(), "Hello, GoFrame! (Greeter)"; got != want {
		t.Fatalf("SayHello() = %q, want %q", got, want)
	}
}

func TestServerSayHelloRequiresName(t *testing.T) {
	server := New("Greeter")

	_, err := server.SayHello(context.Background(), &greeterpb.HelloRequest{})
	if got := gerror.Code(err).Code(); got != int(codes.InvalidArgument) {
		t.Fatalf("SayHello() code = %v, want %v", got, int(codes.InvalidArgument))
	}
}
