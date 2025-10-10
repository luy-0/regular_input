package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvConfig 环境变量配置
type EnvConfig struct {
	// 交易所 API
	BinanceAPIKey    string
	BinanceSecretKey string

	// 消息推送
	TelegramBotToken string
	TelegramChatID   string
	LarkToken        string
	FeishuToken      string
	WeixinFTToken    string

	// 网络代理
	HTTPSProxy string
	HTTPProxy  string

	// 其他配置
	LogLevel string
	Debug    bool
}

// LoadEnvConfig 加载环境变量配置
func LoadEnvConfig() *EnvConfig {
	// 首先尝试从 .env 文件加载
	loadEnvFile()

	return &EnvConfig{
		// 交易所 API
		BinanceAPIKey:    os.Getenv("BINANCE_API_KEY"),
		BinanceSecretKey: os.Getenv("BINANCE_SECRET_KEY"),

		// 消息推送
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		LarkToken:        os.Getenv("LARK_TOKEN"),
		FeishuToken:      os.Getenv("FEISHU_TOKEN"),
		WeixinFTToken:    os.Getenv("WEIXIN_FT_TOKEN"),

		// 网络代理（仅在设置了环境变量时使用）
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		HTTPProxy:  os.Getenv("HTTP_PROXY"),

		// 其他配置
		LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
		Debug:    getEnvOrDefault("DEBUG", "false") == "true",
	}
}

// loadEnvFile 从 .env 文件加载环境变量
func loadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		// .env 文件不存在，跳过
		return
	}

	file, err := os.Open(envFile)
	if err != nil {
		fmt.Printf("警告: 无法打开 .env 文件: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		// 设置环境变量（如果尚未设置）
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("警告: 读取 .env 文件时出错: %v\n", err)
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ValidateEnvConfig 验证环境变量配置
func (env *EnvConfig) ValidateEnvConfig(requiredPushMethods []string) error {
	var missingVars []string

	// 检查必需的环境变量
	if env.BinanceAPIKey == "" {
		missingVars = append(missingVars, "BINANCE_API_KEY")
	}
	if env.BinanceSecretKey == "" {
		missingVars = append(missingVars, "BINANCE_SECRET_KEY")
	}

	// 检查消息推送相关的环境变量
	for _, method := range requiredPushMethods {
		switch strings.ToLower(method) {
		case "telegram":
			if env.TelegramBotToken == "" {
				missingVars = append(missingVars, "TELEGRAM_BOT_TOKEN")
			}
			if env.TelegramChatID == "" {
				missingVars = append(missingVars, "TELEGRAM_CHAT_ID")
			}
		case "lark":
			if env.LarkToken == "" {
				missingVars = append(missingVars, "LARK_TOKEN")
			}
		case "feishu":
			if env.FeishuToken == "" {
				missingVars = append(missingVars, "FEISHU_TOKEN")
			}
		case "wechat":
			if env.WeixinFTToken == "" {
				missingVars = append(missingVars, "WEIXIN_FT_TOKEN")
			}
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("缺少必需的环境变量: %s", strings.Join(missingVars, ", "))
	}

	return nil
}

// GetMissingEnvVars 获取缺失的环境变量列表
func (env *EnvConfig) GetMissingEnvVars(requiredPushMethods []string) []string {
	var missingVars []string

	// 检查交易所 API
	if env.BinanceAPIKey == "" {
		missingVars = append(missingVars, "BINANCE_API_KEY")
	}
	if env.BinanceSecretKey == "" {
		missingVars = append(missingVars, "BINANCE_SECRET_KEY")
	}

	// 检查消息推送相关的环境变量
	for _, method := range requiredPushMethods {
		switch strings.ToLower(method) {
		case "telegram":
			if env.TelegramBotToken == "" {
				missingVars = append(missingVars, "TELEGRAM_BOT_TOKEN")
			}
			if env.TelegramChatID == "" {
				missingVars = append(missingVars, "TELEGRAM_CHAT_ID")
			}
		case "lark":
			if env.LarkToken == "" {
				missingVars = append(missingVars, "LARK_TOKEN")
			}
		case "feishu":
			if env.FeishuToken == "" {
				missingVars = append(missingVars, "FEISHU_TOKEN")
			}
		case "wechat":
			if env.WeixinFTToken == "" {
				missingVars = append(missingVars, "WEIXIN_FT_TOKEN")
			}
		}
	}

	return missingVars
}

// PrintEnvStatus 打印环境变量状态
func (env *EnvConfig) PrintEnvStatus(requiredPushMethods []string) {
	fmt.Println("=== 环境变量状态 ===")

	// 交易所 API
	fmt.Println("\n【交易所 API】")
	status := "❌ 未设置"
	if env.BinanceAPIKey != "" {
		status = "✅ 已设置"
	}
	fmt.Printf("  BINANCE_API_KEY: %s\n", status)

	status = "❌ 未设置"
	if env.BinanceSecretKey != "" {
		status = "✅ 已设置"
	}
	fmt.Printf("  BINANCE_SECRET_KEY: %s\n", status)

	// 消息推送
	fmt.Println("\n【消息推送】")

	for _, method := range requiredPushMethods {
		switch strings.ToLower(method) {
		case "telegram":
			botStatus := "❌ 未设置"
			if env.TelegramBotToken != "" {
				botStatus = "✅ 已设置"
			}
			chatStatus := "❌ 未设置"
			if env.TelegramChatID != "" {
				chatStatus = "✅ 已设置"
			}
			fmt.Printf("  TELEGRAM_BOT_TOKEN: %s\n", botStatus)
			fmt.Printf("  TELEGRAM_CHAT_ID: %s\n", chatStatus)
		case "lark":
			status := "❌ 未设置"
			if env.LarkToken != "" {
				status = "✅ 已设置"
			}
			fmt.Printf("  LARK_TOKEN: %s\n", status)
		case "feishu":
			status := "❌ 未设置"
			if env.FeishuToken != "" {
				status = "✅ 已设置"
			}
			fmt.Printf("  FEISHU_TOKEN: %s\n", status)
		case "wechat":
			status := "❌ 未设置"
			if env.WeixinFTToken != "" {
				status = "✅ 已设置"
			}
			fmt.Printf("  WEIXIN_FT_TOKEN: %s\n", status)
		}
	}

	// 网络代理
	fmt.Println("\n【网络代理】")
	if env.HTTPSProxy != "" {
		fmt.Printf("  HTTPS_PROXY: ✅ 已设置 (%s)\n", (env.HTTPSProxy))
	} else {
		fmt.Println("  HTTPS_PROXY: ❌ 未设置（将使用直连模式）")
	}

	if env.HTTPProxy != "" {
		fmt.Printf("  HTTP_PROXY: ✅ 已设置 (%s)\n", (env.HTTPProxy))
	} else {
		fmt.Println("  HTTP_PROXY: ❌ 未设置（将使用直连模式）")
	}

	// 其他配置
	fmt.Println("\n【其他配置】")
	fmt.Printf("  LOG_LEVEL: %s\n", env.LogLevel)
	fmt.Printf("  DEBUG: %v\n", env.Debug)
}

// CheckEnvFile 检查是否存在 .env 文件
func CheckEnvFile() bool {
	if _, err := os.Stat(".env"); err == nil {
		return true
	}
	return false
}

// PrintEnvSetupInstructions 打印环境变量设置说明
func PrintEnvSetupInstructions() {
	fmt.Println("\n=== 环境变量设置说明 ===")
	fmt.Println("1. 复制 env.example 文件为 .env：")
	fmt.Println("   cp env.example .env")
	fmt.Println("\n2. 编辑 .env 文件，填入实际的环境变量值")
	fmt.Println("\n3. 或者直接设置环境变量：")
	fmt.Println("   export BINANCE_API_KEY=your_api_key")
	fmt.Println("   export BINANCE_SECRET_KEY=your_secret_key")
	fmt.Println("   # ... 其他环境变量")
	fmt.Println("\n4. 运行程序：")
	fmt.Println("   ./regular_input")
}
