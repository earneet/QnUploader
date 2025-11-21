# 七牛云命令行上传工具

一个类似ClaudeCode的命令行工具，支持快捷键打开、拖拽文件上传和文件路径输入上传，实时反馈云端访问路径。

## 功能特性

- 🚀 **交互式上传** - 支持拖拽文件和输入文件路径
- ⌨️ **快捷键支持** - 可配置全局快捷键（开发中）
- 📋 **自动复制链接** - 上传成功后自动复制云端访问链接
- 📊 **进度显示** - 实时显示上传进度
- 🔧 **配置管理** - 支持配置七牛云信息和快捷键
- 🌐 **跨平台** - 支持 Windows、macOS、Linux

## 快速开始

### 1. 安装Go环境

首先确保您的系统已安装 Go 1.21 或更高版本：

```bash
# 检查Go版本
go version
```

### 2. 下载和编译

```bash
# 克隆项目
git clone <repository-url>
cd QiNiuUploader

# 编译项目
go build ./cmd/qiniu-uploader

# 或者安装到系统路径
go install ./cmd/qiniu-uploader
```

### 3. 初始化配置

首次使用需要配置七牛云信息：

```bash
./qu config init
```

按照提示输入：
- Access Key
- Secret Key
- Bucket 名称
- 域名（可选）

### 4. 开始使用

#### 交互式上传模式
```bash
./qu upload
```

进入交互模式后，您可以：
- 拖拽文件到终端窗口
- 输入文件路径上传
- 输入 `list` 查看已上传文件
- 输入 `config` 查看当前配置
- 输入 `quit` 退出

#### 直接上传文件
```bash
# 上传单个文件
./qu upload /path/to/image.jpg

# 或者使用 -f 参数
./qu upload -f /path/to/image.jpg
```

## 使用示例

### 交互模式示例

```bash
$ ./qu upload
🚀 七牛云上传工具 - 交互模式
==================================================
支持以下操作:
  1. 输入文件路径上传 (支持拖拽文件到终端)
  2. 输入 'list' 查看已上传文件
  3. 输入 'config' 显示当前配置
  4. 输入 'quit' 或 'exit' 退出
==================================================

📁 请输入文件路径或命令: /Users/username/Pictures/photo.jpg
📤 正在上传: photo.jpg
[============================                      ] 60.0% 已用: 2s 剩余: 1s
✅ 上传成功!
📁 文件名: photo.jpg
📊 文件大小: 2.34 MB
🔗 访问链接: https://example.com/images/1234567890.jpg
🔑 存储Key: images/1234567890.jpg
📋 云端链接已复制到剪贴板

📁 请输入文件路径或命令: list

📚 已上传文件列表:
--------------------------------------------------------------------------------
 1. photo.jpg
    大小: 2.34 MB | 上传时间: 2025-01-08 14:30:25
    链接: https://example.com/images/1234567890.jpg

 2. screenshot.png
    大小: 1.56 MB | 上传时间: 2025-01-08 14:28:10
    链接: https://example.com/images/1234567891.png
--------------------------------------------------------------------------------

📁 请输入文件路径或命令: quit
👋 再见!
```

### 拖拽上传示例

1. 打开文件管理器
2. 选择要上传的图片文件
3. 拖拽文件到终端窗口
4. 文件路径会自动填充
5. 按回车键开始上传

## 配置说明

### 配置文件位置

配置文件保存在：`~/.config/qu/config.yaml`

### 配置示例

```yaml
qiniu_access_key: "your_access_key"
qiniu_secret_key: "your_secret_key"
qiniu_bucket: "your_bucket_name"
qiniu_domain: "your_domain.com"

hotkey_keys:
  - 85  # U键
hotkey_ctrl: true
hotkey_shift: true
hotkey_alt: false

auto_copy_url: true
show_progress: true
```

### 环境变量

您也可以使用环境变量配置：

```bash
export QINIU_ACCESS_KEY="your_access_key"
export QINIU_SECRET_KEY="your_secret_key"
export QINIU_BUCKET="your_bucket_name"
export QINIU_DOMAIN="your_domain.com"
```

## 命令参考

### 主命令

```bash
qu [command]
```

### 可用命令

- `upload` - 上传文件到七牛云
- `config` - 配置管理
- `service` - 启动后台服务（开发中）
- `version` - 显示版本信息

### Upload 命令

```bash
# 交互式上传
qu upload

# 上传指定文件
qu upload /path/to/file.jpg

# 使用 -f 参数
qu upload -f /path/to/file.jpg
```

### Config 命令

```bash
# 初始化配置
qu config init

# 显示当前配置
qu config show
```

## 支持的文件类型

- JPEG/JPG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)
- BMP (.bmp)

**文件大小限制**: 最大 10MB

## 故障排除

### 常见问题

1. **"七牛云客户端未初始化"**
   - 运行 `qu config init` 初始化配置

2. **"文件不存在"**
   - 检查文件路径是否正确
   - 确保文件有读取权限

3. **"文件大小超过限制"**
   - 当前限制为 10MB
   - 请压缩图片或选择较小的文件

4. **"不支持的文件类型"**
   - 请确保上传的是支持的图片格式

### 调试模式

设置环境变量查看详细日志：

```bash
export QINIU_UPLOADER_DEBUG=true
qu upload
```

## 开发说明

### 项目结构

```
QiNiuUploader/
├── cmd/
│   └── qiniu-uploader/
│       └── main.go          # 程序入口
├── internal/
│   ├── cli/                 # 命令行逻辑
│   │   ├── app.go           # 应用框架
│   │   ├── interactive.go   # 交互式界面
│   │   ├── dragdrop.go      # 拖拽功能
│   │   └── progress.go      # 进度显示
│   └── config/              # 配置管理
│       └── config.go
├── pkg/
│   └── qiniu/               # 七牛云SDK封装
│       └── client.go
└── README.md
```

### 依赖包

- `github.com/spf13/cobra` - 命令行框架
- `github.com/spf13/viper` - 配置管理
- `github.com/qiniu/go-sdk/v7` - 七牛云SDK
- `github.com/joho/godotenv` - 环境变量管理

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！