# generalgateway

GoFrame `ghttp` API 网关，通过 `grpcx` 调用 `greeter` 服务。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/v1/hello` | 调用 Greeter 服务 |
| `POST` | `/v1/calculate` | 调用 Calculator 服务 |
| `GET` | `/status.html` | 健康检查 |
| `GET` | `/api.json` | OpenAPI 文档 |
| `GET` | `/swagger` | Swagger UI |

请求示例：

```bash
curl -X POST http://localhost:8080/v1/hello \
  -H "Content-Type: application/json" \
  -d '{"name":"Codex"}'
```

计算请求：

```bash
curl -X POST http://localhost:8080/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"multiply","left":"6","right":"7"}'
```

本地配置通过 GoFrame 文件注册中心发现 `greeter` 和 `calculator`。测试与容器配置通过 Etcd 注册中心使用相同服务名进行发现。每个服务的服务名字段与地址字段必须且只能配置一个。

本地启动：

```bash
go run ./generalgateway -c=generalgateway/config/config_local.yml
```

`config_local.yml` 使用文件注册中心，`config_test.yml` 使用 Compose 内的 Etcd，正式 `config.yml` 已启用 Etcd，发布前需修改其中的占位地址。

调用链：

```text
POST /v1/hello     -> grpc_client.GreeterClient    -> greeter:50051
POST /v1/calculate -> grpc_client.CalculatorClient -> calculator:50053
```

构建：

```bash
cd generalgateway
bash deploy.sh -b test
IMAGE_REPO=harbor.example.com/goframe-demo bash deploy.sh -d 1.0.0 test
```
