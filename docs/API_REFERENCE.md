# 🔌 完整API参考手册

## 目录
1. [认证](#认证)
2. [核心API](#核心api)
3. [工具API](#工具api)
4. [错误处理](#错误处理)
5. [请求示例](#请求示例)

---

## 认证

### API密钥

所有需要写入权限的端点都需要在请求头中提供API密钥：

```http
X-API-Key: demo-key-12345
```

**默认密钥:**
- `demo-key-12345` (演示密钥, 永久有效)
- `test-key-67890` (测试密钥, 永久有效)

### 创建新密钥

```http
POST /api/v1/auth/create-key
X-API-Key: demo-key-12345
Content-Type: application/json

{
  "name": "my-api-key",
  "expires_in_days": 90
}
```

**响应:**
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

## 核心API

### 1. 健康检查

#### 请求
```http
GET /api/v1/health
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "ok",
    "version": "1.0.0"
  },
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1705862000,
    "duration_ms": 1
  }
}
```

---

### 2. 获取图片列表

#### 请求
```http
GET /api/v1/images
```

**查询参数:**
- `limit` (可选): 返回数量限制, 默认100

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 150,
    "images": [
      {
        "filename": "photo1.jpg",
        "url": "http://localhost:3128/f/photo1.jpg",
        "size": 2097152,
        "mime_type": "image/jpeg",
        "mod_time": "2026-01-22T10:30:00Z"
      },
      {
        "filename": "photo2.png",
        "url": "http://localhost:3128/f/photo2.png",
        "size": 1048576,
        "mime_type": "image/png",
        "mod_time": "2026-01-22T10:25:00Z"
      }
    ]
  }
}
```

---

### 3. 获取随机图片

#### 请求
```http
GET /api/v1/images/random
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "filename": "photo42.jpg",
    "url": "http://localhost:3128/f/photo42.jpg",
    "size": 3145728,
    "mime_type": "image/jpeg"
  }
}
```

---

### 4. 分页获取图片

#### 请求
```http
GET /api/v1/images/paginated?page=1&page_size=20
```

**查询参数:**
- `page`: 页码 (从1开始)
- `page_size`: 每页数量 (1-100, 默认20)

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "page": 1,
    "page_size": 20,
    "total": 150,
    "total_pages": 8,
    "images": [
      {
        "filename": "photo1.jpg",
        "url": "http://localhost:3128/f/photo1.jpg",
        "size": 2097152,
        "mime_type": "image/jpeg",
        "mod_time": "2026-01-22T10:30:00Z"
      }
    ]
  }
}
```

---

### 5. 搜索和过滤

#### 请求
```http
GET /api/v1/images/search?filename=photo&size_min=1000000&type=jpg
```

**查询参数:**
- `filename` (可选): 文件名关键词
- `size_min` (可选): 最小大小 (字节)
- `size_max` (可选): 最大大小 (字节)
- `type` (可选): 文件类型 (jpg/png/gif等)

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 45,
    "images": [
      {
        "filename": "photo_big_1.jpg",
        "url": "http://localhost:3128/f/photo_big_1.jpg",
        "size": 5242880,
        "mime_type": "image/jpeg",
        "mod_time": "2026-01-22T10:30:00Z"
      }
    ]
  }
}
```

---

### 6. 获取元数据

#### 请求
```http
GET /api/v1/images/meta
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_files": 150,
    "total_size": 5368709120,
    "average_size": 35791394,
    "format_breakdown": {
      "jpg": 100,
      "png": 30,
      "gif": 20
    }
  }
}
```

---

### 7. 上传单个文件

#### 请求
```http
POST /api/v1/upload
X-API-Key: demo-key-12345
Content-Type: multipart/form-data

file: [binary data]
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "filename": "photo.jpg",
    "url": "http://localhost:3128/f/photo.jpg",
    "size": 2097152,
    "mime_type": "image/jpeg"
  }
}
```

---

### 8. 批量上传

#### 请求
```http
POST /api/v1/upload/multi
X-API-Key: demo-key-12345
Content-Type: multipart/form-data

files: [binary data 1]
files: [binary data 2]
files: [binary data 3]
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 3,
    "uploaded": 3,
    "files": [
      {
        "index": 0,
        "filename": "photo1.jpg",
        "url": "http://localhost:3128/f/photo1.jpg",
        "size": 2097152,
        "mime_type": "image/jpeg",
        "progress": 33,
        "message": "success"
      },
      {
        "index": 1,
        "filename": "photo2.jpg",
        "url": "http://localhost:3128/f/photo2.jpg",
        "size": 3145728,
        "mime_type": "image/jpeg",
        "progress": 67,
        "message": "success"
      },
      {
        "index": 2,
        "filename": "photo3.jpg",
        "url": "http://localhost:3128/f/photo3.jpg",
        "size": 1048576,
        "mime_type": "image/jpeg",
        "progress": 100,
        "message": "success"
      }
    ]
  }
}
```

---

### 9. 删除单个文件

#### 请求
```http
DELETE /api/v1/images/photo.jpg
X-API-Key: demo-key-12345
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "filename": "photo.jpg",
    "message": "file deleted successfully"
  }
}
```

---

### 10. 批量删除

#### 请求
```http
POST /api/v1/images/delete
X-API-Key: demo-key-12345
Content-Type: application/json

{
  "filenames": ["photo1.jpg", "photo2.png", "photo3.gif"]
}
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 3,
    "deleted": 3,
    "failed": 0,
    "results": [
      {
        "filename": "photo1.jpg",
        "success": true
      },
      {
        "filename": "photo2.png",
        "success": true
      },
      {
        "filename": "photo3.gif",
        "success": true
      }
    ]
  }
}
```

---

## 工具API

### 1. 获取文件统计

#### 请求
```http
GET /api/v1/util/statistics
```

#### 响应
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
      },
      ".gif": {
        "count": 20,
        "size": 429496320,
        "size_str": "0.40 GB",
        "percentage": 8.00
      }
    },
    "largest_file": "photo_4k.jpg",
    "largest_file_size": 52428800
  }
}
```

---

### 2. 获取磁盘使用情况

#### 请求
```http
GET /api/v1/util/disk-usage
```

#### 响应
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

### 3. 导出指定文件

#### 请求
```http
POST /api/v1/util/export
X-API-Key: demo-key-12345
Content-Type: application/json

{
  "filenames": ["photo1.jpg", "photo2.png"]
}
```

#### 响应
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

---

### 4. 导出所有文件

#### 请求
```http
POST /api/v1/util/export-all
X-API-Key: demo-key-12345
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "zip_path": "files/export_20260122_145945.zip",
    "file_count": 150,
    "total_size": 5368709120,
    "size_str": "5.00 GB",
    "compressed": true,
    "download_url": "http://localhost:3128/f/export_20260122_145945.zip"
  }
}
```

---

### 5. 执行清理操作

#### 请求
```http
POST /api/v1/util/cleanup
X-API-Key: demo-key-12345
Content-Type: application/json

{
  "remove_orphan_thumbnails": true,
  "remove_old_files": true,
  "max_file_age_days": 30,
  "remove_empty_dirs": true
}
```

#### 响应
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

---

### 6. 生成缩略图

#### 请求
```http
POST /api/v1/util/generate-thumbnails?filenames=photo1.jpg,photo2.png
X-API-Key: demo-key-12345
```

#### 响应
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

## 错误处理

### 常见错误码

| 代码 | 含义 | 说明 |
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

```json
{
  "code": 400,
  "message": "error",
  "errors": {
    "filename": "invalid filename format"
  },
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1705862000,
    "duration_ms": 2
  }
}
```

### 限流错误

```json
{
  "code": 403,
  "message": "rate limit exceeded",
  "data": {
    "retry_after": 1
  }
}
```

---

## 请求示例

### Shell (curl)

```bash
# 获取统计信息
curl http://localhost:3128/api/v1/util/statistics

# 上传文件
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -F "files=@photo.jpg" \
  http://localhost:3128/api/v1/upload/multi

# 导出文件
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"filenames": ["photo1.jpg"]}' \
  http://localhost:3128/api/v1/util/export
```

### Python

```python
import requests

# 获取统计信息
response = requests.get('http://localhost:3128/api/v1/util/statistics')
print(response.json())

# 上传文件
with open('photo.jpg', 'rb') as f:
    files = {'files': f}
    headers = {'X-API-Key': 'demo-key-12345'}
    response = requests.post(
        'http://localhost:3128/api/v1/upload/multi',
        files=files,
        headers=headers
    )
    print(response.json())
```

### JavaScript

```javascript
// 获取统计信息
fetch('http://localhost:3128/api/v1/util/statistics')
  .then(res => res.json())
  .then(data => console.log(data));

// 上传文件
const formData = new FormData();
formData.append('files', fileInput.files[0]);

fetch('http://localhost:3128/api/v1/upload/multi', {
  method: 'POST',
  headers: {
    'X-API-Key': 'demo-key-12345'
  },
  body: formData
})
  .then(res => res.json())
  .then(data => console.log(data));
```

### Go

```go
package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
)

func main() {
    // 获取统计信息
    resp, _ := http.Get("http://localhost:3128/api/v1/util/statistics")
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))

    // 上传文件
    file, _ := os.Open("photo.jpg")
    defer file.Close()

    req, _ := http.NewRequest("POST", 
        "http://localhost:3128/api/v1/upload",
        file)
    req.Header.Set("X-API-Key", "demo-key-12345")

    client := &http.Client{}
    resp, _ = client.Do(req)
    defer resp.Body.Close()
}
```

