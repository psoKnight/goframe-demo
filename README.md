# goframe-demo

参考 `grpc-demo` 的调用链，使用 GoFrame v2.10 实现的 HTTP + gRPC 微服务示例。

## 架构

```text
HTTP request
  -> generalgateway (GoFrame ghttp)
  -> generalgateway/service
  -> generalgateway/grpc_client (GoFrame grpcx)
     -> greeter       /v1/hello
     -> calculator    /v1/calculate
```

三个服务都使用 GoFrame 生态能力：

- `ghttp` 负责 HTTP 路由、参数解析、OpenAPI 与优雅启停。
- `grpcx` 负责 gRPC 服务端、客户端、服务注册、服务发现、负载均衡、日志、错误和链路追踪拦截器。
- `gcfg` 统一加载 YAML 配置，启动时仍支持参考项目的 `-c` 参数。
- 本地开发默认使用 GoFrame 文件注册中心；配置 `registry.etcd.address` 后切换为 Etcd，支持容器动态扩容和多主机部署。

## 目录

```text
.
├── pbdata/greeterpb/       # protobuf 协议与生成代码
├── pbdata/calculatorpb/    # Calculator protobuf 协议与生成代码
├── common/healthing/       # HTTP 健康检查处理器
├── common/registry/        # 文件注册中心 / Etcd 注册中心初始化
├── greeter/                # gRPC 业务服务
├── calculator/             # 整数计算 gRPC 服务
├── generalgateway/         # HTTP API 网关
└── docker-compose.yml      # 容器环境示例
```

## 本地运行

要求 Go 1.25+、GoFrame CLI、Protocol Buffers 编译器，以及 `protoc-gen-go`、`protoc-gen-go-grpc`。

先启动 gRPC 服务：

```bash
go run ./greeter -c=greeter/config/config_local.yml
```

另一个终端启动 Calculator：

```bash
go run ./calculator -c=calculator/config/config_local.yml
```

再启动 HTTP 网关：

```bash
go run ./generalgateway -c=generalgateway/config/config_local.yml
```

调用：

```bash
curl -X POST http://localhost:8080/v1/hello \
  -H "Content-Type: application/json" \
  -d '{"name":"Codex"}'
```

响应：

```json
{"error_no":0,"error_msg":"成功","data":{"message":"Hello, Codex! (Greeter)"}}
```

计算接口：

```bash
curl -X POST http://localhost:8080/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","left":21,"right":21}'
```

```json
{"error_no":0,"error_msg":"成功","data":{"result":42}}
```

`operation` 支持 `add` 和 `multiply`。`left`、`right` 支持 JSON 数字或字符串，大整数建议使用字符串避免客户端精度损失。超出 `int64` 范围时返回业务参数错误。

健康检查：

```bash
curl http://localhost:50052/status.html
curl http://localhost:50054/status.html
curl http://localhost:8080/status.html
```

GoFrame OpenAPI 与 Swagger UI：

```text
http://localhost:8080/api.json
http://localhost:8080/swagger
```

## 服务发现与直连

本地 `config_local.yml` 使用 GoFrame 文件注册中心，网关按服务名发现本地进程：

```yaml
rpcclient:
  greeterService: "greeter"
  greeterAddress: ""
  calculatorService: "calculator"
  calculatorAddress: ""
```

测试与容器配置通过 Etcd 共享服务发现信息，客户端仍使用服务名：

```yaml
rpcclient:
  greeterService: "greeter"
  greeterAddress: ""
  calculatorService: "calculator"
  calculatorAddress: ""

registry:
  etcd:
    address: "etcd:2379"
    keepaliveTTL: "10s"
    dialTimeout: "5s"
```

每个下游服务必须且只能配置服务名或直连地址中的一个。

`registry.etcd.address` 为空时使用进程本地的 GoFrame 文件注册中心；非空时自动切换为 Etcd。服务通过 Etcd Lease 保活，进程异常退出后实例会在 TTL 到期后自动下线。

配置文件用途：

| 文件 | 用途 |
| --- | --- |
| `config_local.yml` | 本机开发，使用文件注册中心 |
| `config_test.yml` | 测试/Compose，使用 Etcd 服务名 `etcd:2379` |
| `config.yml` | 正式环境，已启用 Etcd，默认占位地址为 `etcd.example.internal:2379` |

正式发布前必须把三个 `config.yml` 中的 Etcd 占位地址替换为实际集群地址。

## 构建与发布

只生成 Linux 发布目录：

```bash
cd greeter
bash deploy.sh -b test

cd ../calculator
bash deploy.sh -b test

cd ../generalgateway
bash deploy.sh -b test
```

发布目录位于各服务的 `docker/release`，包含 Linux 二进制、启动脚本和配置。

构建并推送镜像：

```bash
cd greeter
IMAGE_REPO=harbor.example.com/goframe-demo bash deploy.sh -d 1.0.0 test

cd ../calculator
IMAGE_REPO=harbor.example.com/goframe-demo bash deploy.sh -d 1.0.0 test

cd ../generalgateway
IMAGE_REPO=harbor.example.com/goframe-demo bash deploy.sh -d 1.0.0 test
```

## Docker Compose

先为三个服务生成发布目录，再启动：

```bash
(cd greeter && bash deploy.sh -b test)
(cd calculator && bash deploy.sh -b test)
(cd generalgateway && bash deploy.sh -b test)
docker compose up --build -d
```

Compose 内置单节点 Etcd。Greeter 和 Calculator 不占用宿主机固定端口，所有请求统一从 Gateway 的 `8080` 进入，因此可以动态扩缩容：

```bash
docker compose up -d \
  --scale greeter=3 \
  --scale calculator=2
```

扩容或缩容时 GoFrame 会通过 Etcd Watch 自动更新客户端端点，默认使用轮询负载均衡：

```bash
docker compose up -d --scale greeter=5 --scale calculator=3

docker compose exec etcd etcdctl get --prefix /service
docker compose logs -f generalgateway
```

如需要多个 Gateway，请为每个 Gateway 配置不同的 `server.address`，并在它们前面增加反向代理或负载均衡器。生产环境建议使用高可用 Etcd 集群，而不是 Compose 中用于演示的单节点 Etcd。
```

## 修改协议

```bash
bash pbdata/greeterpb/pb.sh
bash pbdata/calculatorpb/pb.sh
go test ./...
go vet ./...
```

当前示例使用明文 gRPC，生产环境应结合 TLS、身份认证、超时及高可用 Etcd 集群。
# goframe-demo
