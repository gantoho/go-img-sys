# 🚀 快速使用指南

## 服务启动

```bash
# 编译
go build -o build/image-sys.exe ./cmd/image-sys

# 运行
./build/image-sys.exe
```

服务将在 `http://localhost:3128` 启动

## 常见操作示例

### 1️⃣ 健康检查
```bash
curl http://localhost:3128/api/v1/health
```

### 2️⃣ 查看所有图片（不需认证）
```bash
curl http://localhost:3128/api/v1/images
```

### 3️⃣ 分页查询（第一页，每页20张）
```bash
curl "http://localhost:3128/api/v1/images/paginated?page=1&page_size=20"
```

### 4️⃣ 搜索图片
```bash
# 按名称搜索
curl "http://localhost:3128/api/v1/images/search?filename=cat"

# 按类型搜索
curl "http://localhost:3128/api/v1/images/search?type=png"

# 按大小范围搜索（1MB到10MB）
curl "http://localhost:3128/api/v1/images/search?min_size=1048576&max_size=10485760"

# 组合搜索
curl "http://localhost:3128/api/v1/images/search?filename=photo&type=jpg&min_size=1000000"
```

### 5️⃣ 获取图片元数据
```bash
curl http://localhost:3128/api/v1/images/metadata
```

### 6️⃣ 获取随机图片
```bash
# 单个随机图片
curl http://localhost:3128/api/v1/images/random

# 5个随机图片
curl http://localhost:3128/api/v1/images/random/5
```

### 7️⃣ 上传图片（需要API密钥）

#### 方式1：使用Header传递密钥
```bash
curl -X POST -H "X-API-Key: demo-key-12345" \
  -F "files=@photo1.jpg" \
  -F "files=@photo2.png" \
  http://localhost:3128/api/v1/images/upload
```

#### 方式2：使用Query参数传递密钥
```bash
curl -X POST \
  -F "files=@photo.jpg" \
  "http://localhost:3128/api/v1/images/upload?api_key=demo-key-12345"
```

### 8️⃣ 删除图片（需要API密钥）

#### 删除单个图片
```bash
curl -X DELETE \
  -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/images/photo.jpg
```

#### 批量删除
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"filenames": ["photo1.jpg", "photo2.png"]}' \
  http://localhost:3128/api/v1/images/delete
```

### 9️⃣ 创建新API密钥（需要现有密钥认证）
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"expire_days": 30}' \
  http://localhost:3128/api/v1/admin/api-keys
```

响应示例：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "api_key": "api-key-0",
    "expire_days": 30,
    "message": "API key created successfully. Please save it safely!"
  },
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1705862000,
    "duration_ms": 5
  }
}
```

### 🔟 查看所有API密钥（需要认证）
```bash
curl -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/admin/api-keys
```

### 1️⃣1️⃣ 撤销API密钥（需要认证）
```bash
curl -X DELETE \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"api_key": "target-key-to-revoke"}' \
  http://localhost:3128/api/v1/admin/api-keys
```

## 🔐 默认API密钥

| 密钥 | 有效期 | 用途 |
|------|--------|------|
| `demo-key-12345` | 30天 | 演示和测试 |
| `test-key-67890` | 7天 | 短期测试 |

## 📊 API响应格式

所有响应都遵循统一格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 具体数据
  },
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1705862000,
    "duration_ms": 15
  }
}
```

**说明**：
- `code`: HTTP状态码
- `message`: 响应消息
- `data`: 返回的具体数据
- `metadata.version`: API版本
- `metadata.timestamp`: 响应时间戳
- `metadata.duration_ms`: 请求处理耗时（毫秒）

## ⚠️ 常见错误处理

### 401 Unauthorized - 无效或缺少API密钥
```json
{
  "code": 401,
  "message": "invalid or expired API key"
}
```

### 429 Too Many Requests - 超出请求限制
```json
{
  "code": 429,
  "message": "rate limit exceeded"
}
```

### 400 Bad Request - 不支持的图片格式
```json
{
  "code": 400,
  "message": "unsupported image format. Supported formats: jpg, jpeg, png, gif, webp, bmp, ico, svg"
}
```

## 🔄 向后兼容API

旧版API仍然可用（无需更改现有客户端）：

```bash
# 旧API - 无需密钥
GET  /v1/all                    # 获取所有图片
GET  /v1/bgimg                  # 随机图片
GET  /v1/get/:number            # N个随机图片

# 旧API - 需要密钥
POST /v1/upload                 # 上传图片
```

## 📝 日志文件

日志保存位置：
- 标准输出：INFO/WARN/DEBUG 级别
- `logs/error.log`：ERROR/FATAL 级别

## 🛠️ 性能指标

- **请求限流**：100请求/秒
- **并发限制**：10连接/IP地址
- **缓存时间**：5分钟
- **最大上传文件**：100MB
- **支持的文件格式**：8种（jpg, png, gif, webp, bmp, ico, svg）

## 💡 最佳实践

1. **定期创建新密钥**：每个应用/环境使用不同的密钥
2. **监控日志**：定期检查 `logs/error.log`
3. **及时撤销过期密钥**：使用DELETE接口
4. **使用Header传递密钥**：比Query参数更安全
5. **分页查询**：大数据集使用分页而不是一次全量查询
6. **合理设置过期时间**：1-365天根据安全需求选择

## 🐛 故障排除

### 问题：无法连接到服务
```bash
# 检查服务是否运行
netstat -an | findstr :3128

# 查看服务日志
# 查看启动输出是否有错误
```

### 问题：上传文件失败
```bash
# 1. 检查API密钥是否正确
# 2. 检查文件格式是否支持
# 3. 检查文件大小是否超过限制（100MB）
# 4. 检查 ./files 目录权限
```

### 问题：请求被限流
```bash
# 如果收到429错误，等待一段时间后重试
# 降低请求频率
```
