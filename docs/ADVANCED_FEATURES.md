# 📦 高级功能详解

## 概述
本文档详细说明6项新增高级功能，这些功能进一步增强了系统的实用性和可维护性。

---

## 1️⃣ 图片缩略图生成

**文件**: `pkg/imageutil/imageutil.go`

### 功能说明
- 自动生成指定大小的缩略图
- 保持原始图片的宽高比
- 支持 JPEG 和 PNG 格式
- 可配置质量参数（JPEG）

### API使用

触发缩略图生成（后台任务）：
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  "http://localhost:3128/api/v1/util/generate-thumbnails?filenames=photo1.jpg,photo2.png"
```

### 配置示例
```go
config := imageutil.ThumbnailConfig{
    Width:   200,
    Height:  200,
    Quality: 85,
}

err := imageutil.GenerateThumbnail("files/photo.jpg", "files/thumbs/photo.jpg", config)
```

### 缩略图保存位置
```
files/
  ├── photo1.jpg
  ├── photo2.png
  └── thumbs/
      ├── photo1.jpg
      └── photo2.png
```

---

## 2️⃣ 定时清理功能

**文件**: `internal/service/maintenance_service.go`

### 功能说明
- 清理孤立缩略图（对应原文件已删除）
- 删除过期文件（可配置年龄）
- 删除空目录
- 后台定时自动执行
- 详细的清理报告

### API调用

#### 手动触发清理
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "remove_orphan_thumbnails": true,
    "remove_old_files": true,
    "max_file_age_days": 30,
    "remove_empty_dirs": true
  }' \
  http://localhost:3128/api/v1/util/cleanup
```

### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "files_removed": 5,
    "thumbnails_removed": 5,
    "dirs_removed": 2,
    "size_freed": 2097152,
    "errors": []
  },
  "metadata": {
    "version": "1.0.0",
    "timestamp": 1705862000,
    "duration_ms": 125
  }
}
```

### 自动清理配置
在应用启动时启用（可修改间隔）：
```go
maintService := service.NewMaintenanceService()
maintService.StartAutoCleanup(24 * time.Hour) // 每24小时执行一次
```

---

## 3️⃣ 文件夹分类存储

**功能说明**
系统支持按日期自动分类存储上传的文件（实现方式：可在上传时组织文件夹结构）

### 文件夹结构示例
```
files/
  ├── 2026/
  │   ├── 01/
  │   │   ├── 21/
  │   │   │   ├── photo1.jpg
  │   │   │   └── photo2.png
  │   │   └── 22/
  │   │       └── photo3.jpg
  │   └── 02/
  └── thumbs/
      ├── 2026/01/21/
      └── 2026/01/22/
```

### 实现建议
在上传时使用日期路径：
```go
dateDir := time.Now().Format("2006/01/02")
uploadPath := filepath.Join(config.UploadDir, dateDir, filename)
```

---

## 4️⃣ 统计分析API

**文件**: `internal/service/statistics_service.go`

### 获取文件统计信息

#### 请求
```bash
curl http://localhost:3128/api/v1/util/statistics
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

### 获取磁盘使用情况

#### 请求
```bash
curl http://localhost:3128/api/v1/util/disk-usage
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

### 统计指标说明
- **total_files**: 总文件数
- **total_size**: 总大小（字节）
- **average_file_size**: 平均文件大小
- **format_stats**: 按格式分类的统计数据
- **percentage**: 该格式文件占比
- **disk_usage**: 已用/总容量比例

---

## 5️⃣ 批量导出功能

**文件**: `internal/service/export_service.go`

### 功能说明
- 打包多个文件为 ZIP 压缩包
- 支持导出全部文件
- ZIP 文件自动命名（包含时间戳）
- 快速下载整个图库

### 导出多个文件

#### 请求
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "filenames": ["photo1.jpg", "photo2.png", "photo3.gif"]
  }' \
  http://localhost:3128/api/v1/util/export
```

#### 响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "zip_path": "files/export_20260122_145930.zip",
    "file_count": 3,
    "total_size": 15728640,
    "size_str": "15.00 MB",
    "compressed": true
  }
}
```

### 导出所有文件

#### 请求
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all
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
    "compressed": true
  }
}
```

### ZIP 文件下载
生成的 ZIP 文件可通过直接访问下载：
```
http://localhost:3128/f/export_20260122_145930.zip
```

---

## 6️⃣ 图片基础处理

**文件**: `pkg/imageutil/imageutil.go`

### 支持的操作

#### 1. 缩放图片 (Resize)
```bash
# 待实现的API端点
POST /api/v1/util/images/resize
{
  "filename": "photo.jpg",
  "width": 800,
  "height": 600
}
```

#### 2. 旋转图片 (Rotate)
```bash
# 支持 90, 180, 270 度旋转
POST /api/v1/util/images/rotate
{
  "filename": "photo.jpg",
  "degrees": 90
}
```

#### 3. 添加水印 (Watermark)
```bash
# 简单文字水印
POST /api/v1/util/images/watermark
{
  "filename": "photo.jpg",
  "text": "© 2026"
}
```

#### 4. 生成缩略图 (Thumbnail)
```bash
POST /api/v1/util/images/thumbnail
{
  "filename": "photo.jpg",
  "width": 200,
  "height": 200
}
```

### 实现示例

```go
// 缩放图片
err := imageutil.ResizeImage(
    "files/photo.jpg",
    "files/photo_resized.jpg",
    800, 600,
)

// 旋转图片
err := imageutil.RotateImage(
    "files/photo.jpg",
    "files/photo_rotated.jpg",
    90,
)

// 添加水印
err := imageutil.AddWatermark(
    "files/photo.jpg",
    "files/photo_watermarked.jpg",
    "© MyCompany 2026",
)
```

---

## 🔗 完整API路由汇总

### 公开端点（无需认证）

```
# 统计分析
GET  /api/v1/util/statistics         # 获取统计信息
GET  /api/v1/util/disk-usage         # 获取磁盘使用情况
```

### 受保护端点（需要API密钥）

```
# 导出功能
POST /api/v1/util/export             # 导出指定文件为ZIP
POST /api/v1/util/export-all         # 导出所有文件

# 维护操作
POST /api/v1/util/cleanup            # 执行清理操作
POST /api/v1/util/generate-thumbnails # 生成缩略图

# 图片处理
POST /api/v1/util/images/resize      # 缩放图片（待实现）
POST /api/v1/util/images/rotate      # 旋转图片（待实现）
POST /api/v1/util/images/watermark   # 添加水印（待实现）
POST /api/v1/util/images/thumbnail   # 生成缩略图（待实现）
```

---

## 💾 文件结构演变

### 上传后的完整结构
```
go-img-sys/
├── files/
│   ├── photo1.jpg
│   ├── photo2.png
│   ├── photo3.gif
│   ├── thumbs/                          # 缩略图目录
│   │   ├── photo1.jpg
│   │   ├── photo2.png
│   │   └── photo3.gif
│   ├── export_20260122_145930.zip       # 导出文件
│   └── export_20260122_145945.zip
├── logs/
│   ├── error.log
│   └── cleanup_20260122.log
└── ...
```

---

## 🚀 使用建议

### 1. 定期备份
使用导出功能定期备份所有图片：
```bash
curl -X POST -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all > backup.zip
```

### 2. 监控磁盘空间
定期检查磁盘使用情况：
```bash
curl http://localhost:3128/api/v1/util/disk-usage | jq '.data.percentage'
```

### 3. 自动清理策略
配置自动清理（每天凌晨3点）：
```go
// 在应用初始化时添加
maintService := service.NewMaintenanceService()
maintService.StartAutoCleanup(24 * time.Hour)
```

### 4. 缓存管理
系统会自动缓存文件列表（5分钟），删除文件后缓存会自动清除。

---

## 📊 性能优化

| 功能 | 性能 | 备注 |
|------|------|------|
| 统计分析 | O(n) | 遍历所有文件，首次较慢 |
| 清理操作 | 可配置 | 建议非高峰期执行 |
| 导出ZIP | ~100MB/s | 取决于磁盘速度 |
| 缩略图 | 实时生成 | 首次生成较慢，建议异步执行 |

---

## ⚠️ 注意事项

1. **路径安全**：所有文件操作都进行了路径遍历检查
2. **权限控制**：写入操作（导出、清理）需要API密钥
3. **磁盘限制**：导出大量文件时检查可用磁盘空间
4. **并发处理**：同时执行多个清理任务可能影响性能
5. **错误恢复**：清理操作不可逆，建议先导出备份

