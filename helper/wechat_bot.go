package helper

import (
	"fmt"
	"time"

	serverchan "github.com/easychen/serverchan-sdk-golang"
)

// WeChatPusher 微信推送器
type WeChatPusher struct {
	sendKey string
}

// NewWeChatPusher 创建微信推送器
func NewWeChatPusher(sendKey string) *WeChatPusher {
	return &WeChatPusher{
		sendKey: sendKey,
	}
}

// Push 实现 MessagePusher 接口
func (w *WeChatPusher) Push(message string) error {
	// 发送消息
	resp, err := serverchan.ScSend(w.sendKey, "定投通知", message, nil)
	if err != nil {
		return fmt.Errorf("微信推送失败: %w", err)
	}

	// 检查响应
	if resp != nil && resp.Code != 0 {
		return fmt.Errorf("微信推送失败: %s", resp.Message)
	}

	return nil
}

// PushMessage 推送消息（旧接口，保持兼容性）
func (w *WeChatPusher) PushMessage(msg Message) error {
	// 构建消息内容
	content := BuildMessageContent(msg)

	// 发送消息
	resp, err := serverchan.ScSend(w.sendKey, msg.Title, content, nil)
	if err != nil {
		return fmt.Errorf("微信推送失败: %w", err)
	}

	// 检查响应
	if resp != nil && resp.Code != 0 {
		return fmt.Errorf("微信推送失败: %s", resp.Message)
	}

	return nil
}

// HealthCheck 健康检查
func (w *WeChatPusher) TestPush() bool {
	// 发送测试消息进行健康检查
	testMsg := NewMessage("每日定投", "测试微信推送器")
	content := BuildMessageContent(*testMsg)
	err := w.Push(content)
	return err == nil
}

// BuildMessageContent 构建消息内容
func BuildMessageContent(msg Message) string {
	content := msg.Content

	// 添加时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 构建完整内容
	fullContent := fmt.Sprintf("【%s】\n", msg.Title)
	fullContent += fmt.Sprintf("时间: %s\n\n", timestamp)
	fullContent += fmt.Sprintf("来源: %s\n\n", msg.AppID)
	fullContent += fmt.Sprintf("消息ID: %s\n\n", msg.ID)
	fullContent += fmt.Sprintf("内容: \n\n\n%s\n\n", content)

	// 添加元数据
	if len(msg.Metadata) > 0 {
		fullContent += "\n【元数据】\n"
		for key, value := range msg.Metadata {
			fullContent += fmt.Sprintf("%s: %v\n", key, value)
		}
	}

	return fullContent
}

// SetSendKey 设置sendKey
func (w *WeChatPusher) SetSendKey(sendKey string) {
	w.sendKey = sendKey
}

// GetSendKey 获取sendKey
func (w *WeChatPusher) GetSendKey() string {
	return w.sendKey
}
