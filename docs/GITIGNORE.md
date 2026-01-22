# 📋 .gitignore 配置说明

## 📝 什么是 .gitignore？

`.gitignore` 是 Git 的配置文件，用于**指定哪些文件和目录不需要提交到版本控制中**。

## 🎯 作用

- ✅ 避免提交编译产物和临时文件
- ✅ 避免提交敏感信息（密钥、密码等）
- ✅ 避免提交依赖库（如 node_modules）
- ✅ 保持仓库清洁和轻量

## 📁 分类说明

### 1️⃣ 编译产物和二进制文件
```gitignore
*.exe
*.dll
*.so
*.dylib
image-sys
image-sys.exe
build/
bin/
```
**原因**: 这些是编译生成的文件，应该由构建流程生成，不应存储在版本控制中。

### 2️⃣ 测试输出
```gitignore
*.test
*.out
*.cover
*.profile
coverage/
```
**原因**: 测试运行生成的临时文件和覆盖率报告。

### 3️⃣ Go依赖和工作区
```gitignore
vendor/
go.work
go.work.sum
```
**原因**: 
- `vendor/`: 如果使用 vendor 管理依赖，通常不提交（让 go.mod/go.sum 管理）
- `go.work`: Go workspace 配置，本地开发用

### 4️⃣ IDE和编辑器
```gitignore
.vscode/        # VS Code 设置
.idea/          # JetBrains IDE
*.swp           # Vim 交换文件
.DS_Store       # MacOS 文件夹配置
Thumbs.db       # Windows 缩略图数据库
```
**原因**: IDE 配置和编辑器临时文件是个人特定的，不应共享。

### 5️⃣ 环境和配置
```gitignore
.env
.env.local
.env.*.local
.env.production.local
```
**原因**: 包含敏感信息（API密钥、数据库密码等），不应提交到版本控制。

### 6️⃣ 日志文件
```gitignore
logs/
*.log
npm-debug.log*
yarn-error.log*
build-errors.log
```
**原因**: 运行时生成的日志文件，不需要版本控制。

### 7️⃣ 热加载和临时文件
```gitignore
tmp/
temp/
dist/
.air
air.log
```
**原因**: Air 热加载和其他临时构建文件。

### 8️⃣ 操作系统特定文件
```gitignore
.DS_Store          # MacOS
._*                # MacOS 隐藏文件
.Spotlight-V100    # MacOS
ehthumbs.db        # Windows
Thumbs.db          # Windows
```
**原因**: 不同操作系统生成的无关文件。

## ✅ 当前配置涵盖的内容

| 类别 | 文件/目录 | 说明 |
|------|---------|------|
| **编译产物** | `*.exe`, `*.so`, `build/`, `bin/` | 编译生成的文件 |
| **测试输出** | `*.test`, `*.out`, `*.profile`, `coverage/` | 测试相关文件 |
| **依赖** | `vendor/`, `go.work` | Go 依赖管理 |
| **IDE** | `.vscode/`, `.idea/` | 开发工具配置 |
| **敏感信息** | `.env*` | 环境变量文件 |
| **日志** | `logs/`, `*.log` | 应用日志 |
| **临时文件** | `tmp/`, `temp/`, `dist/` | 临时构建文件 |
| **系统文件** | `.DS_Store`, `Thumbs.db` | OS特定文件 |

## 🚀 使用示例

### 添加新的忽略规则
```bash
# 忽略单个文件
.myconfig

# 忽略文件夹
my_temp_folder/

# 忽略特定扩展名
*.bak
*.swp

# 忽略但保留某个特定文件（使用 ! 前缀）
!important_file.txt
```

### 验证 .gitignore 是否生效
```bash
# 显示哪些文件被忽略了
git status --ignored

# 或
git check-ignore -v <filename>
```

## 💡 最佳实践

### ✅ 应该忽略
- 编译产物 (`.exe`, `.so` 等)
- 依赖库 (`node_modules/`, `vendor/`)
- 配置文件 (`.env`)
- IDE 配置 (`.vscode/`, `.idea/`)
- 日志文件
- OS特定文件 (`.DS_Store`)
- 临时文件和缓存

### ❌ 不应该忽略
- 源代码 (`.go` 文件)
- 构建配置 (`go.mod`, `go.sum`)
- 文档 (`.md` 文件)
- 配置文件模板 (`.env.example`)
- 项目特定的脚本

## 📋 常见 Go 项目 .gitignore 模板

完整的 Go 项目 .gitignore 应该包含：

```gitignore
# 编译
*.o
*.a
*.so
*.exe
*.dylib
dist/
build/

# 测试
*.out
*.test

# IDE
.vscode/
.idea/
*.swp

# Go
vendor/
go.work

# 环境
.env
.env.local

# 其他
*.log
logs/
.DS_Store
```

## 🔗 相关资源

- [GitHub 官方 .gitignore 模板](https://github.com/github/gitignore)
- [Go 项目 .gitignore](https://github.com/github/gitignore/blob/main/Go.gitignore)

## 📌 项目特定说明

### files/ 目录
该目录用于存储用户上传的图片。根据需求：
- 如果图片应该被版本控制：**不忽略**
- 如果图片由用户上传存储：**可以忽略**

当前配置未忽略此目录，意味着上传的文件会被提交。如需更改：

```gitignore
# 忽略上传的文件但保留目录
files/*
!files/.gitkeep
```

### logs/ 目录
日志文件已配置为忽略，仓库中不会存储运行时日志。

---

**已完成：.gitignore 已根据 Go 项目最佳实践完善！** ✨
