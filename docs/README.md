# Go Image System - 图片系统

一个使用Go语言和Gin框架开发的高效图片管理系统。支持图片上传、查询、随机获取等功能。

## 项目结构

```
go-img-sys/
├── cmd/                      # 命令行入口
│   └── main.go
├── internal/                 # 内部包（私有包）
│   ├── config/              # 配置管理
│   │   └── config.go
│   ├── handler/             # HTTP处理器
│   │   └── image.go
│   ├── middleware/          # 中间件
│   │   └── cors.go
│   ├── router/              # 路由配置
│   │   └── router.go
│   ├── service/             # 业务逻辑层
│   │   └── image_service.go
│   └── server.go            # 服务器初始化
├── pkg/                      # 公共包（可被外部导入）
│   ├── errors/              # 错误定义
│   │   └── errors.go
│   ├── logger/              # 日志管理
│   │   └── logger.go
│   └── utils/               # 工具函数
│       ├── file.go          # 文件操作
│       └── response.go      # 响应封装
├── files/                    # 图片存储目录
├── logs/                     # 日志输出目录
├── go.mod
├── main.go                  # 应用入口
└── README.md
```

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 运行项目

```bash
# 方式1：直接运行
go run main.go

# 方式2：运行cmd下的main.go
go run cmd/main.go
```

项目启动后，访问 `http://localhost:3128`

## API 文档

### 健康检查

```http
GET /api/v1/health
```

**响应:**
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

### 获取所有图片

```http
GET /api/v1/images
```

**响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 5,
    "data": [
      "localhost:3128/f/image1.jpg",
      "localhost:3128/f/image2.png"
    ]
  }
}
```

### 获取单个随机图片

```http
GET /api/v1/images/random
```

**响应:**
```
image1.jpg
```

### 获取多个随机图片

```http
GET /api/v1/images/random/3
```

**参数:**
- `number`: 需要的图片数量（最多100张）

**响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 3,
    "data": [
      "localhost:3128/f/image1.jpg",
      "localhost:3128/f/image2.png",
      "localhost:3128/f/image3.gif"
    ]
  }
}
```

### 获取指定文件

```http
GET /f/:filename
```

直接返回文件内容

### 上传图片

```http
POST /api/v1/images/upload
Content-Type: multipart/form-data
```

**参数:**
- `files`: 要上传的文件列表

**响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "Upload completed",
    "total_uploaded": 2,
    "uploaded": [
      "localhost:3128/f/photo1.jpg",
      "localhost:3128/f/photo2.jpg"
    ]
  }
}
```

### 遗留API支持（向后兼容）

项目保持对旧API的支持：

```http
GET /v1/                    # 健康检查
GET /v1/all                 # 获取所有图片
GET /v1/bgimg               # 获取随机图片
GET /v1/get/:number         # 获取N个随机图片
POST /v1/upload             # 上传图片
GET /f/:filename            # 获取指定文件
```

## 配置说明

配置文件在 `internal/config/config.go` 中定义：

```go
// 服务器配置
Server: ServerConfig{
    Port:    ":3128",           // 服务器端口
    Env:     "development",     // 运行环境
    Timeout: 30,                // 超时时间（秒）
}

// 文件配置
File: FileConfig{
    UploadDir:  "./files",      // 上传目录
    MaxSize:    100,            // 最大文件大小（MB）
    AllowTypes: []string{...}   // 允许的文件类型
}
```

## 项目特性

✅ **标准化项目结构** - 遵循Go最佳实践  
✅ **分层架构** - Config、Service、Handler 清晰分离  
✅ **错误处理** - 统一的错误定义和响应  
✅ **日志管理** - 集中式日志管理  
✅ **CORS支持** - 完整的跨域资源共享支持  
✅ **文件管理** - 安全的文件上传和访问  
✅ **向后兼容** - 保持对旧API的支持  
✅ **响应统一** - 统一的JSON响应格式  

## 打包二进制文件

### Linux/Mac
```bash
go build -o image-sys main.go
```

### Windows
```bash
go build -o image-sys.exe main.go
```

### 跨平台编译

```bash
# Windows编译Linux版本
GOOS=linux GOARCH=amd64 go build -o image-sys main.go

# Linux编译Windows版本
GOOS=windows GOARCH=amd64 go build -o image-sys.exe main.go
```

## 目录说明

### `internal/` - 内部包
不能被其他项目导入的私有包，包含应用的核心逻辑。

### `pkg/` - 公共包
可被其他项目导入使用的包，如logger、errors等工具库。

### `cmd/` - 命令行应用
应用程序的入口点。

### `files/` - 文件存储
所有上传的图片存储在这个目录。

### `logs/` - 日志目录
应用运行的日志文件。

## 开发建议

1. **添加新功能时**：
   - 在 `internal/service/` 中添加业务逻辑
   - 在 `internal/handler/` 中添加HTTP处理
   - 在 `internal/router/` 中注册新路由

2. **错误处理**：
   - 使用 `pkg/errors` 中的预定义错误
   - 使用 `pkg/utils.ErrorResponse()` 返回错误

3. **日志记录**：
   - 使用 `pkg/logger.GetLogger()` 获取日志实例
   - 记录关键操作和错误

4. **响应格式**：
   - 始终使用 `pkg/utils` 中的响应函数
   - 保持响应格式的一致性

## 常见问题

### Q: 如何修改上传目录？
A: 修改 `internal/config/config.go` 中的 `UploadDir` 配置。

### Q: 如何添加新的中间件？
A: 在 `internal/middleware/` 中创建新文件，然后在 `internal/router/router.go` 中注册。

### Q: 如何限制文件上传大小？
A: 修改 `internal/config/config.go` 中的 `MaxSize` 字段。

## 许可证

MIT

## 作者

GantoHo
