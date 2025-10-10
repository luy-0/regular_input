package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TelegramBot struct for sending messages via Telegram Bot API
type TelegramBot struct {
	Token   string
	ChatID  string
	BaseURL string
	client  *http.Client
	apiURL  string
}

const (
	BaseURL = "https://api.telegram.org/bot"
)

// NewTelegramBot creates a new TelegramBot instance
func NewTelegramBot(token, chatID string) *TelegramBot {
	apiURL := fmt.Sprintf("%s%s", BaseURL, token)
	client := &http.Client{
		Timeout: 30 * time.Second, // 增加超时时间到 30 秒
	}
	return &TelegramBot{
		Token:   token,
		ChatID:  chatID,
		BaseURL: BaseURL,
		client:  client,
		apiURL:  apiURL,
	}
}

// SendText sends a text message to Telegram
func (bot *TelegramBot) SendText(title, content string) (map[string]interface{}, error) {
	parseMode := "HTML"
	payload := map[string]interface{}{
		"chat_id":    bot.ChatID,
		"text":       fmt.Sprintf("【%s】\n\n%s", title, content),
		"parse_mode": parseMode,
	}
	return bot.sendMessage("sendMessage", payload)
}

func (bot *TelegramBot) TestPush() bool {
	// 发送测试消息进行健康检查
	testMsg := NewMessage("每日定投", "测试Telegram推送器")
	content := BuildMessageContent(*testMsg)
	result, err := bot.SendText(testMsg.Title, content)
	if err != nil {
		fmt.Printf("Telegram 错误详情: %v\n", err)
		if _, ok := err.(interface{ Timeout() bool }); ok {
			fmt.Println("提示: Telegram API 可能需要代理访问，请设置 HTTP_PROXY 或 HTTPS_PROXY 环境变量")
		}
		fmt.Printf("响应结果: %+v\n", result)
	}
	return err == nil
}

// sendMessage internal method to send API requests to Telegram
func (bot *TelegramBot) sendMessage(method string, payload map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", bot.apiURL, method)

	data, err := json.Marshal(payload)
	if err != nil {
		return map[string]interface{}{"ok": false, "error": err.Error()}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return map[string]interface{}{"ok": false, "error": err.Error()}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bot.client.Do(req)
	if err != nil {
		return map[string]interface{}{"ok": false, "error": err.Error()}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"ok": false, "error": err.Error()}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]interface{}{"ok": false, "error": err.Error()}, err
	}
	if ok, exists := result["ok"]; exists && !ok.(bool) {
		fmt.Printf("Telegram send message failed: %+v\n", result)
	}
	return result, nil
}
