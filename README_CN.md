# Kairoa CLI

Kairoa 开发者工具箱的命令行版本，包含 50+ 实用工具。

[![Build Status](https://github.com/covoyage/kairoa-cli/workflows/Build/badge.svg)](https://github.com/covoyage/kairoa-cli/actions)
[![Release](https://img.shields.io/github/release/covoyage/kairoa-cli.svg)](https://github.com/covoyage/kairoa-cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue.svg)](https://golang.org)

## 为什么选择 CLI？

Kairoa CLI 专为 **AI Agent 集成** 而设计。通过提供统一的命令行接口，它可以与以下 AI 编程助手无缝集成：

- **OpenClaw** - AI 驱动的开发环境
- **GitHub Copilot** - AI 结对编程助手
- **Claude Code** - Anthropic 的编程助手
- **Cursor** - AI 优先的代码编辑器
- **其他 AI Agent** - 任何可以执行 shell 命令的工具

### AI Agent 的优势

- **结构化输出** - JSON 和格式化输出，便于解析
- **管道友好** - 支持 Unix 管道，实现数据转换链
- **一致的接口** - 所有 50+ 工具遵循相同的 CLI 模式
- **无 GUI 依赖** - 纯 CLI，可在任何环境中运行
- **可脚本化** - 易于集成到自动化工作流中

### AI Agent 使用示例

```bash
# AI agent 可以为数据库记录生成 UUID
kairoa uuid v4 -c 5

# 格式化 JSON API 响应
kairoa http get https://api.example.com/data | kairoa json format

# 生成安全密码
kairoa password -n 32 --no-special

# 计算文件哈希值进行完整性检查
kairoa hash file ./package.json

# 转换数据格式
kairoa data csv2json < data.csv
```

## 功能特性

- **哈希计算**：MD5、SHA1、SHA256、SHA384、SHA512、RIPEMD160
- **UUID 生成**：多种选项生成 UUID
- **Base64 编解码**：Base64 字符串编码和解码
- **JSON 格式化**：格式化、压缩和验证 JSON
- **URL 编解码**：URL 编码和解码字符串
- **JWT 解码**：解码 JWT 令牌
- **时间工具**：转换时间戳和获取当前时间
- **密码生成器**：生成安全的随机密码
- **十六进制编解码**：十六进制编码和解码
- **HMAC 计算**：使用各种算法计算 HMAC
- **DNS 查询**：查询 DNS 记录
- **IP 查询**：查询 IP 地址信息
- **HTTP 客户端**：发送 HTTP 请求
- **WebSocket 客户端**：测试 WebSocket 连接
- **二维码生成器**：生成二维码
- **ASCII 艺术**：将文本转换为 ASCII 艺术
- **颜色转换器**：在各种颜色格式之间转换
- **进制转换器**：在数字进制之间转换
- **罗马数字**：在罗马数字和阿拉伯数字之间转换
- **Cron 解析器**：解析 Cron 表达式
- **SQL 格式化器**：格式化 SQL 查询
- **数据转换器**：在 CSV 和 JSON 之间转换
- **配置转换器**：在配置格式之间转换（JSON、YAML、TOML）
- **Docker 命令**：生成 Docker 命令
- **Git 命令**：生成 Git 命令
- **密码保险箱**：安全存储密码
- **环境管理器**：管理 .env 文件
- **图像处理**：处理和转换图像
- **PDF 工具**：PDF 信息和签名验证
- **端口扫描器**：扫描开放端口
- **TLS 检查器**：检查 TLS/SSL 版本
- **证书查看器**：查看 SSL 证书
- **正则表达式测试器**：测试正则表达式
- **文本处理**：文本统计和差异比较
- **模拟数据生成器**：生成模拟数据
- **IBAN 验证器**：验证 IBAN 号码
- **OTP 生成器**：生成 TOTP/HOTP 代码
- **RSA 密钥生成器**：生成 RSA 密钥对
- **坐标转换器**：在坐标格式之间转换
- **HTTP 状态**：HTTP 状态码参考
- **MIME 类型**：MIME 类型查询
- **用户代理解析器**：解析 User-Agent 字符串
- **基本认证**：生成基本认证头
- **密码强度**：检查密码强度
- **文件权限**：计算 chmod 权限
- **键盘键码**：显示键盘键码
- **更多...**

## 安装

### Homebrew（macOS / Linux）

```bash
brew install covoyage/tap/kairoa
```

### Scoop（Windows）

```powershell
scoop bucket add covoyage https://github.com/covoyage/scoop-bucket
scoop install kairoa
```

### 快速安装脚本（macOS / Linux）

```bash
curl -sSL https://raw.githubusercontent.com/covoyage/kairoa-cli/main/install.sh | bash
```

### 手动下载

从 [GitHub Releases](https://github.com/covoyage/kairoa-cli/releases/latest) 下载对应平台的预编译二进制文件：

| 平台 | 文件 |
|------|------|
| macOS（Apple Silicon） | `kairoa_darwin_arm64.tar.gz` |
| macOS（Intel） | `kairoa_darwin_x86_64.tar.gz` |
| Linux（x86_64） | `kairoa_linux_x86_64.tar.gz` |
| Linux（ARM64） | `kairoa_linux_arm64.tar.gz` |
| Windows（x86_64） | `kairoa_windows_x86_64.zip` |

```bash
# macOS / Linux 示例
tar -xzf kairoa_*.tar.gz
sudo mv kairoa /usr/local/bin/
kairoa version
```

### 从源码构建

```bash
git clone https://github.com/covoyage/kairoa-cli.git
cd kairoa-cli
go build -o kairoa .
```

## 快速开始

```bash
# 计算哈希
kairoa hash text "hello world"

# 生成 UUID
kairoa uuid v4

# Base64 编码 / 解码
kairoa base64 encode "hello"
kairoa base64 decode "aGVsbG8="

# 格式化 JSON
echo '{"a":1}' | kairoa json format

# 查看版本
kairoa version

# 查看所有命令
kairoa --help
```

## 使用方法

### 哈希

```bash
# 计算文本哈希
kairoa hash text "hello world"

# 计算文件哈希
kairoa hash file /path/to/file

# 使用特定算法
kairoa hash text "hello world" -a sha256,md5
```

### UUID

```bash
# 生成单个 UUID
kairoa uuid v4

# 生成多个 UUID
kairoa uuid v4 -c 5

# 生成 ULID
kairoa uuid ulid
```

### Base64

```bash
# 编码
kairoa base64 encode "hello world"

# 解码
kairoa base64 decode "aGVsbG8gd29ybGQ="
```

### JSON

```bash
# 格式化 JSON
echo '{"a":1,"b":2}' | kairoa json format

# 压缩 JSON
echo '{"a": 1, "b": 2}' | kairoa json minify

# 验证 JSON
echo '{"a":1}' | kairoa json validate
```

### 时间

```bash
# 获取当前时间
kairoa time now

# 转换时间戳
kairoa time convert 1609459200
```

### 二维码

```bash
# 生成二维码
kairoa qr "https://example.com" -a

# 保存到文件
kairoa qr "https://example.com" -o qr.png
```

### DNS 查询

```bash
# 查询 DNS 记录
kairoa dns lookup google.com

# 特定记录类型
kairoa dns lookup google.com -t MX
```

### HTTP 客户端

```bash
# GET 请求
kairoa http get https://api.example.com/users

# POST 请求
kairoa http post https://api.example.com/users -d '{"name":"John"}'
```

### 密码生成器

```bash
# 生成密码
kairoa password

# 自定义长度
kairoa password -n 20

# 排除特定字符
kairoa password --no-special
```

### 查看所有命令

```bash
kairoa --help

# 获取特定命令的帮助
kairoa hash --help
```

## 国际化

Kairoa CLI 支持多种语言：

```bash
# 使用中文
kairoa -l zh hash text "hello"

# 设置默认语言
kairoa lang set zh
```

## Shell 补全

### Bash

```bash
kairoa completion bash > /etc/bash_completion.d/kairoa
```

### Zsh

```bash
kairoa completion zsh > "${fpath[1]}/_kairoa"
```

### Fish

```bash
kairoa completion fish > ~/.config/fish/completions/kairoa.fish
```

## 开发

```bash
# 克隆仓库
git clone https://github.com/covoyage/kairoa-cli.git
cd kairoa-cli

# 安装依赖
go mod download

# 构建
go build -o kairoa .

# 运行测试
go test ./...

# 热重载运行（需要 air）
air
```

## 贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 致谢

- 灵感来源于 [Kairoa](https://github.com/covoyage/kairoa) 桌面应用程序
- 使用 [Cobra](https://github.com/spf13/cobra) CLI 框架构建
- 使用各种开源库 - 详见 go.mod 获取完整列表
