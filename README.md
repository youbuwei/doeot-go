# doeot-go (重构起步版)

这是一个简化后的 DDD + HTTP + JSON-RPC 示例服务，用于作为新架构的起点。

## 启动步骤

```bash
# 启动 HTTP 接口服务
go run ./cmd/http-api

# 启动 JSON-RPC 服务
go run ./cmd/json-rpc
```

默认端口：
- HTTP: `8080`
- JSON-RPC: `8090`
- SQLite 数据库文件：`data.db`（自动在当前目录创建）

你可以通过环境变量覆盖默认配置：

```bash
export HTTP_PORT=8080
export RPC_PORT=8090
export DB_PATH=./data.db
```

## 测试接口

### HTTP 创建订单

```bash
curl -X POST http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"u1","amount":99.9}'
```

返回示例：

```json
{"order_id":1}
```

### JSON-RPC 创建订单

```bash
curl -X POST http://localhost:8090/rpc \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","method":"order.create","params":{"user_id":"u1","amount":99.9},"id":1}'
```

返回示例：

```json
{"jsonrpc":"2.0","result":{"order_id":1},"id":1}
```
