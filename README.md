# Wallet Service

基于 [go-zero](https://go-zero.dev) 微服务框架构建的钱包服务，采用 API + RPC 分层架构。支持钱包创建、查询、充值和转账功能。

## 项目结构

```
wallet/
├── service/wallet/
│   ├── api/                          # HTTP API 网关层
│   │   ├── wallet.go                 # API 服务入口
│   │   ├── wallet.api                # go-zero API 定义文件
│   │   ├── doc.go                    # Swagger 静态文件嵌入
│   │   ├── doc/                      # Swagger UI 静态资源
│   │   │   ├── index.html
│   │   │   ├── wallet.json           # OpenAPI 2.0 规范文件
│   │   │   └── swagger-ui-*.js/css
│   │   ├── etc/
│   │   │   └── wallet-api.yaml       # API 服务配置
│   │   └── internal/
│   │       ├── config/config.go      # 配置结构体
│   │       ├── handler/              # HTTP 路由处理器 (goctl 生成)
│   │       ├── logic/                # 业务逻辑层 (调用 RPC)
│   │       ├── svc/servicecontext.go # 服务上下文 (持有 RPC 客户端)
│   │       └── types/types.go        # 请求/响应类型 (goctl 生成)
│   │
│   └── rpc/                          # gRPC 服务层
│       ├── wallet.go                 # RPC 服务入口
│       ├── wallet.proto              # Protobuf 定义文件
│       ├── pb/                       # protoc 生成的 Go 代码
│       ├── wallet/                   # goctl 生成的 RPC 客户端
│       ├── etc/
│       │   └── wallet.yaml           # RPC 服务配置
│       └── internal/
│           ├── config/config.go
│           ├── logic/                # 核心业务逻辑 (钱包操作)
│           ├── server/               # gRPC server 实现 (goctl 生成)
│           └── svc/servicecontext.go # 服务上下文 (内存存储)
│
├── deploy/docker-compose/            # Docker Compose 专用配置
│   ├── wallet-api.yaml
│   └── wallet-rpc.yaml
├── jmeter/                           # JMeter 测试脚本
│   ├── wallet-api-test.jmx           # 功能测试
│   └── wallet-api-stress.jmx         # 压力测试
├── docker-compose.yaml
├── Dockerfile.api
├── Dockerfile.rpc
├── build.sh                          # Docker 镜像构建脚本 (Linux/macOS)
├── start.bat                         # 本地启动脚本 (Windows)
└── stop.bat                          # 本地停止脚本 (Windows)
```

**架构概览：**

```
                HTTP                    gRPC
  Client  ──────────▶  wallet-api  ──────────▶  wallet-rpc
                        :8888                     :8080
                     (网关/路由)               (业务逻辑/存储)
```

- wallet-api：HTTP 网关，负责接收外部请求、参数校验，通过 gRPC 调用 wallet-rpc
- wallet-rpc：核心业务服务，负责钱包的创建、查询、充值、转账等操作，数据存储在内存中（`sync.Map`）

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/wallets` | 创建钱包 |
| GET | `/wallets/:walletId` | 查询钱包 |
| POST | `/wallets/deposit` | 充值 |
| POST | `/wallets/transfer` | 转账 |

## 启动方法

### 前置条件

- Go 1.25+
- Docker & Docker Compose（Docker 方式启动时需要）

### 方式一：脚本启动（Windows）

直接运行 `start.bat`，脚本会自动编译并启动 RPC 和 API 两个服务：

```bat
start.bat
```

脚本执行流程：
1. 编译 wallet-rpc 和 wallet-api 到 `bin/` 目录
2. 先启动 RPC 服务（端口 8080）
3. 等待 3 秒后启动 API 服务（端口 8888）

启动后：
- API 地址：http://localhost:8888
- Swagger UI：http://localhost:8888/swagger/index.html

停止服务：

```bat
stop.bat
```

### 方式二：Docker Compose 启动

**1. 构建镜像**

```bash
# 构建所有镜像
bash build.sh

# 或单独构建
bash build.sh rpc   # 仅构建 wallet-rpc
bash build.sh api   # 仅构建 wallet-api

# 自定义镜像 tag
IMAGE_TAG=v1.0.0 bash build.sh
```

**2. 启动服务**

```bash
docker compose up -d
```

**3. 查看日志**

```bash
docker compose logs -f
```

**4. 停止服务**

```bash
docker compose down
```

## 可扩展性

### API 服务水平扩展

wallet-api 是无状态的 HTTP 网关，所有业务状态都在 wallet-rpc 中维护。因此 API 层可以部署多个实例，通过负载均衡分发流量。

**Docker Compose 多实例部署示例：**

```yaml
services:
  wallet-rpc:
    container_name: wallet-rpc
    image: wallet-rpc:latest
    ports:
      - "8080:8080"
    volumes:
      - ./deploy/docker-compose/wallet-rpc.yaml:/app/etc/wallet.yaml
    restart: unless-stopped

  wallet-api-1:
    image: wallet-api:latest
    ports:
      - "8881:8888"
    volumes:
      - ./deploy/docker-compose/wallet-api.yaml:/app/etc/wallet-api.yaml
    depends_on:
      - wallet-rpc
    restart: unless-stopped

  wallet-api-2:
    image: wallet-api:latest
    ports:
      - "8882:8888"
    volumes:
      - ./deploy/docker-compose/wallet-api.yaml:/app/etc/wallet-api.yaml
    depends_on:
      - wallet-rpc
    restart: unless-stopped

  wallet-api-3:
    image: wallet-api:latest
    ports:
      - "8883:8888"
    volumes:
      - ./deploy/docker-compose/wallet-api.yaml:/app/etc/wallet-api.yaml
    depends_on:
      - wallet-rpc
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "8888:8888"
    volumes:
      - ./deploy/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - wallet-api-1
      - wallet-api-2
      - wallet-api-3
    restart: unless-stopped
```

对应的 Nginx 负载均衡配置（`deploy/nginx.conf`）：

```nginx
events {
    worker_connections 1024;
}

http {
    upstream wallet_api {
        server wallet-api-1:8881;
        server wallet-api-2:8882;
        server wallet-api-3:8883;
    }

    server {
        listen 8888;

        location / {
            proxy_pass http://wallet_api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

**扩展架构：**

```
                          ┌──────────────┐
                     ┌───▶│ wallet-api-1 │───┐
                     │    └──────────────┘   │
  Client ──▶ Nginx ──┤    ┌──────────────┐   ├──▶ wallet-rpc
              :8888   ├───▶│ wallet-api-2 │───┤      :8080
                     │    └──────────────┘   │
                     └───▶│ wallet-api-3 │───┘
                          └──────────────┘
```

go-zero 框架内置了 gRPC 客户端负载均衡，多个 API 实例连接同一个 RPC 服务时会自动处理连接复用和负载分配。

## Swagger 文档

服务启动后，访问 Swagger UI 查看和测试 API：

```
http://localhost:8888/swagger/index.html
```

Swagger UI 通过 Go 的 `embed` 机制嵌入到 API 二进制文件中，无需额外部署静态文件服务。

API 规范文件位于 `service/wallet/api/doc/wallet.json`（OpenAPI 2.0 格式），包含所有接口的请求参数和响应结构定义。

## JMeter 测试

项目提供了两个 JMeter 测试脚本，位于 `jmeter/` 目录：

| 文件 | 用途 |
|------|------|
| `wallet-api-test.jmx` | 功能测试 — 验证各接口的正确性 |
| `wallet-api-stress.jmx` | 压力测试 — 评估服务在高并发下的表现 |

### 运行方式

**前置条件：** 安装 [Apache JMeter](https://jmeter.apache.org/)（5.0+）

**GUI 模式（调试/查看结果）：**

```bash
jmeter -t jmeter/wallet-api-test.jmx
jmeter -t jmeter/wallet-api-stress.jmx
```

**命令行模式（推荐用于压测）：**

```bash
# 功能测试
jmeter -n -t jmeter/wallet-api-test.jmx -l result-test.jtl -e -o report-test/

# 压力测试
jmeter -n -t jmeter/wallet-api-stress.jmx -l result-stress.jtl -e -o report-stress/
```

参数说明：
- `-n`：非 GUI 模式运行
- `-t`：指定测试脚本
- `-l`：输出结果文件
- `-e -o`：生成 HTML 报告到指定目录

> 运行测试前请确保服务已启动（API 地址默认为 `http://localhost:8888`）。
