package main

import (
	"flag"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"goframe-demo/common/healthing"
	"goframe-demo/common/registry"
	"goframe-demo/common/resilience"
	"goframe-demo/greeter/config"
	"goframe-demo/greeter/service"
	"goframe-demo/pbdata/greeterpb"
)

func main() {
	configFile := flag.String("c", "greeter/config/config_local.yml", "configure file path")
	flag.Parse()

	ctx := gctx.GetInitCtx()
	if err := config.Init(*configFile); err != nil {
		g.Log().Fatalf(ctx, "init configuration: %v", err)
	}
	if err := registry.Init(ctx); err != nil {
		g.Log().Fatalf(ctx, "init service registry: %v", err)
	}

	grpcServerConfig := grpcx.Server.NewConfig()
	rateLimiter, err := resilience.NewRateLimiter(ctx, "resilience.grpcRateLimit", "greeter-grpc")
	if err != nil {
		g.Log().Fatalf(ctx, "init gRPC rate limiter: %v", err)
	}
	if rateLimiter != nil {
		grpcServerConfig.Options = append(
			grpcServerConfig.Options,
			grpc.ChainUnaryInterceptor(rateLimiter.UnaryServerInterceptor()),
		)
	}
	grpcServer := grpcx.Server.New(grpcServerConfig)
	greeterpb.RegisterGreeterServer(grpcServer.Server, service.New("Greeter"))

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer.Server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(grpcServer.GetConfig().Name, healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer.Server)

	httpServer := g.Server("greeter-health")
	httpServer.SetAddr(g.Cfg().MustGet(ctx, "healthserver.address").String())
	httpServer.SetDumpRouterMap(false)
	httpServer.BindHandler("GET:/status.html", healthing.Status)
	if err := httpServer.Start(); err != nil {
		g.Log().Fatalf(ctx, "start health server: %v", err)
	}

	grpcServer.Start()
	grpcServer.Wait()
	_ = httpServer.Shutdown()
}
