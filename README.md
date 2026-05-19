# Go Image System

使用 Go 语言和 Gin 框架开发的高效图片管理系统，遵循 Go 官方推荐的标准项目结构。

## 特性

- **图片管理** — 上传、查询（列表/分页/搜索/随机）、删除（单张/批量）
- **图片处理** — 缩略图生成（Catmull-Rom 高质量缩放）、旋转、尺寸调整
- **文件导出** — 按需或全量导出为 ZIP，支持路径遍历防护
- **磁盘清理** — 孤立缩略图清理、过期文件删除、空目录清理
- **统计分析** — 文件数量/大小统计、格式分布、磁盘使用率
- **多层认证** — API Key（SHA256 + 过期/撤销） + JWT（可选）
- **并发安全** — 限流（令牌桶，100 req/s + 10并发/IP，支持环境变量配置）、并发安全的缓存
- **缓存机制** — 内存缓存（TTL + 自动清理 + 优雅关闭）+ fsnotify 文件监听自动失效
- **监控指标** — Prometheus 端点 `/metrics`（HTTP 请求数/延迟/并发、磁盘用量、文件数）
- **统一响应** — 标准 JSON 格式，含版本号、时间戳、请求耗时
- **日志系统** — 分级日志（DEBUG/INFO/WARN/ERROR/FATAL），lumberjack 自动轮转（按大小/时间）
- **断点续传** — 分片上传三步 API：init → upload chunk → complete
- **配置热加载** — POST `/api/v1/system/reload` 运行时重载配置
- **Docker 支持** — 多阶段构建 + 非 root 用户，docker-compose 一键部署

## 项目结构

```
go-img-sys/
├── api/                    # API 测试脚本
│   └── api.http
├── cmd/
│   └── image-sys/
│       └── main.go         # 程序入口
├── deployments/            # 部署配置
│   ├── Dockerfile
│   └── docker-compose.yml
├── docs/
│   └── CODE_REVIEW_OPTIMIZATION.md  # 代码审查与优化记录
├── internal/               # 私有包（不对外导出）
│   ├── app/server.go       # 应用启动与依赖组装
│   ├── config/config.go    # 配置管理
│   ├── handler/image.go    # HTTP 处理器
│   ├── middleware/          # 中间件（CORS/限流/JWT/认证/计时）
│   ├── router/router.go    # 路由注册
│   └── service/             # 业务逻辑层
│       ├── image_service.go
│       ├── export_service.go
│       ├── maintenance_service.go
│       └── statistics_service.go
├── pkg/                    # 公共库（可被外部导入）
│   ├── auth/               # API Key + JWT 认证
│   ├── cache/              # 内存缓存（TTL + 自动清理）
│   ├── errors/             # 统一错误定义
│   ├── imageutil/          # 图片处理（缩放/旋转/缩略图）
│   ├── logger/             # 分级日志
│   └── utils/              # 文件工具 + 统一响应
├── scripts/                # 构建脚本
│   ├── Makefile
│   ├── build.bat
│   └── build.sh
├── .air.toml               # 热加载配置
├── .gitignore
├── go.mod / go.sum
├── openapi.yaml            # OpenAPI 规范
└── README.md
```

## 快速开始

### 环境要求

- Go >= 1.23

> **国内网络用户**：如果遇到 `proxy.golang.org` 超时，请设置 Go 模块代理：
> ```powershell
> # PowerShell
> $env:GOPROXY="https://goproxy.cn,direct"
> # 或永久设置
> go env -w GOPROXY=https://goproxy.cn,direct
> ```

### 编译与运行

```bash
# 编译（产物在 build/ 目录）
go build -o build/image-sys ./cmd/image-sys/

# 运行
./build/image-sys

# 或使用 Makefile（自动生成 OpenAPI 文档）
cd scripts && make run

# 使用构建脚本
# Windows
.\scripts\build.bat run

# Linux/Mac
./scripts/build.sh run
```
#### 交叉编译
windows -> linux amd64
```bash
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o build/image-sys ./cmd/image-sys/
```

### 交叉编译

将项目编译为其他平台的可执行文件，方便部署到服务器：

```bash
cd scripts

# 编译为 Linux（最常用，适合部署到服务器）
make build-linux

# 编译为 Windows
make build-windows

# 编译为 macOS
make build-darwin

# 一次性编译所有平台
make build-all
```

编译产物统一输出到 `build/` 目录：

| 目标平台 | 输出文件 |
|---------|---------|
| 当前平台 | `build/image-sys` |
| Linux (amd64) | `build/image-sys-linux` |
| Windows (amd64) | `build/image-sys-windows.exe` |
| macOS (amd64) | `build/image-sys-darwin` |

> 所有二进制文件均输出到项目根目录的 `build/` 文件夹，便于管理。交叉编译需要 Go 标准库支持目标平台，无需额外安装工具链。

### 热加载开发

```bash
# 安装 air 后（make dev 会自动安装）
cd scripts && make dev
```

### Docker 运行

```bash
docker-compose -f deployments/docker-compose.yml up --build -d
```

## 配置

默认配置位于 `internal/config/config.go`，支持通过**环境变量**和**命令行参数**覆盖。

### 配置项

| 字段 | 环境变量 | 默认值 | 说明 |
|------|---------|--------|------|
| `Server.Port` | `SERVER_PORT` | `:3128` | 服务端口 |
| `Server.Env` | `SERVER_ENV` | `development` | 运行环境（development/release）|
| `Server.Timeout` | `SERVER_TIMEOUT` | `30` | 超时时间（秒） |
| `File.UploadDir` | `UPLOAD_DIR` | `./files` | 上传目录 |
| `File.MaxSize` | `MAX_FILE_SIZE` | `100` | 单文件最大大小（MB） |
| `File.DuplicateStrategy` | `DUPLICATE_STRATEGY` | `rename` | 重名策略 |
| `Auth.JWTSecret` | `JWT_SECRET` | `change-me-in-production` | JWT 签名密钥 |
| `Rate.RequestsPerSec` | `RATE_LIMIT_REQUESTS` | `100` | 每秒请求限制 |
| `Rate.ConcurrentLimit` | `RATE_LIMIT_CONCURRENT` | `10` | 每 IP 并发限制 |

### 命令行参数

```bash
# 指定端口
go run ./cmd/image-sys/ --port :8080

# 指定运行环境
go run ./cmd/image-sys/ --env release

# 组合使用
go run ./cmd/image-sys/ --port :8080 --env release
```

### 环境变量覆盖（最高优先级）

```bash
# Windows PowerShell
$env:SERVER_PORT=":8080"
$env:SERVER_ENV="release"
$env:JWT_SECRET="my-secure-key"
$env:MAX_FILE_SIZE="200"
go run ./cmd/image-sys/

# Linux/Mac
SERVER_PORT=:8080 SERVER_ENV=release JWT_SECRET=my-secure-key go run ./cmd/image-sys/
```

**优先级顺序**：硬编码默认值 < 环境变量 < 命令行参数

**重名策略（DuplicateStrategy）**：
- `rename`（默认）— 自动重命名为 `name_1.ext`、`name_2.ext` …
- `overwrite` — 直接覆盖
- `reject` — 拒绝并返回错误

## API 端点

### 公开接口（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health` | 健康检查（含磁盘空间检测） |
| GET | `/api/v1/images` | 图片列表（URL） |
| GET | `/api/v1/images/metadata` | 图片列表（含元数据） |
| GET | `/api/v1/images/paginated?page=1&page_size=20` | 分页查询 |
| GET | `/api/v1/images/search?filename=&min_size=&max_size=&type=` | 搜索/过滤 |
| GET | `/api/v1/images/random` | 随机一张（返回文件名） |
| GET | `/api/v1/images/random/:number` | 随机 N 张（最多 100，Fisher-Yates 无重复） |
| GET | `/api/v1/util/statistics` | 文件统计 |
| GET | `/api/v1/util/disk-usage` | 磁盘使用情况 |
| POST | `/api/v1/admin/api-keys` | 创建 API Key（无需认证） |
| POST | `/api/v1/admin/api-keys/validate` | 校验 API Key |
| GET | `/metrics` | Prometheus 监控指标 |
| GET | `/f/:filename` | 直接获取图片文件 |
| GET | `/bgimg` | 随机获取一张图片（重定向） |
| POST | `/upload` | 上传图片（multipart, 字段 `files`，无需认证） |

### 受保护接口（需 API Key）

在 Header 中传入 `X-API-Key` 或 Query 传入 `api_key`。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/images/upload` | 上传图片（multipart, 字段 `files`） |
| DELETE | `/api/v1/images/:filename` | 删除单张图片 |
| POST | `/api/v1/images/delete` | 批量删除图片 |
| POST | `/api/v1/images/upload/chunk/init` | 初始化分片上传 |
| POST | `/api/v1/images/upload/chunk` | 上传分片 |
| POST | `/api/v1/images/upload/chunk/complete` | 完成分片合并 |
| POST | `/api/v1/util/export` | 导出指定文件为 ZIP |
| POST | `/api/v1/util/export-all` | 全量导出为 ZIP |
| POST | `/api/v1/util/cleanup` | 磁盘清理（孤立缩略图/过期文件/空目录） |
| POST | `/api/v1/util/generate-thumbnails` | 生成缩略图（后台异步） |
| POST | `/api/v1/system/reload` | 配置热加载 |

### 管理接口（需 API Key）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/api-keys` | 列出所有 API Key 摘要信息 |
| DELETE | `/api/v1/admin/api-keys` | 撤销指定的 API Key |

### JWT 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 登录获取 JWT Token |
| POST | `/api/auth/refresh` | 刷新 JWT Token 有效期 |

### 遗留 API（向后兼容）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/v1/` | 健康检查 |
| GET | `/v1/all` | 所有图片 |
| GET | `/v1/bgimg` | 随机图片 |
| GET | `/v1/get/:number` | 随机 N 张 |
| POST | `/v1/upload` | 上传图片（需 API Key） |

更多请求示例见 [api/api.http](api/api.http)。

## 新增特性

### 📊 Prometheus 监控

启动服务后访问 `http://localhost:3128/metrics`，暴露指标：

| 指标 | 类型 | 说明 |
|------|------|------|
| `http_requests_total` | Counter | HTTP 请求总数（按 method/path/status 细分） |
| `http_request_duration_seconds` | Histogram | HTTP 请求延迟分布 |
| `http_requests_in_flight` | Gauge | 当前正在处理的请求数 |
| `image_upload_total` | Counter | 图片上传总数 |
| `image_delete_total` | Counter | 图片删除总数 |
| `disk_usage_percent` | Gauge | 当前磁盘使用率 |
| `file_count` | Gauge | 当前存储文件数 |

### 📦 断点续传（分片上传）

适用于大文件上传场景，三步完成：

```bash
# 1. 初始化上传会话
curl -X POST http://localhost:3128/api/v1/images/upload/chunk/init \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"filename":"large.mp4","file_size":"104857600"}'

# 2. 上传分片（重复调用，chunk_index 从 0 开始）
curl -X POST http://localhost:3128/api/v1/images/upload/chunk \
  -H "X-API-Key: demo-key-12345" \
  -F "upload_id=xxx" \
  -F "chunk_index=0" \
  -F "file=@chunk_0.bin"

# 3. 完成合并
curl -X POST http://localhost:3128/api/v1/images/upload/chunk/complete \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"upload_id":"xxx"}'
```

### ⚙️ 配置热加载

运行时重载环境变量配置，无需重启服务：

```bash
curl -X POST http://localhost:3128/api/v1/system/reload \
  -H "X-API-Key: demo-key-12345"
```

### 🖼️ 缩略图生成

为已有图片后台异步生成缩略图（200x200，Catmull-Rom 高质量缩放）：

```bash
curl -X POST "http://localhost:3128/api/v1/util/generate-thumbnails?filenames=image1.jpg,image2.png" \
  -H "X-API-Key: demo-key-12345"
```

缩略图存放在 `{uploadDir}/thumbs/` 目录下。

## 认证

### API Key

- Key 使用 SHA256 哈希存储
- 支持过期时间和手动撤销
- 服务启动时自动初始化默认 Key（仅开发用）：
  - `demo-key-12345`（30 天有效期）
  - `test-key-67890`（7 天有效期）

```bash
# 使用 Header 传 Key
curl -H "X-API-Key: demo-key-12345" \
  -F "files=@/path/to/img.jpg" \
  http://localhost:3128/api/v1/images/upload

# 或使用 Query 参数
curl -F "files=@/path/to/img.jpg" \
  "http://localhost:3128/api/v1/images/upload?api_key=demo-key-12345"
```

### JWT

支持通过用户名密码获取 JWT Token：

```bash
curl -X POST http://localhost:3128/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

默认账号：`admin/admin123`（管理员）、`user/user123`（普通用户）。

## 图片上传

- 表单字段名：`files`（支持多文件）
- 单文件大小受 `File.MaxSize` 限制
- 重名冲突处理由 `DuplicateStrategy` 控制

响应示例：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "Upload completed",
    "total_files": 2,
    "total_uploaded": 2,
    "uploaded": [
      {"index": 1, "filename": "a.jpg", "size": 12345, "url": "localhost:3128/f/a.jpg", "progress": 100}
    ]
  },
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1714512345,
    "duration_ms": 123
  }
}
```

## 图片处理

依赖 `golang.org/x/image` 实现高质量图片处理：

- **缩略图生成** — Catmull-Rom 插值算法，保持宽高比
- **图片旋转** — 支持 90°/180°/270° 像素级旋转
- **尺寸调整** — 指定宽高的高质量缩放

## 日志

- 标准输出：所有级别（可配置最低级别）
- 错误日志文件：`logs/error.log`（仅 ERROR / FATAL）
- 日志分级：DEBUG < INFO < WARN < ERROR < FATAL

通过 `logger.SetLogLevel()` 可动态调整日志级别。

## 部署

### 生产环境注意事项

1. **替换默认 API Key** — 在 `pkg/auth/keymanager.go` 中移除开发用 Key
2. **修改 JWT Secret** — 在 `internal/config/config.go` 中修改 `Auth.JWTSecret`
3. **使用反向代理** — 建议在 Nginx 前端做 TLS 终止
4. **目录权限** — 确保 `files/` 和 `logs/` 目录可写

### Docker 部署

```bash
# 构建并启动
docker-compose -f deployments/docker-compose.yml up --build -d

# 查看日志
docker-compose -f deployments/docker-compose.yml logs -f
```

## 开发

```bash
# 格式化
go fmt ./...

# 静态检查
go vet ./...

# 构建验证
go build ./...

# 运行
$env:GOPROXY="https://goproxy.cn,direct"  # 国内网络需要
go run ./cmd/image-sys/
```

## FAQ

**Q: 启动时报 `proxy.golang.org` 连接超时？**
A: 国内网络环境下，设置 Go 代理为 `https://goproxy.cn,direct`。

**Q: 上传返回 401 Unauthorized？**
A: 请求头中缺少 `X-API-Key` 或 Key 已过期。开发环境可用 `demo-key-12345`。

**Q: 构建时报 `main.go: no such file`？**
A: 请使用 `./cmd/image-sys/` 路径编译，而不是 `main.go`。

**Q: 如何修改监听端口？**
A: 修改 `internal/config/config.go` 中的 `Server.Port` 后重新编译。

## OpenAPI 规范与 Swagger UI

项目通过 **swaggo/swag** 注解自动生成 OpenAPI 规范，提供两个运行时端点：

| 端点 | 说明 |
|------|------|
| `/swagger/index.html` | Swagger UI 可视化浏览接口 |
| `/swagger/doc.json` | OpenAPI JSON 规范（供 Swagger UI 加载） |

### 访问 Swagger UI

启动服务后，在浏览器访问：

```
http://localhost:3128/swagger/index.html
```

### 生成方式

OpenAPI 规范通过代码注释自动生成，**新增或修改路由后需要重新生成**。

**触发方式（任选其一）：**
先安装swag
```bash
go install github.com/swaggo/swag/cmd/swag@latest

swag -v
```

```bash
# 方式 1：Makefile（推荐，会自动安装 swag CLI）
cd scripts && make build           # make build 已自动包含 openapi 生成

# 方式 2：go generate
go install github.com/swaggo/swag/cmd/swag@latest
go generate ./...

# 方式 3：直接运行 swag CLI
swag init -g cmd/image-sys/main.go -o docs --parseDependency --parseInternal
```

### 生成文件

执行生成后，`docs/` 目录下会输出三个文件：

| 文件 | 说明 |
|------|------|
| `docs/swagger.yaml` | OpenAPI 3.0 YAML 格式规范 |
| `docs/swagger.json` | OpenAPI 3.0 JSON 格式规范 |
| `docs/docs.go` | 嵌入二进制的 Go 代码（通过 `init()` 注册到 `swag` 包，运行时由 `/swagger/doc.json` 渲染） |

> 注意：`docs/docs.go` 需要在 `internal/router/router.go` 中通过 `_ "github.com/gantoho/go-img-sys/docs"` 导入，`init()` 函数才会执行并注册 API 定义。

### 注解位置

Swagger 注解写在 handler 函数上方，与代码逻辑在一起：

```
internal/handler/image.go  ← 所有 API 端点的注解
cmd/image-sys/main.go      ← 全局 API 信息（标题、版本、认证方式等）
```

### 认证方式声明

- **API Key** — 在 Swagger UI 中点击右上角 `Authorize`，输入 `X-API-Key`
- **JWT** — 先调用 `/api/auth/login` 获取 Token，然后在 Authorize 中输入 `Bearer <token>`

### 故障排查

| 问题 | 原因 | 解决 |
|------|------|------|
| `/swagger/doc.json` 返回 500 `no swag has yet been registered` | `docs` 包未导入，`init()` 未执行 | 确认 `internal/router/router.go` 中有 `_ "github.com/gantoho/go-img-sys/docs"` 导入 |
| `/swagger/index.html` 返回 `Not Found` | 静态文件读取失败 | 重新执行 `make openapi` 生成 `docs/` 文件后重新编译 |
| `docs/docs.go` 报编译错误 `unknown field LeftDelim` | swag CLI 与 `swaggo/swag` 库版本不匹配 | 执行 `go get github.com/swaggo/swag@latest && swag init ...` |

> **注意**：原有的 `openapi.yaml`（手动维护版本）已被自动生成的 `docs/swagger.yaml` 取代。请勿手动编辑 `openapi.yaml`，所有 API 变更应修改 handler 函数中的 Swagger 注解后重新生成。

## 相关文档

- [CODE_REVIEW_OPTIMIZATION.md](docs/CODE_REVIEW_OPTIMIZATION.md) — 代码审查与规范优化记录
- [api/api.http](api/api.http) — API 测试示例
