#!/bin/bash

# 遇到错误立即退出
set -e

APP_NAME="Sowhp"
BUILD_DIR="build"

echo "=========================================="
echo "开始构建 $APP_NAME ..."
echo "=========================================="

# 清理并创建构建目录
if [ -d "$BUILD_DIR" ]; then
    echo "[-] 清理旧的构建文件..."
    rm -rf "$BUILD_DIR"
fi
mkdir -p "$BUILD_DIR"

# 编译参数说明:
# CGO_ENABLED=0: 禁用CGO，生成静态链接的二进制文件，便于分发
# -ldflags "-s -w": 
#   -s: 省略符号表和调试信息
#   -w: 省略DWARF符号表
#   这可以显著减小二进制文件体积

# 1. macOS (Darwin)
echo "[+] 正在编译 macOS (AMD64)..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/${APP_NAME}_darwin_amd64" main.go

echo "[+] 正在编译 macOS (ARM64/M1/M2)..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o "$BUILD_DIR/${APP_NAME}_darwin_arm64" main.go

# 2. Linux
echo "[+] 正在编译 Linux (AMD64)..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/${APP_NAME}_linux_amd64" main.go

# 3. Windows
echo "[+] 正在编译 Windows (AMD64)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/${APP_NAME}_windows_amd64.exe" main.go

echo "=========================================="
echo "构建成功！"
echo "产物目录: $BUILD_DIR"
echo "=========================================="
ls -lh "$BUILD_DIR"
