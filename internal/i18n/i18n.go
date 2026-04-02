package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Locale represents a language locale
type Locale string

const (
	// English locale
	English Locale = "en"
	// Chinese locale
	Chinese Locale = "zh"
)

// Translator handles translations
type Translator struct {
	locale       Locale
	translations map[string]map[string]string
}

// NewTranslator creates a new translator
func NewTranslator() *Translator {
	t := &Translator{
		locale:       English,
		translations: make(map[string]map[string]string),
	}
	t.loadDefaultTranslations()
	return t
}

// SetLocale sets the current locale
func (t *Translator) SetLocale(locale Locale) {
	t.locale = locale
}

// GetLocale returns the current locale
func (t *Translator) GetLocale() Locale {
	return t.locale
}

// T translates a key
func (t *Translator) T(key string, args ...interface{}) string {
	// Try current locale first
	if trans, ok := t.translations[string(t.locale)][key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(trans, args...)
		}
		return trans
	}

	// Fallback to English
	if trans, ok := t.translations[string(English)][key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(trans, args...)
		}
		return trans
	}

	// Return key if not found
	return key
}

// LoadFromFile loads translations from a JSON file
func (t *Translator) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var translations map[string]map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	for locale, trans := range translations {
		if _, ok := t.translations[locale]; !ok {
			t.translations[locale] = make(map[string]string)
		}
		for key, value := range trans {
			t.translations[locale][key] = value
		}
	}

	return nil
}

// loadDefaultTranslations loads built-in translations
func (t *Translator) loadDefaultTranslations() {
	// English translations
	t.translations["en"] = map[string]string{
		// App
		"app.name":        "Kairoa",
		"app.description": "Desktop Utility Tools CLI",

		// Common
		"common.error":           "Error",
		"common.success":         "Success",
		"common.warning":         "Warning",
		"common.info":            "Info",
		"common.done":            "Done",
		"common.cancel":          "Cancel",
		"common.save":            "Save",
		"common.load":            "Load",
		"common.delete":          "Delete",
		"common.edit":            "Edit",
		"common.add":             "Add",
		"common.remove":          "Remove",
		"common.search":          "Search",
		"common.copy":            "Copy",
		"common.paste":           "Paste",
		"common.copied":          "Copied to clipboard",
		"common.notFound":        "Not found",
		"common.invalid":         "Invalid",
		"common.required":        "Required",
		"common.optional":        "Optional",
		"common.enabled":         "Enabled",
		"common.disabled":        "Disabled",
		"common.yes":             "Yes",
		"common.no":              "No",
		"common.true":            "True",
		"common.false":           "False",
		"common.input":           "Input",
		"common.output":          "Output",
		"common.result":          "Result",
		"common.format":          "Format",
		"common.type":            "Type",
		"common.value":           "Value",
		"common.key":             "Key",
		"common.name":            "Name",
		"common.description":     "Description",
		"common.example":         "Example",
		"common.usage":           "Usage",
		"common.help":            "Help",
		"common.version":         "Version",
		"common.author":          "Author",
		"common.license":         "License",

		// Commands - Hash
		"hash.short": "Calculate hash values",
		"hash.long":  `Calculate hash values for text or files using various algorithms (MD5, SHA1, SHA256, SHA384, SHA512, RIPEMD160).`,
		"hash.text":  "Calculate hash of a text string",
		"hash.file":  "Calculate hash of a file",
		"hash.md5":   "Calculate MD5 hash",
		"hash.sha1":  "Calculate SHA1 hash",
		"hash.sha256": "Calculate SHA256 hash",
		"hash.sha512": "Calculate SHA512 hash",
		"hash.algorithm": "Hash algorithm(s) to use",
		"hash.bcrypt": "Hash with Bcrypt",
		"hash.verify": "Verify hash",

		// Commands - UUID
		"uuid.short":    "Generate UUIDs",
		"uuid.long":     `Generate UUIDs (v4, v7) and ULIDs.`,
		"uuid.v4":       "Generate UUID v4 (random)",
		"uuid.v7":       "Generate UUID v7 (time-ordered)",
		"uuid.ulid":     "Generate ULID",
		"uuid.count":    "Number of UUIDs to generate",
		"uuid.uppercase": "Output in uppercase",

		// Commands - Base64
		"base64.short": "Base64 encode/decode",
		"base64.long":  `Encode and decode Base64 strings.`,
		"base64.encode": "Encode to Base64",
		"base64.decode": "Decode from Base64",
		"base64.url":    "Use URL-safe encoding",

		// Commands - JSON
		"json.short":    "JSON formatting utilities",
		"json.long":     `Format, validate, and query JSON data.`,
		"json.format":   "Format JSON",
		"json.minify":   "Minify JSON",
		"json.validate": "Validate JSON",
		"json.query":    "Query JSON with jq syntax",
		"json.sort":     "Sort JSON keys",

		// Commands - URL
		"url.short":    "URL encoding/decoding",
		"url.long":     `Encode and decode URL strings.`,
		"url.encode":   "Encode URL",
		"url.decode":   "Decode URL",
		"url.full":     "Encode full URL",
		"url.component": "Encode URL component",

		// Commands - JWT
		"jwt.short":   "Decode JWT tokens",
		"jwt.long":    `Decode and verify JWT tokens.`,
		"jwt.decode":  "Decode JWT token",
		"jwt.header":  "JWT Header",
		"jwt.payload": "JWT Payload",
		"jwt.signature": "JWT Signature",

		// Commands - Time
		"time.short":    "Time utilities",
		"time.long":     `Convert between timestamps and dates.`,
		"time.now":      "Show current time",
		"time.convert":  "Convert timestamp to date",
		"time.timestamp": "Get current timestamp",
		"time.format":   "Time format",
		"time.timezone": "Timezone",
		"time.unix":     "Unix timestamp (seconds)",
		"time.unixMs":   "Unix timestamp (milliseconds)",

		// Commands - Password
		"password.short":     "Password generator",
		"password.long":      `Generate secure random passwords.`,
		"password.length":    "Password length",
		"password.uppercase": "Include uppercase letters",
		"password.lowercase": "Include lowercase letters",
		"password.numbers":   "Include numbers",
		"password.special":   "Include special characters",
		"password.strength":  "Password strength",

		// Commands - QR Code
		"qr.short":    "Generate QR code",
		"qr.long":     `Generate QR codes from text or URLs.`,
		"qr.text":     "Text to encode",
		"qr.size":     "QR code size",
		"qr.level":    "Error correction level",
		"qr.terminal": "Display in terminal",

		// Commands - HTTP
		"http.short":   "HTTP client utilities",
		"http.long":    `Send HTTP requests and inspect responses.`,
		"http.get":     "Send GET request",
		"http.post":    "Send POST request",
		"http.put":     "Send PUT request",
		"http.delete":  "Send DELETE request",
		"http.headers": "Request headers",
		"http.body":    "Request body",
		"http.follow":  "Follow redirects",

		// Commands - DNS
		"dns.short":   "DNS lookup utilities",
		"dns.long":    `Query DNS records for domains.`,
		"dns.lookup":  "Lookup DNS records",
		"dns.type":    "Record type (A, AAAA, CNAME, MX, TXT, NS)",
		"dns.server":  "DNS server to use",

		// Commands - IP
		"ip.short":    "IP lookup utilities",
		"ip.long":     `Query IP address and domain information.`,
		"ip.lookup":   "Lookup IP information",
		"ip.local":    "Show local IP addresses",
		"ip.public":   "Show public IP address",

		// Commands - Docker
		"docker.short": "Docker command generator",
		"docker.long":  `Generate Docker commands for common operations.`,
		"docker.run":   "Generate docker run command",
		"docker.build": "Generate docker build command",
		"docker.ps":    "Generate docker ps command",
		"docker.logs":  "Generate docker logs command",
		"docker.image": "Docker image",
		"docker.name":  "Container name",
		"docker.ports": "Port mappings",
		"docker.volumes": "Volume mounts",
		"docker.env":   "Environment variables",
		"docker.detach": "Run in detached mode",

		// Commands - Git
		"git.short":    "Git command generator",
		"git.long":     `Generate Git commands for common operations.`,
		"git.commit":   "Generate git commit command",
		"git.push":     "Generate git push command",
		"git.pull":     "Generate git pull command",
		"git.branch":   "Generate git branch command",
		"git.merge":    "Generate git merge command",
		"git.rebase":   "Generate git rebase command",
		"git.clone":    "Generate git clone command",
		"git.message":  "Commit message",
		"git.remote":   "Remote name",
		"git.branchName": "Branch name",

		// Commands - AI Chat
		"aichat.short":      "AI chat client",
		"aichat.long":       `Chat with AI models (OpenAI, Anthropic, etc.).`,
		"aichat.send":       "Send message to AI",
		"aichat.interactive": "Interactive chat session",
		"aichat.models":     "List available models",
		"aichat.apiKey":     "API key",
		"aichat.baseUrl":    "API base URL",
		"aichat.model":      "Model name",
		"aichat.stream":     "Stream response",

		// Commands - ASCII
		"ascii.short": "ASCII art generator",
		"ascii.long":  `Convert text to ASCII art with various fonts.`,
		"ascii.text":  "Convert text to ASCII art",
		"ascii.list":  "List available fonts",
		"ascii.font":  "Font style",

		// Commands - IBAN
		"iban.short":     "IBAN validator and formatter",
		"iban.long":      `Validate and format IBAN numbers.`,
		"iban.validate":  "Validate IBAN",
		"iban.format":    "Format IBAN",
		"iban.valid":     "Valid IBAN",
		"iban.invalid":   "Invalid IBAN",

		// Commands - OTP
		"otp.short":     "OTP generator",
		"otp.long":      `Generate TOTP and HOTP codes.`,
		"otp.totp":      "Generate TOTP code",
		"otp.hotp":      "Generate HOTP code",
		"otp.verify":    "Verify TOTP code",
		"otp.secret":    "Secret key",
		"otp.digits":    "Number of digits",
		"otp.period":    "Time period (seconds)",
		"otp.algorithm": "Hash algorithm",

		// Commands - RSA
		"rsa.short":     "RSA key pair generator",
		"rsa.long":      `Generate RSA public/private key pairs.`,
		"rsa.generate":  "Generate RSA key pair",
		"rsa.info":      "Show RSA key information",
		"rsa.bits":      "Key size in bits",
		"rsa.format":    "Key format",
		"rsa.privateKey": "Private key",
		"rsa.publicKey": "Public key",

		// Commands - Mock
		"mock.short":   "Mock data generator",
		"mock.long":    `Generate mock data for testing.`,
		"mock.user":    "Generate mock user data",
		"mock.employee": "Generate mock employee data",
		"mock.address": "Generate mock address data",
		"mock.product": "Generate mock product data",
		"mock.count":   "Number of records",

		// Commands - Roman
		"roman.short":      "Roman numeral converter",
		"roman.long":       `Convert between Roman numerals and Arabic numbers.`,
		"roman.toArabic":   "Convert Roman to Arabic",
		"roman.fromArabic": "Convert Arabic to Roman",
		"roman.convert":    "Auto-detect and convert",

		// Commands - Chmod
		"chmod.short":    "File permission calculator",
		"chmod.long":     `Calculate and convert file permissions.`,
		"chmod.calc":     "Calculate permissions",
		"chmod.octal":    "Convert octal to symbolic",
		"chmod.symbolic": "Convert symbolic to octal",
		"chmod.owner":    "Owner permissions",
		"chmod.group":    "Group permissions",
		"chmod.other":    "Other permissions",
		"chmod.read":     "Read permission",
		"chmod.write":    "Write permission",
		"chmod.execute":  "Execute permission",

		// Commands - HTTP Status
		"httpstatus.short":  "HTTP status code reference",
		"httpstatus.long":   `Look up HTTP status codes.`,
		"httpstatus.lookup": "Look up status code",
		"httpstatus.list":   "List all status codes",
		"httpstatus.search": "Search status codes",
		"httpstatus.code":   "Status code",
		"httpstatus.name":   "Status name",
		"httpstatus.category": "Category",

		// Commands - MIME
		"mime.short":    "MIME type lookup",
		"mime.long":     `Look up MIME types by file extension.`,
		"mime.lookup":   "Look up MIME type",
		"mime.list":     "List all MIME types",
		"mime.extension": "File extension",
		"mime.type":     "MIME type",

		// Commands - User Agent
		"useragent.short": "User-Agent parser",
		"useragent.long":  `Parse and analyze User-Agent strings.`,
		"useragent.parse": "Parse User-Agent",
		"useragent.list":  "List common User-Agents",
		"useragent.browser": "Browser",
		"useragent.os":    "Operating System",
		"useragent.device": "Device",

		// Commands - Basic Auth
		"basicauth.short":    "Basic Auth generator",
		"basicauth.long":     `Generate Basic Authentication headers.`,
		"basicauth.generate": "Generate Basic Auth",
		"basicauth.decode":   "Decode Basic Auth",
		"basicauth.username": "Username",
		"basicauth.password": "Password",
		"basicauth.header":   "Authorization header",

		// Commands - Password Strength
		"passwordstrength.short": "Password strength checker",
		"passwordstrength.long":  `Analyze password strength.`,
		"passwordstrength.score": "Strength score",
		"passwordstrength.level": "Strength level",
		"passwordstrength.weak":  "Weak",
		"passwordstrength.fair":  "Fair",
		"passwordstrength.strong": "Strong",
		"passwordstrength.veryStrong": "Very Strong",

		// Commands - Coordinate
		"coordinate.short":    "Coordinate converter",
		"coordinate.long":     `Convert between coordinate formats.`,
		"coordinate.toDMS":    "Convert to DMS",
		"coordinate.toDecimal": "Convert to decimal",
		"coordinate.distance": "Calculate distance",
		"coordinate.latitude": "Latitude",
		"coordinate.longitude": "Longitude",

		// Commands - Env
		"env.short":      "Environment variable manager",
		"env.long":       `Manage environment variables from .env files.`,
		"env.load":       "Load .env file",
		"env.get":        "Get variable",
		"env.set":        "Set variable",
		"env.list":       "List variables",
		"env.validate":   "Validate .env file",
		"env.file":       "Env file path",

		// Commands - Vault
		"vault.short":      "Password vault manager",
		"vault.long":       `Securely store passwords with encryption.`,
		"vault.init":       "Initialize vault",
		"vault.list":       "List entries",
		"vault.add":        "Add entry",
		"vault.get":        "Get entry",
		"vault.remove":     "Remove entry",
		"vault.password":   "Master password",
		"vault.title":      "Entry title",
		"vault.category":   "Category",

		// Commands - Image
		"image.short":      "Image processing utilities",
		"image.long":       `Process and convert images.`,
		"image.info":       "Show image info",
		"image.convert":    "Convert image format",
		"image.base64":     "Image to base64",
		"image.base64decode": "Base64 to image",
		"image.format":     "Image format",
		"image.dimensions": "Dimensions",

		// Commands - PDF
		"pdf.short":        "PDF utilities",
		"pdf.long":         `PDF signature verification and info.`,
		"pdf.info":         "Show PDF info",
		"pdf.sign":         "Show signature info",

		// Errors
		"error.invalidInput":    "Invalid input: %s",
		"error.fileNotFound":    "File not found: %s",
		"error.readFile":        "Failed to read file: %s",
		"error.writeFile":       "Failed to write file: %s",
		"error.decode":          "Failed to decode: %s",
		"error.encode":          "Failed to encode: %s",
		"error.invalidFormat":   "Invalid format: %s",
		"error.required":        "Required: %s",
		"error.notFound":        "Not found: %s",
		"error.alreadyExists":   "Already exists: %s",
		"error.permission":      "Permission denied: %s",
		"error.network":         "Network error: %s",
		"error.timeout":         "Timeout: %s",
		"error.unknown":         "Unknown error: %s",
	}

	// Chinese translations
	t.translations["zh"] = map[string]string{
		// App
		"app.name":        "Kairoa",
		"app.description": "桌面实用工具 CLI",

		// Common
		"common.error":           "错误",
		"common.success":         "成功",
		"common.warning":         "警告",
		"common.info":            "信息",
		"common.done":            "完成",
		"common.cancel":          "取消",
		"common.save":            "保存",
		"common.load":            "加载",
		"common.delete":          "删除",
		"common.edit":            "编辑",
		"common.add":             "添加",
		"common.remove":          "移除",
		"common.search":          "搜索",
		"common.copy":            "复制",
		"common.paste":           "粘贴",
		"common.copied":          "已复制到剪贴板",
		"common.notFound":        "未找到",
		"common.invalid":         "无效",
		"common.required":        "必填",
		"common.optional":        "可选",
		"common.enabled":         "已启用",
		"common.disabled":        "已禁用",
		"common.yes":             "是",
		"common.no":              "否",
		"common.true":            "真",
		"common.false":           "假",
		"common.input":           "输入",
		"common.output":          "输出",
		"common.result":          "结果",
		"common.format":          "格式",
		"common.type":            "类型",
		"common.value":           "值",
		"common.key":             "键",
		"common.name":            "名称",
		"common.description":     "描述",
		"common.example":         "示例",
		"common.usage":           "用法",
		"common.help":            "帮助",
		"common.version":         "版本",
		"common.author":          "作者",
		"common.license":         "许可证",

		// Commands - Hash
		"hash.short": "计算哈希值",
		"hash.long":  `使用各种算法（MD5、SHA1、SHA256、SHA384、SHA512、RIPEMD160）计算文本或文件的哈希值。`,
		"hash.text":  "计算文本字符串的哈希值",
		"hash.file":  "计算文件的哈希值",
		"hash.md5":   "计算 MD5 哈希",
		"hash.sha1":  "计算 SHA1 哈希",
		"hash.sha256": "计算 SHA256 哈希",
		"hash.sha512": "计算 SHA512 哈希",
		"hash.algorithm": "要使用的哈希算法",
		"hash.bcrypt": "使用 Bcrypt 哈希",
		"hash.verify": "验证哈希",

		// Commands - UUID
		"uuid.short":    "生成 UUID",
		"uuid.long":     `生成 UUID（v4、v7）和 ULID。`,
		"uuid.v4":       "生成 UUID v4（随机）",
		"uuid.v7":       "生成 UUID v7（时间排序）",
		"uuid.ulid":     "生成 ULID",
		"uuid.count":    "生成数量",
		"uuid.uppercase": "大写输出",

		// Commands - Base64
		"base64.short": "Base64 编码/解码",
		"base64.long":  `对字符串进行 Base64 编码和解码。`,
		"base64.encode": "编码为 Base64",
		"base64.decode": "从 Base64 解码",
		"base64.url":    "使用 URL 安全编码",

		// Commands - JSON
		"json.short":    "JSON 格式化工具",
		"json.long":     `格式化、验证和查询 JSON 数据。`,
		"json.format":   "格式化 JSON",
		"json.minify":   "压缩 JSON",
		"json.validate": "验证 JSON",
		"json.query":    "使用 jq 语法查询 JSON",
		"json.sort":     "排序 JSON 键",

		// Commands - URL
		"url.short":    "URL 编码/解码",
		"url.long":     `对 URL 字符串进行编码和解码。`,
		"url.encode":   "编码 URL",
		"url.decode":   "解码 URL",
		"url.full":     "编码完整 URL",
		"url.component": "编码 URL 组件",

		// Commands - JWT
		"jwt.short":   "解码 JWT 令牌",
		"jwt.long":    `解码和验证 JWT 令牌。`,
		"jwt.decode":  "解码 JWT 令牌",
		"jwt.header":  "JWT 头部",
		"jwt.payload": "JWT 载荷",
		"jwt.signature": "JWT 签名",

		// Commands - Time
		"time.short":    "时间工具",
		"time.long":     `在时间戳和日期之间转换。`,
		"time.now":      "显示当前时间",
		"time.convert":  "转换时间戳为日期",
		"time.timestamp": "获取当前时间戳",
		"time.format":   "时间格式",
		"time.timezone": "时区",
		"time.unix":     "Unix 时间戳（秒）",
		"time.unixMs":   "Unix 时间戳（毫秒）",

		// Commands - Password
		"password.short":     "密码生成器",
		"password.long":      `生成安全的随机密码。`,
		"password.length":    "密码长度",
		"password.uppercase": "包含大写字母",
		"password.lowercase": "包含小写字母",
		"password.numbers":   "包含数字",
		"password.special":   "包含特殊字符",
		"password.strength":  "密码强度",

		// Commands - QR Code
		"qr.short":    "生成二维码",
		"qr.long":     `从文本或 URL 生成二维码。`,
		"qr.text":     "要编码的文本",
		"qr.size":     "二维码大小",
		"qr.level":    "纠错级别",
		"qr.terminal": "在终端显示",

		// Commands - HTTP
		"http.short":   "HTTP 客户端工具",
		"http.long":    `发送 HTTP 请求并检查响应。`,
		"http.get":     "发送 GET 请求",
		"http.post":    "发送 POST 请求",
		"http.put":     "发送 PUT 请求",
		"http.delete":  "发送 DELETE 请求",
		"http.headers": "请求头",
		"http.body":    "请求体",
		"http.follow":  "跟随重定向",

		// Commands - DNS
		"dns.short":   "DNS 查询工具",
		"dns.long":    `查询域名的 DNS 记录。`,
		"dns.lookup":  "查询 DNS 记录",
		"dns.type":    "记录类型（A、AAAA、CNAME、MX、TXT、NS）",
		"dns.server":  "要使用的 DNS 服务器",

		// Commands - IP
		"ip.short":    "IP 查询工具",
		"ip.long":     `查询 IP 地址和域名信息。`,
		"ip.lookup":   "查询 IP 信息",
		"ip.local":    "显示本地 IP 地址",
		"ip.public":   "显示公网 IP 地址",

		// Commands - Docker
		"docker.short": "Docker 命令生成器",
		"docker.long":  `生成常用操作的 Docker 命令。`,
		"docker.run":   "生成 docker run 命令",
		"docker.build": "生成 docker build 命令",
		"docker.ps":    "生成 docker ps 命令",
		"docker.logs":  "生成 docker logs 命令",
		"docker.image": "Docker 镜像",
		"docker.name":  "容器名称",
		"docker.ports": "端口映射",
		"docker.volumes": "卷挂载",
		"docker.env":   "环境变量",
		"docker.detach": "后台运行",

		// Commands - Git
		"git.short":    "Git 命令生成器",
		"git.long":     `生成常用操作的 Git 命令。`,
		"git.commit":   "生成 git commit 命令",
		"git.push":     "生成 git push 命令",
		"git.pull":     "生成 git pull 命令",
		"git.branch":   "生成 git branch 命令",
		"git.merge":    "生成 git merge 命令",
		"git.rebase":   "生成 git rebase 命令",
		"git.clone":    "生成 git clone 命令",
		"git.message":  "提交信息",
		"git.remote":   "远程名称",
		"git.branchName": "分支名称",

		// Commands - AI Chat
		"aichat.short":      "AI 聊天客户端",
		"aichat.long":       `与 AI 模型聊天（OpenAI、Anthropic 等）。`,
		"aichat.send":       "发送消息给 AI",
		"aichat.interactive": "交互式聊天会话",
		"aichat.models":     "列出可用模型",
		"aichat.apiKey":     "API 密钥",
		"aichat.baseUrl":    "API 基础 URL",
		"aichat.model":      "模型名称",
		"aichat.stream":     "流式响应",

		// Commands - ASCII
		"ascii.short": "ASCII 艺术生成器",
		"ascii.long":  `使用各种字体将文本转换为 ASCII 艺术。`,
		"ascii.text":  "将文本转换为 ASCII 艺术",
		"ascii.list":  "列出可用字体",
		"ascii.font":  "字体样式",

		// Commands - IBAN
		"iban.short":     "IBAN 验证器和格式化工具",
		"iban.long":      `验证和格式化 IBAN 号码。`,
		"iban.validate":  "验证 IBAN",
		"iban.format":    "格式化 IBAN",
		"iban.valid":     "有效的 IBAN",
		"iban.invalid":   "无效的 IBAN",

		// Commands - OTP
		"otp.short":     "OTP 生成器",
		"otp.long":      `生成 TOTP 和 HOTP 验证码。`,
		"otp.totp":      "生成 TOTP 验证码",
		"otp.hotp":      "生成 HOTP 验证码",
		"otp.verify":    "验证 TOTP 验证码",
		"otp.secret":    "密钥",
		"otp.digits":    "位数",
		"otp.period":    "时间周期（秒）",
		"otp.algorithm": "哈希算法",

		// Commands - RSA
		"rsa.short":     "RSA 密钥对生成器",
		"rsa.long":      `生成 RSA 公钥/私钥对。`,
		"rsa.generate":  "生成 RSA 密钥对",
		"rsa.info":      "显示 RSA 密钥信息",
		"rsa.bits":      "密钥大小（位）",
		"rsa.format":    "密钥格式",
		"rsa.privateKey": "私钥",
		"rsa.publicKey": "公钥",

		// Commands - Mock
		"mock.short":   "模拟数据生成器",
		"mock.long":    `生成用于测试的模拟数据。`,
		"mock.user":    "生成模拟用户数据",
		"mock.employee": "生成模拟员工数据",
		"mock.address": "生成模拟地址数据",
		"mock.product": "生成模拟产品数据",
		"mock.count":   "记录数量",

		// Commands - Roman
		"roman.short":      "罗马数字转换器",
		"roman.long":       `在罗马数字和阿拉伯数字之间转换。`,
		"roman.toArabic":   "罗马数字转阿拉伯数字",
		"roman.fromArabic": "阿拉伯数字转罗马数字",
		"roman.convert":    "自动检测并转换",

		// Commands - Chmod
		"chmod.short":    "文件权限计算器",
		"chmod.long":     `计算和转换文件权限。`,
		"chmod.calc":     "计算权限",
		"chmod.octal":    "八进制转符号表示",
		"chmod.symbolic": "符号表示转八进制",
		"chmod.owner":    "所有者权限",
		"chmod.group":    "组权限",
		"chmod.other":    "其他用户权限",
		"chmod.read":     "读取权限",
		"chmod.write":    "写入权限",
		"chmod.execute":  "执行权限",

		// Commands - HTTP Status
		"httpstatus.short":  "HTTP 状态码参考",
		"httpstatus.long":   `查询 HTTP 状态码。`,
		"httpstatus.lookup": "查询状态码",
		"httpstatus.list":   "列出所有状态码",
		"httpstatus.search": "搜索状态码",
		"httpstatus.code":   "状态码",
		"httpstatus.name":   "状态名称",
		"httpstatus.category": "类别",

		// Commands - MIME
		"mime.short":    "MIME 类型查询",
		"mime.long":     `按文件扩展名查询 MIME 类型。`,
		"mime.lookup":   "查询 MIME 类型",
		"mime.list":     "列出所有 MIME 类型",
		"mime.extension": "文件扩展名",
		"mime.type":     "MIME 类型",

		// Commands - User Agent
		"useragent.short": "User-Agent 解析器",
		"useragent.long":  `解析和分析 User-Agent 字符串。`,
		"useragent.parse": "解析 User-Agent",
		"useragent.list":  "列出常见 User-Agents",
		"useragent.browser": "浏览器",
		"useragent.os":    "操作系统",
		"useragent.device": "设备",

		// Commands - Basic Auth
		"basicauth.short":    "Basic Auth 生成器",
		"basicauth.long":     `生成 Basic Authentication 头部。`,
		"basicauth.generate": "生成 Basic Auth",
		"basicauth.decode":   "解码 Basic Auth",
		"basicauth.username": "用户名",
		"basicauth.password": "密码",
		"basicauth.header":   "Authorization 头部",

		// Commands - Password Strength
		"passwordstrength.short": "密码强度检查器",
		"passwordstrength.long":  `分析密码强度。`,
		"passwordstrength.score": "强度分数",
		"passwordstrength.level": "强度等级",
		"passwordstrength.weak":  "弱",
		"passwordstrength.fair":  "一般",
		"passwordstrength.strong": "强",
		"passwordstrength.veryStrong": "非常强",

		// Commands - Coordinate
		"coordinate.short":    "坐标转换器",
		"coordinate.long":     `在坐标格式之间转换。`,
		"coordinate.toDMS":    "转换为度分秒",
		"coordinate.toDecimal": "转换为十进制",
		"coordinate.distance": "计算距离",
		"coordinate.latitude": "纬度",
		"coordinate.longitude": "经度",

		// Commands - Env
		"env.short":      "环境变量管理器",
		"env.long":       `管理 .env 文件中的环境变量。`,
		"env.load":       "加载 .env 文件",
		"env.get":        "获取变量",
		"env.set":        "设置变量",
		"env.list":       "列出变量",
		"env.validate":   "验证 .env 文件",
		"env.file":       "环境文件路径",

		// Commands - Vault
		"vault.short":      "密码保险箱管理器",
		"vault.long":       `使用加密安全存储密码。`,
		"vault.init":       "初始化保险箱",
		"vault.list":       "列出条目",
		"vault.add":        "添加条目",
		"vault.get":        "获取条目",
		"vault.remove":     "删除条目",
		"vault.password":   "主密码",
		"vault.title":      "条目标题",
		"vault.category":   "类别",

		// Commands - Image
		"image.short":      "图像处理工具",
		"image.long":       `处理和转换图像。`,
		"image.info":       "显示图像信息",
		"image.convert":    "转换图像格式",
		"image.base64":     "图像转 base64",
		"image.base64decode": "base64 转图像",
		"image.format":     "图像格式",
		"image.dimensions": "尺寸",

		// Commands - PDF
		"pdf.short":        "PDF 工具",
		"pdf.long":         `PDF 签名验证和信息。`,
		"pdf.info":         "显示 PDF 信息",
		"pdf.sign":         "显示签名信息",

		// Errors
		"error.invalidInput":    "无效输入：%s",
		"error.fileNotFound":    "文件未找到：%s",
		"error.readFile":        "读取文件失败：%s",
		"error.writeFile":       "写入文件失败：%s",
		"error.decode":          "解码失败：%s",
		"error.encode":          "编码失败：%s",
		"error.invalidFormat":   "无效格式：%s",
		"error.required":        "必填：%s",
		"error.notFound":        "未找到：%s",
		"error.alreadyExists":   "已存在：%s",
		"error.permission":      "权限被拒绝：%s",
		"error.network":         "网络错误：%s",
		"error.timeout":         "超时：%s",
		"error.unknown":         "未知错误：%s",
	}
}

// Global translator instance
var globalTranslator = NewTranslator()

// SetLocale sets the global locale
func SetLocale(locale Locale) {
	globalTranslator.SetLocale(locale)
}

// GetLocale gets the global locale
func GetLocale() Locale {
	return globalTranslator.GetLocale()
}

// T translates a key using the global translator
func T(key string, args ...interface{}) string {
	return globalTranslator.T(key, args...)
}

// DetectLocale detects locale from environment
func DetectLocale() Locale {
	// Check LANG environment variable
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}

	// Check for Chinese
	if strings.Contains(lang, "zh") || strings.Contains(lang, "CN") {
		return Chinese
	}

	// Default to English
	return English
}

// Init initializes the translator with detected locale
func Init() {
	SetLocale(DetectLocale())
}
