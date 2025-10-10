package helper

import (
	"fmt"
	"time"
)

// PushMethod 推送方式枚举
type PushMethod int

const (
	WeChat PushMethod = iota // 微信推送
	Email                    // 邮件推送
	SMS                      // 短信推送
	Logger                   // 日志推送
)

// Message 消息体定义
type Message struct {
	ID      string `json:"id"`      // 消息唯一标识，自动生成格式：{app_id}_YYMMDD_{gen_id}
	AppID   string `json:"app_id"`  // 发送方ID，标志消息来源
	Title   string `json:"title"`   // 消息标题
	Content string `json:"content"` // 消息内容
	// Level      MessageLevel           `json:"level"`       // 紧急程度
	Metadata  map[string]interface{} `json:"metadata"`   // 扩展元数据
	CreatedAt time.Time              `json:"created_at"` // 创建时间
	SentAt    time.Time              `json:"sent_at"`    // 最终成功发送时间
	// SendStatus SendStatus             `json:"send_status"` // 发送状态
	// mu sync.RWMutex `json:"-"` // 用于metadata操作的互斥锁
}

// NewMessage 创建新消息
func NewMessage(title, content string) *Message {
	return &Message{
		ID:        generateMessageID(title),
		AppID:     title,
		Title:     title,
		Content:   content,
		Metadata:  make(map[string]interface{}),
		CreatedAt: time.Now(),
	}
}

// SetMetadata 设置元数据
func (m *Message) SetMetadata(key string, value interface{}) *Message {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	m.Metadata[key] = value
	return m
}

// SetMetadata 设置元数据
func (m *Message) SetMetadataMap(metadata map[string]interface{}) *Message {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	if metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	for key, value := range metadata {
		m.Metadata[key] = value
	}
	return m
}

// GetMetadata 获取元数据
func (m *Message) GetMetadata(key string) (interface{}, bool) {
	// 	m.mu.RLock()
	// defer m.mu.RUnlock()
	if m.Metadata == nil {
		return nil, false
	}
	value, exists := m.Metadata[key]
	return value, exists
}

// generateMessageID 生成消息ID
func generateMessageID(title string) string {
	now := time.Now()
	dateStr := now.Format("060102_150405")                // YYMMDD格式
	nanoStr := fmt.Sprintf("%06d", now.Nanosecond()/1000) // 微秒部分

	return fmt.Sprintf("%s_%s_%s", title, dateStr, nanoStr)
}

// MessagePusher 消息推送器接口
type MessagePusher interface {
	Push(message string) error
	TestPush() bool
}

// FormattedMessagePusher 格式化消息推送器接口
type FormattedMessagePusher interface {
	MessagePusher
	PushFormatted(msg *FormattedMessage) error
}
