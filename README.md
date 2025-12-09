# doeot-go HTTP 版本示例（无 RPC）

这是一个只包含 HTTP 接口的最小可运行示例，用来先把新架构跑起来：

- DDD 分层：`internal/domain` / `internal/app` / `internal/infra` / `internal/interfaces/http`
- handler：单接口单文件，继承 `Get/Post/Put/Delete` 约束
- 依赖注入：`internal/di/container.go`（基于 dig）
- 数据库：MySQL + GORM
- 暂时 **不包含 JSON-RPC 相关代码**，后续再单独加回去

## 准备 MySQL

可以用 Docker 快速起一个本地 MySQL：

```bash
docker run -d --name doeot-mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=doeot \
  -p 3306:3306 \
  mysql:8.0
```

数据库默认配置：

- 数据库名：`doeot`
- 用户：`root`
- 密码：`123456`

可以使用 `schema/sql/orders.sql` 手动建表，也可以直接依赖 GORM 的 AutoMigrate（第一次运行会自动建表）。

## 配置

通过环境变量控制：

- `HTTP_PORT`：HTTP 端口，默认 `8080`
- `MYSQL_DSN`：MySQL DSN，默认：

```text
root:123456@tcp(127.0.0.1:3306)/doeot?charset=utf8mb4&parseTime=True&loc=Local
```

## 安装依赖

在项目根目录执行：

```bash
go mod tidy
```

## 启动 HTTP 服务

```bash
go run ./cmd/http-api
```

## HTTP CRUD 示例

### 创建订单

```bash
curl -X POST http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"u1","amount":99.9}'
```

### 获取订单

```bash
curl http://localhost:8080/orders/1
```

### 列出订单（可选按 user_id 过滤）

```bash
curl "http://localhost:8080/orders?user_id=u1"
```

### 更新订单

```bash
curl -X PUT http://localhost:8080/orders/1 \
  -H 'Content-Type: application/json' \
  -d '{"amount":120.5,"status":"paid"}'
```

### 删除订单

```bash
curl -X DELETE http://localhost:8080/orders/1
```
