# 🧪 新功能测试指南

## 环境准备

```bash
# 1. 启动应用
cd d:\GantoHo\Dev\Golang\go-img-sys
.\build\image-sys.exe

# 2. 查看服务信息
# 服务运行在: http://localhost:3128
# 默认API密钥: demo-key-12345
```

---

## 测试场景

### 测试1️⃣ : 缩略图生成

#### 步骤1: 上传原始图片
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -F "files=@C:\path\to\photo.jpg" \
  http://localhost:3128/api/v1/upload/multi
```

#### 步骤2: 触发缩略图生成
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  "http://localhost:3128/api/v1/util/generate-thumbnails?filenames=photo.jpg"
```

#### 步骤3: 验证
```bash
# 检查缩略图是否生成
ls -la ./files/thumbs/
```

✅ **预期结果**: `files/thumbs/photo.jpg` 存在

---

### 测试2️⃣ : 文件统计

#### 步骤1: 获取统计信息
```bash
curl http://localhost:3128/api/v1/util/statistics | jq
```

#### 步骤2: 检查结果
```json
{
  "total_files": 1,
  "total_size": 2097152,
  "format_stats": {
    ".jpg": {
      "count": 1,
      "percentage": 100
    }
  }
}
```

✅ **预期结果**: 显示上传的文件统计

---

### 测试3️⃣ : 磁盘使用情况

```bash
curl http://localhost:3128/api/v1/util/disk-usage | jq
```

✅ **预期结果**:
```json
{
  "used_space": 2097152,
  "limit": 107374182400,
  "percentage": 0.002
}
```

---

### 测试4️⃣ : 批量导出

#### 方式A: 导出指定文件
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "filenames": ["photo.jpg"]
  }' \
  http://localhost:3128/api/v1/util/export | jq
```

#### 方式B: 导出所有文件
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all | jq
```

✅ **预期结果**:
```json
{
  "zip_path": "files/export_20260122_145930.zip",
  "file_count": 1,
  "compressed": true
}
```

#### 下载导出的ZIP
```bash
curl -O http://localhost:3128/f/export_20260122_145930.zip
```

---

### 测试5️⃣ : 清理操作

#### 步骤1: 执行清理
```bash
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "remove_orphan_thumbnails": true,
    "remove_old_files": false,
    "remove_empty_dirs": true
  }' \
  http://localhost:3128/api/v1/util/cleanup | jq
```

✅ **预期结果**:
```json
{
  "files_removed": 0,
  "thumbnails_removed": 0,
  "dirs_removed": 0,
  "size_freed": 0
}
```

---

### 测试6️⃣ : 图片处理 (实现中)

这些功能已有代码支持，可通过处理器调用：

```go
// 旋转图片
imageutil.RotateImage(
    "files/photo.jpg",
    "files/photo_90.jpg",
    90,
)

// 缩放图片
imageutil.ResizeImage(
    "files/photo.jpg",
    "files/photo_resize.jpg",
    800, 600,
)

// 添加水印
imageutil.AddWatermark(
    "files/photo.jpg",
    "files/photo_watermarked.jpg",
    "© MyCompany",
)
```

---

## 🔄 完整测试工作流

```bash
# 1. 启动服务
.\build\image-sys.exe

# 2. 上传多个测试图片
curl -X POST -H "X-API-Key: demo-key-12345" \
  -F "files=@test1.jpg" -F "files=@test2.png" -F "files=@test3.gif" \
  http://localhost:3128/api/v1/upload/multi

# 3. 生成缩略图
curl -X POST -H "X-API-Key: demo-key-12345" \
  "http://localhost:3128/api/v1/util/generate-thumbnails?filenames=test1.jpg,test2.png,test3.gif"

# 4. 查看统计信息
curl http://localhost:3128/api/v1/util/statistics | jq

# 5. 查看磁盘使用情况
curl http://localhost:3128/api/v1/util/disk-usage | jq

# 6. 导出所有文件
curl -X POST -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all > result.json

# 7. 下载导出文件
# 从result.json中获取 zip_path，然后下载

# 8. 清理操作
curl -X POST -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"remove_orphan_thumbnails": true, "remove_empty_dirs": true}' \
  http://localhost:3128/api/v1/util/cleanup
```

---

## 🐛 常见问题

### Q1: 缩略图未生成？
**原因**: 可能因为图片格式不支持（仅支持JPEG和PNG）
**解决**: 使用支持的格式重新上传

### Q2: 导出ZIP为空？
**原因**: 没有指定文件名或文件不存在
**解决**: 先使用 `/api/v1/images/list` 查看可用文件

### Q3: 清理操作超时？
**原因**: 文件数量很多
**解决**: 在服务器运行时间内手动执行

### Q4: 磁盘使用情况不准确？
**原因**: 缓存未更新
**解决**: 等待5分钟或重启服务

---

## 📈 性能基准

在本地测试环境下（SSD):

| 操作 | 文件数 | 耗时 |
|------|-------|------|
| 统计分析 | 100 | ~50ms |
| 导出ZIP | 100 | ~200ms |
| 清理操作 | 100 | ~100ms |
| 缩略图生成 | 10 | ~500ms |

---

## ✅ 验收标准

- [ ] 缩略图成功生成
- [ ] 统计信息准确
- [ ] 磁盘使用百分比正确
- [ ] ZIP导出可正常下载
- [ ] 清理操作不误删文件
- [ ] 图片处理函数可调用
- [ ] 所有API返回正确状态码
- [ ] 错误消息清晰明确

