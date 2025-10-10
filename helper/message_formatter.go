package helper

import (
	"fmt"
	"strings"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeSuccess MessageType = "success"
	MessageTypeError   MessageType = "error"
	MessageTypeInfo    MessageType = "info"
	MessageTypeDebug   MessageType = "debug"
)

// FormattedMessage 格式化消息结构
type FormattedMessage struct {
	Type      MessageType
	Title     string
	Content   string
	Timestamp time.Time
	Metadata  map[string]interface{}
	Emoji     string
	Color     string
}

// NewFormattedMessage 创建格式化消息
func NewFormattedMessage(msgType MessageType, title, content string) *FormattedMessage {
	return &FormattedMessage{
		Type:      msgType,
		Title:     title,
		Content:   content,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
		Emoji:     getEmojiByType(msgType),
		Color:     getColorByType(msgType),
	}
}

// getEmojiByType 根据消息类型获取表情符号
func getEmojiByType(msgType MessageType) string {
	switch msgType {
	case MessageTypeSuccess:
		return "✅"
	case MessageTypeError:
		return "❌"
	case MessageTypeInfo:
		return "ℹ️"
	case MessageTypeDebug:
		return "🧪"
	default:
		return "📢"
	}
}

// getColorByType 根据消息类型获取颜色（用于支持颜色的平台）
func getColorByType(msgType MessageType) string {
	switch msgType {
	case MessageTypeSuccess:
		return "green"
	case MessageTypeError:
		return "red"
	case MessageTypeInfo:
		return "blue"
	case MessageTypeDebug:
		return "yellow"
	default:
		return "default"
	}
}

// FormatDCAReport 格式化定投报告消息
func FormatDCAReport(isDebug bool, btcPrice, ahr999Value, amount float64, orderResult string, err error) *FormattedMessage {
	var msgType MessageType
	var title string
	var content strings.Builder

	// 确定消息类型和标题
	if isDebug {
		msgType = MessageTypeDebug
		title = "定投任务（调试模式）"
	} else if err != nil {
		msgType = MessageTypeError
		title = "定投任务失败"
	} else {
		msgType = MessageTypeSuccess
		title = "定投任务成功"
	}

	// 构建内容
	content.WriteString(fmt.Sprintf("💰 **BTC 价格**: $%.2f\n", btcPrice))

	if ahr999Value > 0 {
		content.WriteString(fmt.Sprintf("📊 **AHR999 值**: %.4f\n", ahr999Value))

		// 添加AHR999区间说明
		interval := getAHR999Interval(ahr999Value)
		content.WriteString(fmt.Sprintf("📈 **投资区间**: %s\n", interval))
	}

	content.WriteString(fmt.Sprintf("💵 **定投金额**: %.2f USDT\n", amount))

	if !isDebug && orderResult != "" {
		content.WriteString("\n📋 **订单详情**:\n")
		content.WriteString(formatOrderResult(orderResult))
	}

	if err != nil {
		content.WriteString(fmt.Sprintf("\n⚠️ **错误信息**: %v", err))
	}

	msg := NewFormattedMessage(msgType, title, content.String())
	msg.Metadata["btc_price"] = btcPrice
	msg.Metadata["amount"] = amount
	msg.Metadata["is_debug"] = isDebug

	return msg
}

// getAHR999Interval 获取AHR999区间描述
func getAHR999Interval(value float64) string {
	switch {
	case value < 0.45:
		return "极度低估区间 🚀"
	case value >= 0.45 && value < 0.6:
		return "低估区间 📈"
	case value >= 0.6 && value < 0.8:
		return "较低估区间 📊"
	case value >= 0.8 && value < 0.9:
		return "略低估区间 📉"
	case value >= 0.9 && value < 1.1:
		return "正常区间 ⚖️"
	case value >= 1.1 && value < 1.2:
		return "略高估区间 ⚠️"
	case value >= 1.2 && value < 1.4:
		return "高估区间 ⛔"
	case value >= 1.4 && value < 1.6:
		return "较高估区间 🛑"
	case value >= 1.6 && value < 1.8:
		return "高估区间 🚫"
	default:
		return "极度高估区间 ⛔"
	}
}

// formatOrderResult 格式化订单结果
func formatOrderResult(result string) string {
	if result == "" {
		return "无订单信息"
	}

	// 简化订单结果显示
	if len(result) > 200 {
		return result[:200] + "..."
	}

	return result
}

// ToTelegramFormat 转换为Telegram格式
func (m *FormattedMessage) ToTelegramFormat() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s **%s**\n", m.Emoji, m.Title))
	sb.WriteString(fmt.Sprintf("🕐 %s\n\n", m.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(m.Content)

	if len(m.Metadata) > 0 {
		sb.WriteString("\n\n📊 **统计信息**:\n")
		for key, value := range m.Metadata {
			sb.WriteString(fmt.Sprintf("• %s: %v\n", key, value))
		}
	}

	return sb.String()
}

// ToLarkFormat 转换为Lark/飞书格式
func (m *FormattedMessage) ToLarkFormat() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s\n", m.Emoji, m.Title))
	sb.WriteString(fmt.Sprintf("时间: %s\n\n", m.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(m.Content)

	if len(m.Metadata) > 0 {
		sb.WriteString("\n\n统计信息:\n")
		for key, value := range m.Metadata {
			sb.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
		}
	}

	return sb.String()
}

// ToWeChatFormat 转换为微信格式
func (m *FormattedMessage) ToWeChatFormat() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("【%s】%s\n", m.Title, m.Emoji))
	sb.WriteString(fmt.Sprintf("时间: %s\n\n", m.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(m.Content)

	if len(m.Metadata) > 0 {
		sb.WriteString("\n\n统计信息:\n")
		for key, value := range m.Metadata {
			sb.WriteString(fmt.Sprintf("• %s: %v\n", key, value))
		}
	}

	return sb.String()
}

// ToPlainText 转换为纯文本格式
func (m *FormattedMessage) ToPlainText() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s\n", m.Emoji, m.Title))
	sb.WriteString(fmt.Sprintf("时间: %s\n\n", m.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(m.Content)

	if len(m.Metadata) > 0 {
		sb.WriteString("\n\n统计信息:\n")
		for key, value := range m.Metadata {
			sb.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
		}
	}

	return sb.String()
}
