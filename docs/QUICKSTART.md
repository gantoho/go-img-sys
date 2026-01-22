# 🚀 快速开始指南

## 项目已成功改进并重新规划！

您的Go图片系统已按照Go社区最佳实践进行了完整的改进和重新组织。

## ✨ 主要改进

### 项目结构
```
go-img-sys/
├── cmd/                      # CLI 入口
├── internal/                 # 私有内部包
│   ├── app/                 # 应用核心
│   ├── config/              # 配置管理
│   ├── handler/             # HTTP处理器
│   ├── middleware/          # 中间件
│   ├── router/              # 路由定义
│   └── service/             # 业务逻辑
├── pkg/                      # 公共可复用包
│   ├── errors/              # 错误处理
│   ├── logger/              # 日志系统
│   └── utils/               # 工具函数
└── files/                    # 上传文件目录
```

### 新增功能
- ✅ 集中式配置管理
- ✅ 完整的日志系统（含日志文件）
- ✅ 统一的错误处理
- ✅ 工具库（文件、响应）
- ✅ 分层架构（Service、Handler分离）
- ✅ 改进的CORS中间件
- ✅ 新的API路由 (/api/v1/*)
- ✅ 向后兼容旧API (/v1/*)
- ✅ Docker和Makefile支持
- ✅ 跨平台构建脚本

## 📦 快速启动

### 1️⃣ 第一次运行（仅需一次）

```bash
# 确保Go版本 >= 1.22
go version

# 下载依赖
go mod tidy
```

### 2️⃣ 编译项目

#### Windows
```bash
# 方式1：使用脚本
.\build.bat build

# 方式2：直接编译
go build -o image-sys.exe main.go

# 然后运行
.\image-sys.exe
```

#### Linux/Mac
```bash
# 方式1：使用脚本
./build.sh build
./build.sh run

# 方式2：使用Makefile
make build
make run

# 方式3：直接编译和运行
go run main.go
```

#### Docker
```bash
# 使用Docker Compose
docker-compose up

# 或自己构建
docker build -t image-sys:latest .
docker run -p 3128:3128 -v $(pwd)/files:/root/files image-sys:latest
```

### 3️⃣ 访问应用

打开浏览器访问：
```
http://localhost:3128/api/v1/health
```

应该看到：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "ok",
    "version": "1.0.0"
  }
}
```

## 🔗 常用API

### 获取所有图片
```bash
curl http://localhost:3128/api/v1/images
```

### 获取随机图片
```bash
# 获取1个
curl http://localhost:3128/api/v1/images/random

# 获取5个
curl http://localhost:3128/api/v1/images/random/5
```

### 上传图片
```bash
curl -X POST http://localhost:3128/api/v1/images/upload \
  -F "files=@/path/to/image.jpg"
```

### 直接访问图片
```
http://localhost:3128/f/image.jpg
```

### 健康检查
```bash
curl http://localhost:3128/api/v1/health
```

## 📝 文档位置

- **API文档**: [README_NEW.md](README_NEW.md)
- **改进详情**: [IMPROVEMENTS.md](IMPROVEMENTS.md)
- **API测试**: [api.http](api.http) （VS Code REST Client）
- **原始文档**: [README.md](README.md)

## 🛠️ 常用命令

### 使用 Makefile (推荐 Unix/Linux/Mac)
```bash
make help          # 显示帮助
make build         # 编译
make run           # 编译并运行
make clean         # 清理
make test          # 运行测试
make docker-build  # Docker构建
make docker-run    # Docker运行
```

### 使用脚本

#### Windows (.bat)
```bash
.\build.bat build   # 编译
.\build.bat run     # 运行
.\build.bat clean   # 清理
```

#### Linux/Mac (.sh)
```bash
./build.sh build    # 编译
./build.sh run      # 运行
./build.sh clean    # 清理
./build.sh test     # 测试
```

## 📂 关键文件说明

| 文件 | 说明 |
|------|------|
| `main.go` | 应用入口 |
| `internal/app/server.go` | 服务器初始化 |
| `internal/config/config.go` | 配置管理 |
| `internal/service/image_service.go` | 业务逻辑 |
| `internal/handler/image.go` | HTTP处理器 |
| `internal/router/router.go` | 路由定义 |
| `pkg/logger/logger.go` | 日志系统 |
| `pkg/errors/errors.go` | 错误定义 |
| `Dockerfile` | Docker镜像 |
| `Makefile` | 构建工具 |
| `api.http` | API测试文件 |

## 🔧 配置修改

编辑 `internal/config/config.go` 来修改配置：

```go
Config{
    Server: ServerConfig{
        Port:    ":3128",           // 改变端口
        Env:     "development",     // 改变环境
        Timeout: 30,
    },
    File: FileConfig{
        UploadDir:  "./files",      // 改变上传目录
        MaxSize:    100,            // 改变文件大小限制(MB)
        AllowTypes: []string{...},  // 改变允许的文件类型
    },
}
```

## 📊 项目统计

- Go文件数：15+
- 代码行数：2000+
- 支持的API端点：12+
- 文档页面：4+

## 🎯 下一步建议

1. **开发测试 API**
   ```bash
   # 使用 VS Code 扩展 "REST Client"
   # 打开 api.http 文件，点击 Send Request
   ```

2. **查看日志**
   ```bash
   # 应用运行时会在 logs/ 目录生成日志
   cat logs/error.log
   ```

3. **上传文件**
   ```bash
   # 在浏览器上传或使用 curl
   curl -X POST http://localhost:3128/api/v1/images/upload \
     -F "files=@test.jpg"
   ```

4. **生产部署**
   ```bash
   # 使用 Docker Compose
   docker-compose up -d
   ```

## ❓ 常见问题

**Q: 如何修改监听端口？**
A: 编辑 `internal/config/config.go` 中的 `Port: ":8080"`

**Q: 上传的文件存在哪里？**
A: 默认存在 `./files` 目录，可在配置中修改

**Q: 如何生成日志？**
A: 应用会自动在 `logs/` 目录生成日志文件

**Q: 旧的 API 还能用吗？**
A: 可以！所有 `/v1/*` 的端点都被保留支持

**Q: 如何进行开发测试？**
A: 使用 VS Code "REST Client" 扩展打开 `api.http` 文件

## 🚀 性能优化建议

1. 使用 Gin 的生产模式：设置环境变量 `GIN_MODE=release`
2. 启用 HTTP 缓存头
3. 添加请求速率限制
4. 定期清理旧日志
5. 使用 CDN 加速图片服务

## 🔐 安全建议

1. 验证上传的文件类型和大小
2. 实现身份验证和授权
3. 启用 HTTPS
4. 定期备份文件
5. 监控错误日志

## 📞 支持

如有问题，请查看：
1. [README_NEW.md](README_NEW.md) - 详细文档
2. [IMPROVEMENTS.md](IMPROVEMENTS.md) - 改进说明
3. [api.http](api.http) - API示例

---

**项目已准备好！开始享受您改进后的图片系统吧！** 🎉
