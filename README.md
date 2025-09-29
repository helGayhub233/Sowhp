# Sowhp 网站截图工具
<img src="https://raw.githubusercontent.com/helGayhub233/Sowhp/refs/heads/main/images/example_1.png" alt="示例图片" width="410" height="300"><img src="https://raw.githubusercontent.com/helGayhub233/Sowhp/refs/heads/main/images/example_2.png" alt="示例图片" width="400" height="300">
## 工具介绍

本项目是基于 [sh1yan/Sowhp](https://github.com/sh1yan/Sowhp) 的优化版本。

## 功能特点

- **批量截图**：支持批量处理URL列表，快速生成网页截图
- **智能报告**：自动生成HTML和CSV格式的详细报告
- **智能重试**：HTTPS失败时自动尝试HTTP请求
- **实时进度**：显示截图进度和处理状态

## 改进优化

### ✨ 日志系统优化
- 格式化输出显示
- 优化错误消息格式

### 🔧 错误处理改进
- 增强的重试机制，支持SSL证书错误和超时错误的智能重试
- 详细的错误分类：SSL错误、超时错误、DNS错误、连接拒绝等
- 改进的状态码获取和响应内容展示


## 环境要求

- **Go语言**：Go 1.25 或更高版本

## 安装编译

### 克隆项目
```bash
git clone <your-repo-url>
cd Sowhp
```

### 手动编译
```bash
# 基础编译
go build -o sowhp

# 优化编译（减小文件体积）
go build -ldflags="-s -w" -trimpath -o sowhp

# 使用 UPX 进一步压缩（可选）
upx -9 sowhp
```

## 使用说明

### 基本用法
```bash
./sowhp -f urls.txt
```

### 参数说明
- `-f`：指定包含URL列表的文本文件路径（必需参数）
- `-log`：设置日志输出详细程度（可选，默认值：3）
  - `1`：仅错误
  - `2`：错误和警告
  - `3`：错误、警告和信息（默认）
  - `4`：包含调试信息

### 使用示例
```bash
# 基础使用
./sowhp -f urls.txt

# 启用调试模式
./sowhp -f urls.txt -log 4

# 处理大量URL时使用信息模式
./sowhp -f large_urls.txt -log 3
```

### URL文件格式
创建一个文本文件，每行一个URL：
```
https://www.example.com
http://192.168.1.1:8080
example.org
192.168.1.100
```

## 输出说明

程序运行后会在 `result` 目录下生成以下文件：
- **截图文件**：保存在 `data/` 子目录中
- **HTML报告**：包含截图预览和详细信息的网页报告
- **CSV报告**：生成CSV格式的处理结果用于批处理

## 更新记录

### 功能改进
- ✅ 优化日志系统，使用符号化显示
- ✅ 改进错误处理和重试机制
- ✅ 优化进度条显示为白色
- ✅ 改进HTML报告生成和样式
- ✅ 优化文件夹命名规则（年月日_序号格式）
- ✅ 清理代码注释，提高代码简洁性
- ✅ 修复HTML响应内容显示问题
- ✅ 统一错误消息格式
- ✅ 优化URL格式处理

## 致谢

- 感谢 [sh1yan](https://github.com/sh1yan) 开发的原版 [Sowhp](https://github.com/sh1yan/Sowhp) 工具
