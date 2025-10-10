package auto_buy

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"regular_input/ahr999"
	"regular_input/exchange_api"
	"regular_input/helper"
)

// AutoBuyTask 定投任务
type AutoBuyTask struct {
	config      *Config
	exchangeAPI *exchange_api.Client
	pushers     []helper.MessagePusher
}

// Config 定投任务配置
type Config struct {
	Name             string             `json:"name"`               // 任务名称
	BaseAmount       float64            `json:"base_amount"`        // 基础金额
	UseAhr999        bool               `json:"use_ahr999"`         // 是否使用 AHR999 指标
	Ahr999TimerTable map[string]float64 `json:"ahr999_timer_table"` // AHR999 倍数表
	Symbol           string             `json:"symbol"`             // 交易对，如 BTCUSDT
	Debug            bool               `json:"debug"`              // 是否为调试模式
	PushMethods      []string           `json:"push_methods"`       // 推送方式
}

// NewAutoBuyTask 创建定投任务
func NewAutoBuyTask(config *Config, apiKey, secretKey, proxyUrl string) *AutoBuyTask {
	task := &AutoBuyTask{
		config:  config,
		pushers: []helper.MessagePusher{},
	}

	// 初始化交易所客户端
	if config.Debug {
		// 调试模式：不需要 API 密钥
		task.exchangeAPI = exchange_api.NewClientWithoutAuth(proxyUrl)
		log.Println("[定投任务] 调试模式：使用无认证客户端")
	} else {
		// 生产模式：需要 API 密钥
		task.exchangeAPI = exchange_api.NewClient(apiKey, secretKey, proxyUrl)
		log.Println("[定投任务] 生产模式：使用认证客户端")
	}

	// 初始化消息推送器
	task.initPushers()

	return task
}

// initPushers 初始化消息推送器
func (t *AutoBuyTask) initPushers() {
	for _, method := range t.config.PushMethods {
		switch method {
		case "telegram":
			token := os.Getenv("TELEGRAM_BOT_TOKEN")
			chatID := os.Getenv("TELEGRAM_CHAT_ID")
			if token != "" && chatID != "" {
				t.pushers = append(t.pushers, helper.NewTelegramBot(token, chatID))
				log.Println("[定投任务] 已启用 Telegram 推送")
			} else {
				log.Printf("[定投任务] Telegram 推送未配置 (需要 TELEGRAM_BOT_TOKEN 和 TELEGRAM_CHAT_ID)")
			}
		case "lark":
			token := os.Getenv("LARK_TOKEN")
			if token != "" {
				t.pushers = append(t.pushers, helper.NewLarkBot(token, false)) // Lark
				log.Println("[定投任务] 已启用 Lark 推送")
			} else {
				log.Printf("[定投任务] Lark 推送未配置 (需要 LARK_TOKEN)")
			}
		case "feishu":
			token := os.Getenv("FEISHU_TOKEN")
			if token != "" {
				t.pushers = append(t.pushers, helper.NewLarkBot(token, true)) // 飞书
				log.Println("[定投任务] 已启用飞书推送")
			} else {
				log.Printf("[定投任务] 飞书推送未配置 (需要 FEISHU_TOKEN)")
			}
		case "wechat":
			sendKey := os.Getenv("WEIXIN_FT_TOKEN")
			if sendKey != "" {
				t.pushers = append(t.pushers, helper.NewWeChatPusher(sendKey))
				log.Println("[定投任务] 已启用微信推送")
			} else {
				log.Printf("[定投任务] 微信推送未配置 (需要 WEIXIN_FT_TOKEN)")
			}
		}
	}

	if len(t.pushers) == 0 {
		log.Println("[定投任务] 未配置任何消息推送器")
	} else {
		log.Printf("[定投任务] 已配置 %d 个消息推送器", len(t.pushers))
	}
}

// Execute 执行定投任务
func (t *AutoBuyTask) Execute(ctx context.Context) error {
	log.Printf("[定投任务] 开始执行定投任务: %s", t.config.Name)

	var amount float64
	var ahr999Value float64
	var btcPrice float64
	var err error

	// 1. 获取当前 BTC 价格和 AHR999 值
	if t.config.UseAhr999 {
		btcPrice, ahr999Value, err = ahr999.GetAhr999()
		if err != nil {
			errMsg := fmt.Sprintf("获取 AHR999 数据失败: %v", err)
			log.Printf("[定投任务] %s", errMsg)
			errorMsg := helper.NewFormattedMessage(helper.MessageTypeError, "定投任务失败", errMsg)
			t.pushFormattedMessage(errorMsg)
			return fmt.Errorf(errMsg)
		}
		log.Printf("[定投任务] 当前 BTC 价格: %.2f, AHR999 值: %.4f", btcPrice, ahr999Value)

		// 2. 根据 AHR999 值计算定投金额
		timerTable := ahr999.Ahr999TimerTable(t.config.Ahr999TimerTable)
		amount, multiplier, rangeStr, err := ahr999.GetRecommendedAmount(
			t.config.BaseAmount,
			ahr999Value,
			timerTable,
		)
		if err != nil {
			errMsg := fmt.Sprintf("计算定投金额失败: %v", err)
			log.Printf("[定投任务] %s", errMsg)
			errorMsg := helper.NewFormattedMessage(helper.MessageTypeError, "定投任务失败", errMsg)
			t.pushFormattedMessage(errorMsg)
			return fmt.Errorf(errMsg)
		}
		log.Printf("[定投任务] AHR999 区间: %s, 倍数: %.2f, 定投金额: %.2f USDT",
			rangeStr, multiplier, amount)
	} else {
		// 不使用 AHR999，直接使用基础金额
		amount = t.config.BaseAmount
		btcPrice, err = t.exchangeAPI.GetBTCPrice(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("获取 BTC 价格失败: %v", err)
			log.Printf("[定投任务] %s", errMsg)
			errorMsg := helper.NewFormattedMessage(helper.MessageTypeError, "定投任务失败", errMsg)
			t.pushFormattedMessage(errorMsg)
			return fmt.Errorf(errMsg)
		}
		log.Printf("[定投任务] 当前 BTC 价格: %.2f, 定投金额: %.2f USDT", btcPrice, amount)
	}

	// 3. 执行买入操作
	if t.config.Debug {
		// 调试模式：只记录日志，不实际买入
		formattedMsg := helper.FormatDCAReport(true, btcPrice, ahr999Value, amount, "", nil)
		log.Printf("[定投任务] 调试模式：%s", formattedMsg.ToPlainText())
		t.pushFormattedMessage(formattedMsg)
		return nil
	}

	// 生产模式：实际买入
	symbol := t.config.Symbol
	if symbol == "" {
		symbol = "BTCUSDT" // 默认 BTC
	}

	orderResult := t.exchangeAPI.BuyCoinByMarketPrice(ctx, symbol, amount)
	log.Printf("[定投任务] 买入结果: %s", orderResult)

	// 4. 发送通知
	formattedMsg := helper.FormatDCAReport(false, btcPrice, ahr999Value, amount, orderResult, nil)
	t.pushFormattedMessage(formattedMsg)

	return nil
}

// GetName 获取任务名称
func (t *AutoBuyTask) GetName() string {
	return t.config.Name
}

// formatMessage 格式化消息
func (t *AutoBuyTask) formatMessage(isDebug bool, btcPrice, ahr999Value, amount float64, orderResult string, err error) string {
	now := time.Now().Format("2006-01-02 15:04:05")

	var message string
	if isDebug {
		message = fmt.Sprintf("🧪 定投任务（调试模式）\n时间: %s\n", now)
	} else {
		if err != nil {
			message = fmt.Sprintf("❌ 定投任务失败\n时间: %s\n错误: %v\n", now, err)
		} else {
			message = fmt.Sprintf("✅ 定投任务成功\n时间: %s\n", now)
		}
	}

	message += fmt.Sprintf("BTC 价格: $%.2f\n", btcPrice)

	if t.config.UseAhr999 {
		message += fmt.Sprintf("AHR999 值: %.4f\n", ahr999Value)
	}

	message += fmt.Sprintf("定投金额: %.2f USDT\n", amount)

	if !isDebug && orderResult != "" {
		message += fmt.Sprintf("订单详情:\n%s", orderResult)
	}

	return message
}

// pushMessage 推送消息
func (t *AutoBuyTask) pushMessage(message string) {
	if len(t.pushers) == 0 {
		return
	}

	for _, pusher := range t.pushers {
		if err := pusher.Push(message); err != nil {
			log.Printf("[定投任务] 消息推送失败: %v", err)
		}
	}
}

// pushFormattedMessage 推送格式化消息
func (t *AutoBuyTask) pushFormattedMessage(msg *helper.FormattedMessage) {
	if len(t.pushers) == 0 {
		return
	}

	for _, pusher := range t.pushers {
		var err error

		// 尝试使用格式化推送
		if formattedPusher, ok := pusher.(helper.FormattedMessagePusher); ok {
			err = formattedPusher.PushFormatted(msg)
		} else {
			// 回退到普通推送
			err = pusher.Push(msg.ToPlainText())
		}

		if err != nil {
			log.Printf("[定投任务] 消息推送失败: %v", err)
		}
	}
}
