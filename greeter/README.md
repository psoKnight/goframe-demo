# greeter

GoFrame `grpcx` 业务服务，监听 gRPC 业务端口和独立 HTTP 健康检查端口。

| 端口 | 用途 |
| --- | --- |
| `50051` | gRPC `Greeter.SayHello` |
| `50052` | HTTP `GET /status.html` |

启动：

```bash
go run ./greeter -c=greeter/config/config_local.yml
```

本地使用 `config_local.yml` 和 GoFrame 文件注册中心；`config_test.yml` 使用 Compose 内的 Etcd；正式 `config.yml` 已启用 Etcd，发布前需修改其中的 Etcd 占位地址。网关按 `grpc.name` 服务名 `greeter` 发现多个实例。

服务通过 GoFrame `grpcx` 自定义 Unary Interceptor 实现每实例、每 RPC 方法的令牌桶限流，配置位于 `resilience.grpcRateLimit`。

服务同时注册标准 gRPC Health 服务和 gRPC Reflection，便于健康探测与调试。

构建：

```bash
cd greeter
bash deploy.sh -b test
IMAGE_REPO=harbor.example.com/goframe-demo bash deploy.sh -d 1.0.0 test
```
