package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LarkBot struct for sending messages via Lark (Feishu) Webhook API
type LarkBot struct {
	isUseFeishu bool
	Token       string
	BaseURL     string
	WebhookURL  string
	client      *http.Client
}

const (
	FeiShuBaseURL = "https://www.feishu.cn/flow/api/trigger-webhook/"
	LarkBaseURL   = "https://open.larksuite.com/open-apis/bot/v2/hook/"
)

// NewLarkBot creates a new LarkBot instance
func NewLarkBot(token string, isUseFeishu ...bool) *LarkBot {
	base := LarkBaseURL
	if len(isUseFeishu) > 0 && isUseFeishu[0] {
		base = FeiShuBaseURL
	}

	webhookURL := fmt.Sprintf("%s/%s", trimSuffix(base, "/"), token)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &LarkBot{
		Token:      token,
		BaseURL:    base,
		WebhookURL: webhookURL,
		client:     client,
	}
}

// SendText sends a text message to Lark
func (bot *LarkBot) SendText(title, content string) (map[string]interface{}, error) {
	text := fmt.Sprintf("【%s】\n\n%s", title, content)
	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": text,
		},
	}
	return bot.sendMessage(payload)
}

// Push 实现 MessagePusher 接口
func (bot *LarkBot) Push(message string) error {
	_, err := bot.SendText("定投通知", message)
	return err
}

// PushFormatted 推送格式化消息
func (bot *LarkBot) PushFormatted(msg *FormattedMessage) error {
	formattedMessage := msg.ToLarkFormat()
	_, err := bot.SendText(msg.Title, formattedMessage)
	return err
}

// TestPush performs a health check by sending a test message
func (bot *LarkBot) TestPush() bool {
	// 发送测试消息进行健康检查
	str := "飞书"
	if bot.isUseFeishu {
		str = "飞书"
	} else {
		str = "Lark"
	}
	testMsg := NewMessage("每日定投", "测试"+str+"推送器")
	content := BuildMessageContent(*testMsg)
	_, err := bot.SendText(testMsg.Title, content)
	return err == nil
}

// sendMessage internal method to send API requests to Lark Webhook
func (bot *LarkBot) sendMessage(payload map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return map[string]interface{}{"code": -1, "error": err.Error()}, err
	}

	req, err := http.NewRequest("POST", bot.WebhookURL, bytes.NewBuffer(data))
	if err != nil {
		return map[string]interface{}{"code": -1, "error": err.Error()}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bot.client.Do(req)
	if err != nil {
		return map[string]interface{}{"code": -1, "error": err.Error()}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"code": -1, "error": err.Error()}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]interface{}{"code": -1, "error": err.Error()}, err
	}

	// 检查飞书 API 响应
	if code, exists := result["code"]; exists {
		if codeInt, ok := code.(float64); ok && codeInt != 0 {
			fmt.Printf("Lark send message failed: %+v\n", result)
		}
	}

	return result, nil
}

// trimSuffix removes suffix from string
func trimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}
