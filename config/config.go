package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 总配置结构
type Config struct {
	TaskConfig    TaskConfig    `json:"task_config"`
	ParamsConfig  ParamsConfig  `json:"params_config"`
	MessageConfig MessageConfig `json:"message_config"`
}

// TaskConfig 任务配置
type TaskConfig struct {
	Name     string `json:"name"`      // 任务名称
	Schedule string `json:"schedule"`  // cron 定时表达式
	LogLevel string `json:"log_level"` // 日志级别
}

// ParamsConfig 参数配置
type ParamsConfig struct {
	Debug            bool               `json:"debug"`              // 是否调试模式
	BaseAmount       float64            `json:"base_amount"`        // 基础金额
	UseAhr999        bool               `json:"use_ahr999"`         // 是否使用 AHR999 指标
	Ahr999TimerTable map[string]float64 `json:"ahr999_timer_table"` // AHR999 倍数表
}

// MessageConfig 消息推送配置
type MessageConfig struct {
	Enabled    bool     `json:"enabled"`     // 是否启用消息推送
	PushMethod []string `json:"push_method"` // 推送方式列表
}

// LoadConfig 从文件加载配置
func LoadConfig(filepath string) (*Config, error) {
	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// LoadConfigOrDefault 加载配置，如果失败则返回默认配置
func LoadConfigOrDefault(filepath string) *Config {
	config, err := LoadConfig(filepath)
	if err != nil {
		fmt.Printf("警告: 加载配置失败，使用默认配置: %v\n", err)
		return GetDefaultConfig()
	}
	return config
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		TaskConfig: TaskConfig{
			Name:     "regular-buy",
			Schedule: "0 0 7 * * *",
			LogLevel: "info",
		},
		ParamsConfig: ParamsConfig{
			Debug:      false,
			BaseAmount: 100,
			UseAhr999:  true,
			Ahr999TimerTable: map[string]float64{
				"<0.45":    8,
				"0.45-0.6": 4,
				"0.6-0.8":  2,
				"0.8-0.9":  1,
				"0.9-1.1":  0.5,
				"1.1-1.2":  0.25,
				"1.2-1.4":  0.125,
				"1.4-1.6":  0,
				"1.6-1.8":  0,
				">1.8":     0,
			},
		},
		MessageConfig: MessageConfig{
			Enabled:    true,
			PushMethod: []string{"telegram", "lark", "feishu", "wechat"},
		},
	}
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	// 验证任务名称
	if c.TaskConfig.Name == "" {
		return fmt.Errorf("任务名称不能为空")
	}

	// 验证 cron 表达式
	if c.TaskConfig.Schedule == "" {
		return fmt.Errorf("定时表达式不能为空")
	}

	// 验证基础金额
	if c.ParamsConfig.BaseAmount <= 0 {
		return fmt.Errorf("基础金额必须大于 0")
	}

	// 验证推送方式
	if c.MessageConfig.Enabled && len(c.MessageConfig.PushMethod) == 0 {
		return fmt.Errorf("启用消息推送时必须指定至少一个推送方式")
	}

	return nil
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig(filepath string) error {
	// 将配置序列化为 JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetPushMethods 获取启用的推送方式
func (c *Config) GetPushMethods() []string {
	if !c.MessageConfig.Enabled {
		return []string{}
	}
	return c.MessageConfig.PushMethod
}

// IsPushMethodEnabled 检查指定的推送方式是否启用
func (c *Config) IsPushMethodEnabled(method string) bool {
	if !c.MessageConfig.Enabled {
		return false
	}
	for _, m := range c.MessageConfig.PushMethod {
		if m == method {
			return true
		}
	}
	return false
}

// GetAhr999Multiplier 根据 AHR999 值获取对应的倍数
func (c *Config) GetAhr999Multiplier(ahr999Value float64) float64 {
	if !c.ParamsConfig.UseAhr999 {
		return 1.0
	}

	// 根据范围查找对应的倍数
	switch {
	case ahr999Value < 0.45:
		return c.ParamsConfig.Ahr999TimerTable["<0.45"]
	case ahr999Value >= 0.45 && ahr999Value < 0.6:
		return c.ParamsConfig.Ahr999TimerTable["0.45-0.6"]
	case ahr999Value >= 0.6 && ahr999Value < 0.8:
		return c.ParamsConfig.Ahr999TimerTable["0.6-0.8"]
	case ahr999Value >= 0.8 && ahr999Value < 0.9:
		return c.ParamsConfig.Ahr999TimerTable["0.8-0.9"]
	case ahr999Value >= 0.9 && ahr999Value < 1.1:
		return c.ParamsConfig.Ahr999TimerTable["0.9-1.1"]
	case ahr999Value >= 1.1 && ahr999Value < 1.2:
		return c.ParamsConfig.Ahr999TimerTable["1.1-1.2"]
	case ahr999Value >= 1.2 && ahr999Value < 1.4:
		return c.ParamsConfig.Ahr999TimerTable["1.2-1.4"]
	case ahr999Value >= 1.4 && ahr999Value < 1.6:
		return c.ParamsConfig.Ahr999TimerTable["1.4-1.6"]
	case ahr999Value >= 1.6 && ahr999Value < 1.8:
		return c.ParamsConfig.Ahr999TimerTable["1.6-1.8"]
	default: // >= 1.8
		return c.ParamsConfig.Ahr999TimerTable[">1.8"]
	}
}

// CalculateAmount 计算实际投资金额
func (c *Config) CalculateAmount(ahr999Value float64) float64 {
	multiplier := c.GetAhr999Multiplier(ahr999Value)
	return c.ParamsConfig.BaseAmount * multiplier
}

// String 返回配置的字符串表示
func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}
