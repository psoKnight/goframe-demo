# calculator

GoFrame `grpcx` 整数计算微服务，提供加法和乘法，并检查 `int64` 溢出。

| RPC | 请求 | 结果 |
| --- | --- | --- |
| `Calculator.Add` | `{ left: int64, right: int64 }` | `{ result: int64 }` |
| `Calculator.Multiply` | `{ left: int64, right: int64 }` | `{ result: int64 }` |

| 端口 | 用途 |
| --- | --- |
| `50053` | gRPC 业务端口 |
| `50054` | HTTP `GET /status.html` |

启动：

```bash
go run ./calculator -c=calculator/config/config_local.yml
```

本地使用 `config_local.yml` 和 GoFrame 文件注册中心；`config_test.yml` 使用 Compose 内的 Etcd；正式 `config.yml` 已启用 Etcd，发布前需修改其中的 Etcd 占位地址。Gateway 按 `grpc.name` 服务名发现多个实例。服务同时注册标准 gRPC Health 服务和 gRPC Reflection。

通过网关调用：

```bash
curl -X POST http://localhost:8080/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"multiply","left":6,"right":7}'
```

结果超出 `int64` 范围时返回 gRPC `OutOfRange`。

构建：

```bash
cd calculator
bash deploy.sh -b test
IMAGE_REPO=harbor.example.com/goframe-demo bash deploy.sh -d 1.0.0 test
```
