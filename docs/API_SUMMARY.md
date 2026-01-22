# 🔌 API接口汇总文档

## 📑 快速导航

- [认证](#认证)
- [基础API](#基础api-13个)
- [工具API](#工具api-6个)
- [请求示例](#请求示例)
- [错误响应](#错误响应)

---

## 🔐 认证

所有需要写入权限的API需在请求头中提供API密钥：

```
X-API-Key: demo-key-12345
```

**默认密钥:**
- `demo-key-12345` - 演示密钥 (永久)
- `test-key-67890` - 测试密钥 (永久)

---

## 基础API (13个)

### 1️⃣ 健康检查

| 方法 | 端点 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/health` | 公开 | 检查服务状态 |

**响应示例:**
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

---

### 2️⃣ 图片查询 (6个)

| # | 方法 | 端点 | 权限 | 说明 |
|---|------|------|------|------|
| 2.1 | GET | `/api/v1/images` | 公开 | 获取所有图片列表 |
| 2.2 | GET | `/api/v1/images/list` | 公开 | 获取图片列表(限制数量) |
| 2.3 | GET | `/api/v1/images/random` | 公开 | 获取随机一张图片 |
| 2.4 | GET | `/api/v1/images/paginated` | 公开 | 分页获取图片 |
| 2.5 | GET | `/api/v1/images/search` | 公开 | 搜索和过滤图片 |
| 2.6 | GET | `/api/v1/images/meta` | 公开 | 获取图片元数据统计 |

**参数说明:**

**2.4 分页参数:**
```
page: 页码 (从1开始)
page_size: 每页数量 (1-100, 默认20)
```

**2.5 搜索参数:**
```
filename: 文件名关键词
size_min: 最小大小(字节)
size_max: 最大大小(字节)
type: 文件类型 (jpg/png/gif等)
```

**响应示例 (2.1):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 10,
    "images": [
      {
        "filename": "photo.jpg",
        "url": "http://localhost:3128/f/photo.jpg",
        "size": 2097152,
        "mime_type": "image/jpeg",
        "mod_time": "2026-01-22T10:30:00Z"
      }
    ]
  }
}
```

---

### 3️⃣ 文件操作 (4个)

| # | 方法 | 端点 | 权限 | 说明 |
|---|------|------|------|------|
| 3.1 | POST | `/api/v1/upload` | 需密钥 | 上传单个文件 |
| 3.2 | POST | `/api/v1/upload/multi` | 需密钥 | 批量上传文件 |
| 3.3 | DELETE | `/api/v1/images/:filename` | 需密钥 | 删除指定文件 |
| 3.4 | POST | `/api/v1/images/delete` | 需密钥 | 批量删除文件 |

**3.1/3.2 请求:**
```bash
Content-Type: multipart/form-data
files: [二进制文件内容]
```

**3.3 请求:**
```bash
DELETE /api/v1/images/photo.jpg
X-API-Key: demo-key-12345
```

**3.4 请求:**
```json
{
  "filenames": ["photo1.jpg", "photo2.png"]
}
```

**响应示例 (3.2):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 2,
    "uploaded": 2,
    "files": [
      {
        "index": 0,
        "filename": "photo1.jpg",
        "url": "http://localhost:3128/f/photo1.jpg",
        "size": 2097152,
        "progress": 50
      },
      {
        "index": 1,
        "filename": "photo2.png",
        "url": "http://localhost:3128/f/photo2.png",
        "size": 1048576,
        "progress": 100
      }
    ]
  }
}
```

---

### 4️⃣ 密钥管理 (3个)

| # | 方法 | 端点 | 权限 | 说明 |
|---|------|------|------|------|
| 4.1 | POST | `/api/v1/auth/create-key` | 需密钥 | 创建新API密钥 |
| 4.2 | GET | `/api/v1/auth/keys` | 需密钥 | 列出所有密钥 |
| 4.3 | POST | `/api/v1/auth/revoke-key` | 需密钥 | 撤销密钥 |

**4.1 请求:**
```json
{
  "name": "my-api-key",
  "expires_in_days": 90
}
```

**4.1 响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "key": "new-api-key-xyz",
    "name": "my-api-key",
    "created_at": 1705862000,
    "expires_at": 1714638000
  }
}
```

---

## 工具API (6个)

### 5️⃣ 统计分析 (2个)

| # | 方法 | 端点 | 权限 | 说明 |
|---|------|------|------|------|
| 5.1 | GET | `/api/v1/util/statistics` | 公开 | 获取文件统计信息 |
| 5.2 | GET | `/api/v1/util/disk-usage` | 公开 | 获取磁盘使用情况 |

**5.1 响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_files": 150,
    "total_size": 5368709120,
    "total_size_str": "5.00 GB",
    "average_file_size": 35791394,
    "format_stats": {
      ".jpg": {
        "count": 100,
        "size": 3865099776,
        "size_str": "3.60 GB",
        "percentage": 71.95
      },
      ".png": {
        "count": 30,
        "size": 1073741824,
        "size_str": "1.00 GB",
        "percentage": 20.00
      }
    },
    "largest_file": "photo_4k.jpg",
    "largest_file_size": 52428800
  }
}
```

**5.2 响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "used_space": 5368709120,
    "used_space_str": "5.00 GB",
    "limit": 107374182400,
    "limit_str": "100.00 GB",
    "percentage": 5.00
  }
}
```

---

### 6️⃣ 导出功能 (2个)

| # | 方法 | 端点 | 权限 | 说明 |
|---|------|------|------|------|
| 6.1 | POST | `/api/v1/util/export` | 需密钥 | 导出指定文件为ZIP |
| 6.2 | POST | `/api/v1/util/export-all` | 需密钥 | 导出所有文件为ZIP |

**6.1 请求:**
```json
{
  "filenames": ["photo1.jpg", "photo2.png"]
}
```

**6.1 响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "zip_path": "files/export_20260122_145930.zip",
    "file_count": 2,
    "total_size": 3145728,
    "size_str": "3.00 MB",
    "compressed": true,
    "download_url": "http://localhost:3128/f/export_20260122_145930.zip"
  }
}
```

**6.2 请求:** (无参数)
```bash
POST /api/v1/util/export-all
X-API-Key: demo-key-12345
```

---

### 7️⃣ 维护操作 (2个)

| # | 方法 | 端点 | 权限 | 说明 |
|---|------|------|------|------|
| 7.1 | POST | `/api/v1/util/cleanup` | 需密钥 | 执行清理操作 |
| 7.2 | POST | `/api/v1/util/generate-thumbnails` | 需密钥 | 生成缩略图 |

**7.1 请求:**
```json
{
  "remove_orphan_thumbnails": true,
  "remove_old_files": true,
  "max_file_age_days": 30,
  "remove_empty_dirs": true
}
```

**7.1 响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "files_removed": 5,
    "thumbnails_removed": 5,
    "dirs_removed": 2,
    "size_freed": 2097152,
    "size_freed_str": "2.00 MB"
  }
}
```

**7.2 请求:**
```bash
POST /api/v1/util/generate-thumbnails?filenames=photo1.jpg,photo2.png
X-API-Key: demo-key-12345
```

**7.2 响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "thumbnail generation started",
    "files_to_process": 2
  }
}
```

---

## 📋 接口统计

| 分类 | 数量 | 说明 |
|------|------|------|
| **基础API** | 13个 | 健康检查、查询、上传、删除、密钥 |
| **工具API** | 6个 | 统计、导出、清理、缩略图 |
| **总计** | **19个** | **完整功能** |

| 权限 | 数量 | 说明 |
|------|------|------|
| **公开** | 15个 | 不需要API密钥 |
| **需密钥** | 4个 | 写入操作需要API密钥 |
| **总计** | **19个** | **安全保护** |

---

## 请求示例

### 使用curl

```bash
# 1. 健康检查
curl http://localhost:3128/api/v1/health

# 2. 获取统计信息
curl http://localhost:3128/api/v1/util/statistics

# 3. 上传单个文件
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -F "files=@photo.jpg" \
  http://localhost:3128/api/v1/upload

# 4. 批量上传
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -F "files=@photo1.jpg" \
  -F "files=@photo2.png" \
  http://localhost:3128/api/v1/upload/multi

# 5. 搜索图片
curl "http://localhost:3128/api/v1/images/search?filename=photo&type=jpg"

# 6. 分页查询
curl "http://localhost:3128/api/v1/images/paginated?page=1&page_size=20"

# 7. 导出指定文件
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"filenames": ["photo1.jpg"]}' \
  http://localhost:3128/api/v1/util/export

# 8. 导出所有文件
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all

# 9. 执行清理
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "remove_orphan_thumbnails": true,
    "remove_old_files": true,
    "max_file_age_days": 30
  }' \
  http://localhost:3128/api/v1/util/cleanup

# 10. 创建API密钥
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-key", "expires_in_days": 90}' \
  http://localhost:3128/api/v1/auth/create-key
```

### 使用Python

```python
import requests
import json

# 配置
BASE_URL = "http://localhost:3128/api/v1"
API_KEY = "demo-key-12345"

# 1. 健康检查
response = requests.get(f"{BASE_URL}/health")
print(response.json())

# 2. 获取统计
response = requests.get(f"{BASE_URL}/util/statistics")
print(response.json())

# 3. 上传文件
files = {'files': open('photo.jpg', 'rb')}
headers = {'X-API-Key': API_KEY}
response = requests.post(f"{BASE_URL}/upload/multi", files=files, headers=headers)
print(response.json())

# 4. 搜索图片
response = requests.get(f"{BASE_URL}/images/search", params={
    'filename': 'photo',
    'type': 'jpg'
})
print(response.json())

# 5. 导出文件
headers = {'X-API-Key': API_KEY}
data = {'filenames': ['photo1.jpg', 'photo2.png']}
response = requests.post(f"{BASE_URL}/util/export", json=data, headers=headers)
print(response.json())
```

### 使用JavaScript

```javascript
const BASE_URL = "http://localhost:3128/api/v1";
const API_KEY = "demo-key-12345";

// 1. 健康检查
fetch(`${BASE_URL}/health`)
  .then(res => res.json())
  .then(data => console.log(data));

// 2. 获取统计
fetch(`${BASE_URL}/util/statistics`)
  .then(res => res.json())
  .then(data => console.log(data));

// 3. 上传文件
const formData = new FormData();
formData.append('files', fileInput.files[0]);
fetch(`${BASE_URL}/upload/multi`, {
  method: 'POST',
  headers: {'X-API-Key': API_KEY},
  body: formData
})
  .then(res => res.json())
  .then(data => console.log(data));

// 4. 搜索图片
fetch(`${BASE_URL}/images/search?filename=photo&type=jpg`)
  .then(res => res.json())
  .then(data => console.log(data));

// 5. 导出文件
fetch(`${BASE_URL}/util/export`, {
  method: 'POST',
  headers: {
    'X-API-Key': API_KEY,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({filenames: ['photo1.jpg']})
})
  .then(res => res.json())
  .then(data => console.log(data));
```

---

## 错误响应

### 常见错误码

| 代码 | 说明 | 原因 |
|------|------|------|
| 200 | 成功 | 请求成功 |
| 400 | 请求错误 | 参数验证失败 |
| 401 | 未认证 | 缺少或无效API密钥 |
| 403 | 禁止访问 | 权限不足或限流 |
| 404 | 未找到 | 文件或端点不存在 |
| 409 | 冲突 | 文件已存在 |
| 413 | 请求体过大 | 文件大小超出限制 |
| 415 | 不支持的媒体类型 | 文件格式不支持 |
| 500 | 服务器错误 | 内部错误 |

### 错误响应示例

**400 - 参数错误:**
```json
{
  "code": 400,
  "message": "error",
  "errors": {
    "filename": "invalid filename format"
  }
}
```

**401 - 认证失败:**
```json
{
  "code": 401,
  "message": "unauthorized",
  "data": {
    "error": "invalid or expired API key"
  }
}
```

**403 - 限流:**
```json
{
  "code": 403,
  "message": "rate limit exceeded",
  "data": {
    "retry_after": 1
  }
}
```

**404 - 文件不存在:**
```json
{
  "code": 404,
  "message": "error",
  "data": {
    "error": "file not found"
  }
}
```

---

## 📊 响应格式统一

所有API响应遵循统一格式：

```json
{
  "code": 200,
  "message": "success|error",
  "data": {},
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1705862000,
    "duration_ms": 10
  }
}
```

**字段说明:**
- `code`: HTTP状态码
- `message`: 消息描述
- `data`: 响应数据
- `metadata.version`: API版本
- `metadata.timestamp`: 响应时间戳
- `metadata.duration_ms`: 请求耗时(毫秒)

---

## 🚀 快速参考

### 最常用的接口

```bash
# 查看服务状态
curl http://localhost:3128/api/v1/health

# 查看文件统计
curl http://localhost:3128/api/v1/util/statistics

# 上传文件
curl -X POST -H "X-API-Key: demo-key-12345" \
  -F "files=@photo.jpg" \
  http://localhost:3128/api/v1/upload/multi

# 获取文件列表
curl http://localhost:3128/api/v1/images

# 导出所有文件
curl -X POST -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all
```

---

## 📞 获取帮助

- 详细API参考: [docs/API_REFERENCE.md](docs/API_REFERENCE.md)
- 高级功能说明: [docs/ADVANCED_FEATURES.md](docs/ADVANCED_FEATURES.md)
- 测试指南: [docs/TESTING_GUIDE.md](docs/TESTING_GUIDE.md)

