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

## 🔧 配置

编辑 `internal/config/config.go` 修改配置：

```go
Config{
    Server: ServerConfig{
        Port:    ":3128",          # 服务器端口
        Env:     "development",    # 运行环境
        Timeout: 30,
    },
    File: FileConfig{
        UploadDir:  "./files",     # 上传目录
        MaxSize:    100,           # 最大文件大小(MB)
    },
}
```

## 🔐 API密钥配置

系统提供两个默认API密钥供开发测试使用：
- `demo-key-12345` (30天有效期)
- `test-key-67890` (7天有效期)

使用密钥调用受保护的API：

```bash
# Header方式
curl -H "X-API-Key: demo-key-12345" http://localhost:3128/api/v1/images/upload

# Query参数方式
curl "http://localhost:3128/api/v1/images/upload?api_key=demo-key-12345"
```

## 📋 性能指标

| 指标 | 值 |
|------|-----|
| 请求限流 | 100请求/秒 |
| 并发连接数 | 10/IP |
| 缓存有效期 | 5分钟 |
| 最大上传文件 | 100MB |
| 支持的图片格式 | 8种 (jpg, png, gif, webp等) |
| 日志级别 | 5级 (DEBUG/INFO/WARN/ERROR/FATAL) |

## 🧪 测试API

使用VS Code REST Client扩展，打开 `api/api.http` 测试所有端点。

## 📝 许可证

MIT

## 👨‍💻 开发者

GantoHo
