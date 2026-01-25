# Go Image System

一个使用Go语言和Gin框架开发的高效图片管理系统。遵循Go官方推荐的标准项目结构。

## 🚀 快速开始

### 编译
```bash
go build -o build/image-sys.exe ./cmd/image-sys
```

### 运行
```bash
./build/image-sys.exe
```

### 使用构建脚本
```bash
# Windows
.\scripts\build.bat build
.\scripts\build.bat run

# Linux/Mac
./scripts/build.sh run

# 或使用Makefile
cd scripts && make run
```

### Docker运行
```bash
docker-compose -f deployments/docker-compose.yml up
```

## 📚 文档

所有文档位于 `docs/` 目录：
- [README.md](docs/README.md) - 完整项目文档
- [QUICKSTART.md](docs/QUICKSTART.md) - 快速开始指南
- [QUICK_START_GUIDE.md](docs/QUICK_START_GUIDE.md) - 详细API使用指南
- [FEATURES_SUMMARY.md](docs/FEATURES_SUMMARY.md) - 功能完善总结
- [IMPROVEMENTS.md](docs/IMPROVEMENTS.md) - 改进详情
- [CLEANUP_REPORT.md](docs/CLEANUP_REPORT.md) - 清理报告

## 📁 项目结构

按照Go官方推荐的标准结构组织 (golang-standards/project-layout)：

```
├── api/                    # API文档和规范
├── build/                  # 编译输出 (gitignored)
├── cmd/
│   └── image-sys/         # 可执行程序入口
│       └── main.go        # 程序主入口
├── configs/               # 配置文件
├── deployments/           # Docker和部署配置
│   ├── Dockerfile
│   └── docker-compose.yml
├── docs/                  # 文档
├── internal/              # 私有包 (不对外导出)
│   ├── app/              # 应用核心
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP处理器
│   ├── middleware/       # 中间件
│   ├── router/           # 路由
│   └── service/          # 业务逻辑
├── pkg/                   # 公共包 (可被导入)
│   ├── errors/           # 错误定义
│   ├── logger/           # 日志系统
│   └── utils/            # 工具函数
├── scripts/              # 构建脚本
│   ├── build.bat
│   ├── build.sh
│   └── Makefile
├── tests/                # 集成测试
├── files/                # 上传文件存储
├── logs/                 # 日志文件 (gitignored)
├── .air.toml             # 热加载配置
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🔗 API端点

### 新API (推荐)
```
GET  /api/v1/health              # 健康检查
GET  /api/v1/images              # 获取所有图片
GET  /api/v1/images/metadata     # 获取图片元数据
GET  /api/v1/images/paginated    # 分页查询图片
GET  /api/v1/images/search       # 搜索/过滤图片
GET  /api/v1/images/random       # 获取随机图片
GET  /api/v1/images/random/:num  # 获取N个随机图片
POST /api/v1/images/upload       # 上传图片 (需密钥)
DELETE /api/v1/images/:filename  # 删除图片 (需密钥)
POST /api/v1/images/delete       # 批量删除 (需密钥)
POST /api/v1/admin/api-keys      # 创建新密钥 (需认证)
GET  /api/v1/admin/api-keys      # 查看密钥 (需认证)
DELETE /api/v1/admin/api-keys    # 撤销密钥 (需认证)
GET  /f/:filename                # 直接获取文件
```

### 遗留API (向后兼容)
```
GET  /v1/                   # 健康检查
GET  /v1/all                # 获取所有图片
GET  /v1/bgimg              # 获取随机图片
GET  /v1/get/:number        # 获取N个随机图片
POST /v1/upload             # 上传图片 (需密钥)
```

## ✨ 项目特性

- ✅ 标准Go项目结构（遵循golang-standards）
- ✅ 清晰的分层架构 (Config/Service/Handler)
- ✅ 完整的错误处理机制
- ✅ 完善的日志系统（DEBUG/INFO/WARN/ERROR/FATAL）
- ✅ 跨域资源共享 (CORS) 支持
- ✅ Docker容器化支持
- ✅ 热加载开发支持
- ✅ 向后兼容的API端点
- ✅ **图片格式验证** (jpg, png, gif, webp, bmp, ico, svg)
- ✅ **分页查询** (支持自定义页数和大小)
- ✅ **图片搜索和过滤** (按名称、大小、类型)
- ✅ **上传进度跟踪**
- ✅ **内存缓存机制** (TTL + 自动清理)
- ✅ **分级日志系统** (可配置日志级别)
- ✅ **统一API响应格式** (含元数据和耗时)
- ✅ **请求限流** (100请求/秒，10并发/IP)
- ✅ **图片元数据返回** (大小、MIME类型、修改时间)
- ✅ **批量删除功能**
- ✅ **API密钥认证系统** (SHA256加密、过期机制)

# go-img-sys — Go 图片管理系统

一个用 Go (Gin) 编写的轻量级图片管理与分发服务，提供上传、查询、导出、清理、APIKey 管理等功能，适合小到中型图片存储场景和自托管部署。

本 README 覆盖项目结构、构建运行、配置项、API 列表、认证机制、部署与常见故障排查（含 Apifox 使用提示）。

--

## 主要特性

- 多种图片读取接口：列表 / 分页 / 搜索 / 随机值
- 文件上传（支持多文件、进度与冲突策略）
- 批量删除、按条件清理与导出为 ZIP
- 基于 API Key 的访问控制（SHA256 存储 + 过期/撤销）
- 日志系统（文件与 stdout，分级）
- 内置限流、中间件（CORS、耗时统计、速率限制）
- Docker 化与 docker-compose 支持

--

## 项目布局（重要路径）

- `cmd/image-sys` — 可执行入口（`main.go`）
- `internal/app` — 应用启动、配置与依赖初始化
- `internal/router` — 路由注册（API 路径在这里定义）
- `internal/handler` — HTTP 处理器（核心业务接口实现）
- `internal/service` — 业务逻辑实现（文件管理、导出等）
- `internal/config` — 默认配置（端口、上传目录、重复文件策略）
- `internal/middleware` — 认证、限流、CORS、计时等中间件
- `pkg/auth` — API Key 管理（生成/校验/默认 key）
- `pkg/logger` — 日志初始化与封装
- `pkg/utils` — 常用工具（路径、文件判断、统一响应）
- `deployments/` — `Dockerfile` 与 `docker-compose.yml`
- `api/api.http` — REST 测试示例（REST Client / curl 参考）

上传文件存储目录：默认 `./files`（可在配置中修改）。日志目录：`./logs`。

--

## 快速开始（本地）

1. 克隆仓库并进入项目根目录

2. 编译（Windows 示例）：

```powershell
go build -o build/image-sys.exe ./cmd/image-sys
```

3. 运行：

```powershell
./build/image-sys.exe
```

4. 默认服务监听：`:3128`，控制台会打印初始化信息（包括默认开发用 API keys）。

脚本支持：`.\scripts\build.bat`（Windows）或 `./scripts/build.sh`（Unix）。

--

## 使用 Docker

构建镜像并启动（在项目根目录）：

```bash
docker-compose -f deployments/docker-compose.yml up --build -d
```

这会把 `./files` 与 `./logs` 映射到容器内 `/root/files` 和 `/root/logs`。

--

## 配置说明

默认配置位于 `internal/config/config.go`，主要字段：

- `Server.Port`：服务端口（默认 `:3128`）
- `Server.Env`：运行环境（`development` / `release`）
- `File.UploadDir`：上传目录（默认 `./files`）
- `File.MaxSize`：单文件最大大小（MB，默认 `100`）
- `File.AllowTypes`：允许的 MIME 类型列表
- `File.DuplicateStrategy`：文件重名处理策略（`overwrite`、`rename`、`reject`；默认 `rename`）

示例（修改 `internal/config/config.go` 后重启生效）：

```go
AppConfig.File.DuplicateStrategy = "overwrite"
```

DuplicateStrategy 行为说明：

- `rename`（默认）：若存在则生成 `name_1.ext`、`name_2.ext`... 直到找到未被占用的名称。
- `overwrite`：直接覆盖已存在文件。
- `reject`：返回失败并在响应里列出被拒绝的文件。

--

## API 概览（重要端点）

注意：受保护的写操作需要带 API Key（`X-API-Key` header 或 `api_key` query）。

常用端点：

- GET  `/api/v1/health` — 健康检查
- GET  `/api/v1/images` — 列表所有图片（返回带 URL 的数据）
- GET  `/api/v1/images/metadata` — 返回包含元数据的列表
- GET  `/api/v1/images/paginated` — 分页查询（`page` / `page_size`）
- GET  `/api/v1/images/search` — 按名称/大小/类型搜索（支持 `filename`, `min_size`, `max_size`, `type` 等查询）
- GET  `/api/v1/images/random` — 随机图片（文本返回文件名或 URL）
- GET  `/api/v1/images/random/:number` — 获取 N 个随机图片（最大 100）
- POST `/api/v1/images/upload` — 上传（multipart/form-data，字段名 `files`，受保护）
- DELETE `/api/v1/images/:filename` — 删除单个文件（受保护）
- POST `/api/v1/images/delete` — 批量删除（JSON body: { "filenames": [...] }，受保护）

管理类（API Key 管理，受保护）：

- POST `/api/v1/admin/api-keys` — 创建 API Key（body: {"expire_days": <int>}）
- GET  `/api/v1/admin/api-keys` — 列出 Key 信息（不返回明文）
- DELETE `/api/v1/admin/api-keys` — 撤销 Key（body: {"api_key": "<plain>"}）

直接文件访问：

- GET `/f/:filename` — 直接从 `UploadDir` 返回文件。

兼容旧路径（向后兼容）：`/v1/*` 系列接口也存在以支持历史客户端。

更多请求示例见 `api/api.http`。

--

## 认证（API Key）

实现概要：

- 客户端需要在请求中提供明文 API Key（Header `X-API-Key` 或 query `api_key`）。
- 服务端使用 SHA256 对明文 Key 散列后与内存中存储的哈希值比对（见 `pkg/auth/keymanager.go`）。
- Key 支持过期与撤销。
- 启动时会初始化两个开发用默认 Key（仅用于快速本地测试）：
    - `demo-key-12345`（30 天）
    - `test-key-67890`（7 天）

示例（curl）：

```bash
# 使用 Header
curl -H "X-API-Key: demo-key-12345" -F "files=@/path/to/img.jpg" http://localhost:3128/api/v1/images/upload

# 使用 query（不推荐在浏览器地址栏暴露）
curl -F "files=@/path/to/img.jpg" "http://localhost:3128/api/v1/images/upload?api_key=demo-key-12345"
```

Apifox / Postman 使用提示：

- 在 Apifox 中不要直接把 `{{apiKey}}` 放到 header 而不在环境中定义。请在 Apifox 的 Environment 中创建 `apiKey` 并填入实际值（例如 `demo-key-12345`），或者直接在请求头里使用明文值测试。

--

## 文件上传行为

- HTTP 表单字段名：`files`（支持多文件）
- 限制：单文件大小受 `File.MaxSize` 控制（单位 MB）
- 重名冲突由 `DuplicateStrategy` 控制（见上文）
- 成功响应包含已上传文件的 `filename`, `size`, `url` 等信息；失败文件会被列在 `failed` 字段中

示例响应结构（成功/部分失败）：

```json
{
    "message": "Upload completed",
    "total_files": 2,
    "total_uploaded": 1,
    "uploaded": [{"filename":"a.jpg","url":"localhost:3128/f/a.jpg","size":12345}],
    "failed": [{"filename":"b.jpg","error":"file already exists"}]
}
```

--

## 日志与监控

- 日志目录：`./logs`，错误会写入 `logs/error.log`，普通信息输出到 stdout。
- 日志分级：`DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`（可通过 `pkg/logger` 调整）

--

## 部署与运维注意

- 确保 `files` 与 `logs` 目录对运行用户可写。
- 生产环境应替换默认 API Key 并移除开发演示 Key。
- 建议在反向代理（如 Nginx）前端做 TLS 终止并限制请求体大小。

--

## 常见问题（FAQ）

- Apifox 返回 401 且你确认服务端有默认 Key：通常是因为在 Apifox 请求头里使用了未定义的变量（例如 `{{apiKey}}` 未在 Environment 中设置）。解决方案：在 Apifox 环境中添加 `apiKey` 变量或直接填入明文 Key。
- Windows PowerShell 的 `curl` 可能映射到 `Invoke-WebRequest`，请使用 `curl.exe` 或在 Git Bash / WSL 中使用原生 curl。

--

## 开发者与贡献

欢迎提交 Issue 与 PR。主要代码在 `internal/` 下，公共库在 `pkg/` 下。提交前请保证：

- 按需运行 `go fmt` 与 `go vet`。
- 新功能加入对应单元/集成测试（若可能）。

--

## 许可证

本项目仓库当前未包含 LICENSE 文件。请在需要开源许可时添加合适的 `LICENSE`。

--

如果你希望我把 README 中的某一部分扩展为示例脚本（例如详细的上传 curl 测试、Apifox 环境导入片段或 Docker 部署步骤），告诉我需要的部分，我会补充示例并把它加入到仓库中。
