# Sowhp 网站截图工具

## 工具介绍

  本项目是基于 [sh1yan/Sowhp](https://github.com/sh1yan/Sowhp) 的优化修改版本。

## 功能特点

- **批量截图**：支持批量处理URL列表，快速生成网页截图
- **智能报告**：自动生成HTML和文本格式的详细报告
- **智能重试**：HTTPS失败时自动尝试HTTP请求
- **实时进度**：显示截图进度和处理状态

## 主要改进和优化

### ✨ 日志系统优化
- 格式化输出显示：`[√]` 成功、`[×]` 错误、`[!]` 警告、`[?]` 调试、`[i]` 详细信息
- 统一括号颜色为白色，符号保持原有颜色
- 优化错误消息格式：`[!] [错误详情] 访问 URL 时遇到错误，重试中...`

### 🔧 错误处理改进
- 增强的重试机制，支持SSL证书错误和超时错误的智能重试
- 详细的错误分类：SSL错误、超时错误、DNS错误、连接拒绝等
- 改进的状态码获取和响应内容展示


## 环境要求

- **Go语言**：Go 1.25 或更高版本
- **Chrome浏览器**：用于网页截图功能
- **操作系统**：支持 Windows、macOS、Linux

## 安装和编译

### 1. 克隆项目
```bash
git clone <your-repo-url>
cd Sowhp
```

### 2. 编译项目
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

## 项目结构

```
├── main.go              # 程序入口
├── core/
│   └── working.go       # 核心工作流程
├── concert/
│   ├── directoriescreat.go  # 目录创建
│   └── logger/          # 日志系统
│       ├── config.go
│       ├── level.go
│       ├── log.go
│       └── logsync.go
├── scripts/
│   ├── Reporting.go     # 报告生成
│   ├── conditionals.go  # 条件判断
│   └── ipurlrelated.go  # URL处理和截图
└── result/              # 输出目录
    └── result_YYYYMMDD_NNNN/
        ├── data/        # 截图文件
        ├── *.html       # HTML报告
        └── *.txt        # 文本报告
```

## 输出说明

程序运行后会在 `result` 目录下生成以下文件：
- **截图文件**：保存在 `data/` 子目录中
- **HTML报告**：包含截图预览和详细信息的网页报告
- **文本报告**：纯文本格式的处理结果

## 更新记录

### 最新版本改进
- ✅ 优化日志系统，使用符号化显示
- ✅ 改进错误处理和重试机制
- ✅ 优化进度条显示为白色
- ✅ 改进HTML报告生成和样式
- ✅ 优化文件夹命名规则（年月日_序号格式）
- ✅ 清理代码注释，提高代码简洁性
- ✅ 修复HTML响应内容显示问题
- ✅ 统一错误消息格式
- ✅ 优化URL格式处理

### 原版功能保持
- ✅ 批量网页截图功能
- ✅ 多种URL格式支持
- ✅ 并发处理能力
- ✅ Chrome浏览器集成
- ✅ 跨平台支持

## 致谢

- 感谢 [sh1yan](https://github.com/sh1yan) 开发的原版 [Sowhp](https://github.com/sh1yan/Sowhp) 工具
