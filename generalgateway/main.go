package main

import (
	"flag"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"goframe-demo/common/healthing"
	"goframe-demo/common/registry"
	"goframe-demo/common/resilience"
	"goframe-demo/generalgateway/config"
	"goframe-demo/generalgateway/controller"
	"goframe-demo/generalgateway/grpc_client"
)

func main() {
	configFile := flag.String("c", "generalgateway/config/config_local.yml", "configure file path")
	flag.Parse()

	ctx := gctx.GetInitCtx()
	if err := config.Init(*configFile); err != nil {
		g.Log().Fatalf(ctx, "init configuration: %v", err)
	}
	if err := registry.Init(ctx); err != nil {
		g.Log().Fatalf(ctx, "init service registry: %v", err)
	}

	rpcConfig, err := config.RPCClient(ctx)
	if err != nil {
		g.Log().Fatalf(ctx, "load rpc client configuration: %v", err)
	}
	if err = grpc_client.Init(ctx, rpcConfig); err != nil {
		g.Log().Fatalf(ctx, "init gRPC client: %v", err)
	}
	defer grpc_client.Close()

	rateLimiter, err := resilience.NewRateLimiter(ctx, "resilience.httpRateLimit", "gateway-http")
	if err != nil {
		g.Log().Fatalf(ctx, "init HTTP rate limiter: %v", err)
	}

	server := g.Server("generalgateway")
	if err = server.SetConfigWithMap(g.Cfg().MustGet(ctx, "server").Map()); err != nil {
		g.Log().Fatalf(ctx, "configure HTTP server: %v", err)
	}
	server.BindHandler("GET:/status.html", healthing.Status)
	server.Group("/v1", func(group *ghttp.RouterGroup) {
		if rateLimiter != nil {
			group.Middleware(rateLimiter.HTTPMiddleware)
		}
		group.Middleware(controller.Response)
		group.Bind(controller.New())
	})
	server.Run()
}
