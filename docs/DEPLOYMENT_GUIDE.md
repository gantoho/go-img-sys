# 🚀 部署和运维指南

## 目录
1. [本地开发](#本地开发)
2. [编译构建](#编译构建)
3. [Docker部署](#docker部署)
4. [生产部署](#生产部署)
5. [监控和维护](#监控和维护)
6. [故障排查](#故障排查)

---

## 本地开发

### 环境要求

- Go 1.18+
- 64位操作系统 (Windows/Linux/macOS)
- 至少2GB可用磁盘空间
- 网络访问权限 (用于mod下载)

### 快速启动

```bash
# 1. 克隆项目
git clone <repository-url> go-img-sys
cd go-img-sys

# 2. 下载依赖
go mod download
go mod tidy

# 3. 运行应用
go run ./cmd/image-sys/main.go

# 4. 访问服务
curl http://localhost:3128/api/v1/health
```

### 开发配置

创建 `configs/config.yml`:

```yaml
# 服务器配置
server:
  port: 3128
  max_upload_size: 104857600  # 100MB

# 文件配置
files:
  upload_dir: "./files"
  max_files: 10000

# API限流配置
rate_limit:
  requests_per_second: 100
  concurrent_per_ip: 10

# 日志配置
logging:
  level: debug
  output: ./logs/app.log

# 缓存配置
cache:
  ttl_seconds: 300
```

---

## 编译构建

### Windows编译

```bash
# 编译为exe
go build -o build/image-sys.exe ./cmd/image-sys

# 运行编译后的程序
.\build\image-sys.exe
```

### Linux/macOS编译

```bash
# 编译为可执行文件
go build -o build/image-sys ./cmd/image-sys

# 运行
./build/image-sys
```

### 交叉编译

```bash
# 为Linux编译
GOOS=linux GOARCH=amd64 go build -o build/image-sys-linux ./cmd/image-sys

# 为Windows编译
GOOS=windows GOARCH=amd64 go build -o build/image-sys.exe ./cmd/image-sys

# 为macOS编译
GOOS=darwin GOARCH=amd64 go build -o build/image-sys-mac ./cmd/image-sys
```

### 优化构建

```bash
# 生产优化构建 (减小体积和增加性能)
go build -ldflags="-s -w" -o build/image-sys ./cmd/image-sys
```

---

## Docker部署

### Dockerfile

项目已包含 `Dockerfile`, 使用以下命令:

```bash
# 构建镜像
docker build -t go-img-sys:latest .

# 运行容器
docker run -d \
  --name img-sys \
  -p 3128:3128 \
  -v /data/images:/app/files \
  -v /data/logs:/app/logs \
  go-img-sys:latest

# 查看日志
docker logs -f img-sys

# 停止容器
docker stop img-sys

# 删除容器
docker rm img-sys
```

### Docker Compose

使用 `deployments/docker-compose.yml`:

```bash
# 启动服务
docker-compose -f deployments/docker-compose.yml up -d

# 查看日志
docker-compose -f deployments/docker-compose.yml logs -f image-sys

# 停止服务
docker-compose -f deployments/docker-compose.yml down
```

### 卷挂载最佳实践

```bash
# 保存数据到主机
docker run -d \
  --name img-sys \
  -p 3128:3128 \
  -v /host/path/images:/app/files \
  -v /host/path/logs:/app/logs \
  go-img-sys:latest
```

---

## 生产部署

### 系统要求

- CPU: 2核+
- 内存: 2GB+
- 磁盘: 10GB+ (根据图片数量调整)
- 网络: 100Mbps+

### Linux systemd服务

创建 `/etc/systemd/system/image-sys.service`:

```ini
[Unit]
Description=Go Image System
After=network.target

[Service]
Type=simple
User=nobody
WorkingDirectory=/opt/image-sys
ExecStart=/opt/image-sys/image-sys
Restart=always
RestartSec=10

# 日志输出
StandardOutput=journal
StandardError=journal

# 资源限制
MemoryLimit=2G
CPUQuota=50%

[Install]
WantedBy=multi-user.target
```

### 启动服务

```bash
# 重新加载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start image-sys

# 设置开机自启
sudo systemctl enable image-sys

# 查看状态
sudo systemctl status image-sys

# 查看日志
sudo journalctl -u image-sys -f
```

### Nginx反向代理

配置 `/etc/nginx/sites-available/image-sys`:

```nginx
upstream image_sys {
    server localhost:3128;
}

server {
    listen 80;
    server_name images.example.com;

    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name images.example.com;

    # SSL证书
    ssl_certificate /etc/letsencrypt/live/images.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/images.example.com/privkey.pem;

    # SSL配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 日志
    access_log /var/log/nginx/image-sys-access.log;
    error_log /var/log/nginx/image-sys-error.log;

    # 请求限制
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20;

    # 缓冲设置
    proxy_buffering on;
    proxy_buffer_size 4k;
    proxy_buffers 8 4k;

    # 代理设置
    location / {
        proxy_pass http://image_sys;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 静态文件缓存
    location ~* \.(jpg|png|gif|webp|ico)$ {
        proxy_pass http://image_sys;
        expires 7d;
        add_header Cache-Control "public, immutable";
    }
}
```

启用站点:

```bash
sudo ln -s /etc/nginx/sites-available/image-sys /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 监控和维护

### 性能监控

创建监控脚本 `scripts/monitor.sh`:

```bash
#!/bin/bash

while true; do
    echo "=== Image System Status ==="
    
    # 检查进程
    ps aux | grep image-sys | grep -v grep
    
    # 检查端口
    netstat -tuln | grep 3128
    
    # 检查磁盘使用
    curl -s http://localhost:3128/api/v1/util/disk-usage | jq '.data'
    
    # 检查文件数量
    curl -s http://localhost:3128/api/v1/util/statistics | jq '.data | {total_files, total_size_str}'
    
    sleep 60
done
```

### 日志分析

```bash
# 查看错误日志
tail -f logs/error.log

# 统计错误数量
grep -c "ERROR" logs/error.log

# 查看特定时间段的日志
grep "2026-01-22" logs/app.log | head -20

# 实时监控
tail -f logs/app.log | grep "POST\|DELETE"
```

### 定期维护任务

```bash
#!/bin/bash
# 每周清理脚本 (scripts/cleanup.sh)

# 备份数据库
tar -czf backup_$(date +%Y%m%d).tar.gz files/

# 执行清理
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "remove_orphan_thumbnails": true,
    "remove_old_files": true,
    "max_file_age_days": 90,
    "remove_empty_dirs": true
  }' \
  http://localhost:3128/api/v1/util/cleanup

# 创建完整备份
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/util/export-all \
  -o backup_full_$(date +%Y%m%d).zip

echo "Cleanup completed at $(date)" >> logs/maintenance.log
```

在crontab中定时执行:

```bash
# 每周日凌晨2点执行清理
0 2 * * 0 /opt/image-sys/scripts/cleanup.sh
```

---

## 故障排查

### 常见问题

#### 1. 端口被占用

```bash
# Windows
netstat -ano | findstr :3128
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :3128
kill -9 <PID>
```

#### 2. 权限问题

```bash
# Linux/Mac - 给予执行权限
chmod +x ./build/image-sys

# 给予文件夹权限
chmod -R 755 ./files
chmod -R 755 ./logs
```

#### 3. 磁盘空间不足

```bash
# 查看磁盘使用
df -h

# 清理旧日志
find logs -name "*.log" -mtime +30 -delete

# 执行清理操作
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"remove_old_files": true, "max_file_age_days": 60}' \
  http://localhost:3128/api/v1/util/cleanup
```

#### 4. 高内存占用

```bash
# 重启服务
sudo systemctl restart image-sys

# 或删除缓存
curl -X POST \
  -H "X-API-Key: demo-key-12345" \
  http://localhost:3128/api/v1/cache/clear
```

#### 5. 图片无法上传

```bash
# 检查文件夹权限
ls -la files/

# 检查磁盘空间
df -h files/

# 查看错误日志
tail -f logs/error.log | grep "upload"
```

### 日志级别设置

```bash
# 通过环境变量设置
export LOG_LEVEL=debug
./build/image-sys

# 在配置文件中设置
logging:
  level: info  # debug/info/warn/error/fatal
```

### 健康检查命令

```bash
# 基础健康检查
curl http://localhost:3128/api/v1/health

# 完整健康检查脚本
#!/bin/bash

echo "Checking Image System Health..."

# 检查API响应
if curl -f http://localhost:3128/api/v1/health > /dev/null 2>&1; then
    echo "✓ API is healthy"
else
    echo "✗ API is down"
    exit 1
fi

# 检查磁盘
USAGE=$(curl -s http://localhost:3128/api/v1/util/disk-usage | jq '.data.percentage')
echo "Disk usage: $USAGE%"

# 检查文件数
TOTAL=$(curl -s http://localhost:3128/api/v1/util/statistics | jq '.data.total_files')
echo "Total files: $TOTAL"

# 检查进程
if pgrep -f "image-sys" > /dev/null; then
    echo "✓ Process is running"
else
    echo "✗ Process is not running"
    exit 1
fi

echo "Health check completed successfully"
```

---

## 性能优化建议

### 应用层优化

1. **增加缓存TTL**: 如果文件变化不频繁
   ```go
   cache.SetTTL(15 * time.Minute) // 15分钟
   ```

2. **调整页面大小**: 对于大量文件
   ```bash
   GET /api/v1/images/paginated?page_size=50
   ```

3. **启用GZIP压缩**: 在Nginx中
   ```nginx
   gzip on;
   gzip_types application/json;
   gzip_min_length 1000;
   ```

### 系统层优化

1. **增加文件描述符限制**
   ```bash
   ulimit -n 65536
   ```

2. **优化网络参数**
   ```bash
   sysctl -w net.core.somaxconn=65535
   sysctl -w net.ipv4.tcp_max_syn_backlog=65535
   ```

3. **使用SSD存储**: 关键IO操作存储

### 监控建议

- CPU使用率 < 50%
- 内存使用率 < 60%
- 磁盘使用率 < 80%
- API响应时间 < 200ms
- 错误率 < 0.1%

---

## 备份和恢复

### 自动备份

创建备份脚本 `scripts/backup.sh`:

```bash
#!/bin/bash

BACKUP_DIR="/backup/image-sys"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份文件
tar -czf $BACKUP_DIR/files_$DATE.tar.gz files/

# 保留最近30天的备份
find $BACKUP_DIR -name "files_*.tar.gz" -mtime +30 -delete

echo "Backup completed: files_$DATE.tar.gz"
```

### 恢复

```bash
# 解压备份
tar -xzf /backup/image-sys/files_20260122_100000.tar.gz

# 重启服务
sudo systemctl restart image-sys
```

