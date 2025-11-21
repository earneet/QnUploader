#!/bin/bash

# 七牛云上传工具编译脚本

echo "🚀 开始编译七牛云上传工具..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go编译器"
    echo "请先安装Go 1.21或更高版本"
    echo "下载地址: https://golang.org/dl/"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ 错误: Go版本过低"
    echo "当前版本: $GO_VERSION"
    echo "需要版本: $REQUIRED_VERSION 或更高"
    exit 1
fi

echo "✅ Go版本检查通过: $GO_VERSION"

# 下载依赖
echo "📦 下载依赖包..."
if ! go mod download; then
    echo "❌ 依赖下载失败"
    exit 1
fi

# 编译程序
echo "🔨 编译程序..."
if ! go build -o qiniu-uploader ./cmd/qiniu-uploader; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功!"
echo ""
echo "📋 使用说明:"
echo "   1. 初始化配置: ./qiniu-uploader config init"
echo "   2. 上传文件: ./qiniu-uploader upload"
echo "   3. 查看帮助: ./qiniu-uploader --help"
echo ""
echo "💡 提示: 您可以将程序移动到系统PATH目录，方便使用"
echo "   sudo mv qiniu-uploader /usr/local/bin/"

# 检查文件权限
chmod +x qiniu-uploader

echo ""
echo "🎉 七牛云上传工具已准备就绪!"