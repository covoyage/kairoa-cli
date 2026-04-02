# Kairoa CLI

A command-line version of the Kairoa developer toolbox with 50+ utilities.

[English](README.md) | [中文](README_CN.md)

[![Build Status](https://github.com/covoyage/kairoa-cli/workflows/Build/badge.svg)](https://github.com/covoyage/kairoa-cli/actions)
[![Release](https://img.shields.io/github/release/covoyage/kairoa-cli.svg)](https://github.com/covoyage/kairoa-cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org)

## Why CLI?

Kairoa CLI is specifically designed for **AI Agent integration**. By providing a unified command-line interface, it enables seamless integration with AI coding assistants like:

- **OpenClaw** - AI-powered development environment
- **GitHub Copilot** - AI pair programmer
- **Claude Code** - Anthropic's coding assistant
- **Cursor** - AI-first code editor
- **Other AI Agents** - Any tool that can execute shell commands

### Benefits for AI Agents

- **Structured Output** - JSON and formatted output for easy parsing
- **Pipe-friendly** - Works with Unix pipes for data transformation chains
- **Consistent Interface** - All 50+ tools follow the same CLI patterns
- **No GUI Dependencies** - Pure CLI, runs in any environment
- **Scriptable** - Easy to integrate into automation workflows

### Example: AI Agent Usage

```bash
# AI agent can generate UUIDs for database records
kairoa uuid v4 -c 5

# Format JSON API responses
kairoa http get https://api.example.com/data | kairoa json format

# Generate secure passwords
kairoa password -n 32 --no-special

# Calculate file hashes for integrity checks
kairoa hash file ./package.json

# Convert data formats
kairoa data csv2json < data.csv
```

## Features

- **Hash Calculation**: MD5, SHA1, SHA256, SHA384, SHA512, RIPEMD160
- **UUID Generation**: Generate UUIDs with various options
- **Base64 Encoding/Decoding**: Encode and decode Base64 strings
- **JSON Formatting**: Format, minify, and validate JSON
- **URL Encoding/Decoding**: URL encode and decode strings
- **JWT Decoding**: Decode JWT tokens
- **Time Utilities**: Convert timestamps and get current time
- **Password Generator**: Generate secure random passwords
- **Hex Encoding/Decoding**: Hex encode and decode strings
- **HMAC Calculation**: Calculate HMAC with various algorithms
- **DNS Lookup**: Query DNS records
- **IP Lookup**: Query IP address information
- **HTTP Client**: Send HTTP requests
- **WebSocket Client**: Test WebSocket connections
- **QR Code Generator**: Generate QR codes
- **ASCII Art**: Convert text to ASCII art
- **Color Converter**: Convert between color formats
- **Base Converter**: Convert between number bases
- **Roman Numerals**: Convert between Roman numerals and Arabic numbers
- **Cron Parser**: Parse cron expressions
- **SQL Formatter**: Format SQL queries
- **Data Converter**: Convert between CSV and JSON
- **Config Converter**: Convert between config formats (JSON, YAML, TOML)
- **Docker Commands**: Generate Docker commands
- **Git Commands**: Generate Git commands
- **Password Vault**: Securely store passwords
- **Environment Manager**: Manage .env files
- **Image Processing**: Process and convert images
- **PDF Utilities**: PDF information and signature verification
- **Port Scanner**: Scan open ports
- **TLS Checker**: Check TLS/SSL versions
- **Certificate Viewer**: View SSL certificates
- **Regex Tester**: Test regular expressions
- **Text Processing**: Text statistics and diff
- **Mock Data Generator**: Generate mock data
- **IBAN Validator**: Validate IBAN numbers
- **OTP Generator**: Generate TOTP/HOTP codes
- **RSA Key Generator**: Generate RSA key pairs
- **Coordinate Converter**: Convert between coordinate formats
- **HTTP Status**: HTTP status code reference
- **MIME Types**: MIME type lookup
- **User-Agent Parser**: Parse User-Agent strings
- **Basic Auth**: Generate Basic Authentication headers
- **Password Strength**: Check password strength
- **File Permissions**: Calculate chmod permissions
- **Keyboard Keycodes**: Show keyboard key codes
- **And more...**

## Installation

### Quick Install (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/covoyage/kairoa-cli/main/install.sh | bash
```

### Homebrew (macOS/Linux)

```bash
brew tap covoyage/tap
brew install kairoa
```

### Manual Installation

#### macOS

```bash
# Intel Mac
curl -L -o kairoa.tar.gz https://github.com/covoyage/kairoa-cli/releases/latest/download/kairoa_darwin_x86_64.tar.gz

# Apple Silicon Mac
curl -L -o kairoa.tar.gz https://github.com/covoyage/kairoa-cli/releases/latest/download/kairoa_darwin_arm64.tar.gz

tar -xzf kairoa.tar.gz
sudo mv kairoa /usr/local/bin/

# 验证安装
kairoa --version
```

#### Linux

```bash
curl -L -o kairoa.tar.gz https://github.com/covoyage/kairoa-cli/releases/latest/download/kairoa_linux_x86_64.tar.gz
tar -xzf kairoa.tar.gz
sudo mv kairoa /usr/local/bin/

# Verify installation
kairoa --version
```

#### Windows

Download the latest release from [GitHub Releases](https://github.com/covoyage/kairoa-cli/releases) and extract to a directory in your PATH.

### Build from Source

```bash
git clone https://github.com/covoyage/kairoa-cli.git
cd kairoa-cli
go build -o kairoa .
```

## Usage

### Hash

```bash
# Calculate hash of text
kairoa hash text "hello world"

# Calculate hash of file
kairoa hash file /path/to/file

# Use specific algorithms
kairoa hash text "hello world" -a sha256,md5
```

### UUID

```bash
# Generate a single UUID
kairoa uuid v4

# Generate multiple UUIDs
kairoa uuid v4 -c 5

# Generate ULID
kairoa uuid ulid
```

### Base64

```bash
# Encode
kairoa base64 encode "hello world"

# Decode
kairoa base64 decode "aGVsbG8gd29ybGQ="
```

### JSON

```bash
# Format JSON
echo '{"a":1,"b":2}' | kairoa json format

# Minify JSON
echo '{"a": 1, "b": 2}' | kairoa json minify

# Validate JSON
echo '{"a":1}' | kairoa json validate
```

### Time

```bash
# Get current time
kairoa time now

# Convert timestamp
kairoa time convert 1609459200
```

### QR Code

```bash
# Generate QR code
kairoa qr "https://example.com" -a

# Save to file
kairoa qr "https://example.com" -o qr.png
```

### DNS Lookup

```bash
# Lookup DNS records
kairoa dns lookup google.com

# Specific record type
kairoa dns lookup google.com -t MX
```

### HTTP Client

```bash
# GET request
kairoa http get https://api.example.com/users

# POST request
kairoa http post https://api.example.com/users -d '{"name":"John"}'
```

### Password Generator

```bash
# Generate password
kairoa password

# Custom length
kairoa password -n 20

# Exclude certain characters
kairoa password --no-special
```

### View All Commands

```bash
kairoa --help

# Get help for specific command
kairoa hash --help
```

## Internationalization

Kairoa CLI supports multiple languages:

```bash
# Use Chinese
kairoa -l zh hash text "hello"

# Set default language
kairoa lang set zh
```

## Shell Completion

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

## Development

```bash
# Clone repository
git clone https://github.com/covoyage/kairoa-cli.git
cd kairoa-cli

# Install dependencies
go mod download

# Build
go build -o kairoa .

# Run tests
go test ./...

# Run with hot reload (requires air)
air
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the [Kairoa](https://github.com/covoyage/kairoa) desktop application
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses various open-source libraries - see go.mod for full list
