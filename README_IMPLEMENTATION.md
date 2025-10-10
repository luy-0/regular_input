# 定投系统实现说明

## 项目结构

```
regular_input/
├── main.go                 # 主程序入口
├── config.go              # 配置管理
├── config.json            # 配置文件
├── test_debug.go          # 测试脚本
├── scheduler/             # 定时任务调度器
│   └── scheduler.go
├── auto_buy/              # 定投任务模块
│   └── auto_buy.go
├── ahr999/                # AHR999 指标计算
│   ├── ahr999.go
│   └── calculate_amount.go
├── exchange_api/          # 交易所 API
│   ├── client.go
│   └── example.go
└── helper/                # 消息推送
    ├── messages.go
    ├── telegram_bot.go
    ├── lark_bot.go
    └── wechat_bot.go
```

## 核心模块

### 1. 调度器模块 (scheduler/)

基于 `github.com/robfig/cron/v3` 实现的定时任务调度器：

- **Scheduler**: 主调度器，管理所有定时任务
- **Task**: 任务接口，所有任务都需要实现此接口
- **ScheduledTask**: 已调度的任务包装器

**主要功能：**
- 添加/移除任务
- 启动/停止调度器
- 立即执行任务（用于测试）
- 支持秒级精度的 cron 表达式

### 2. 定投任务模块 (auto_buy/)

整合了 AHR999 计算、金额计算和交易执行：

- **AutoBuyTask**: 定投任务实现
- **Config**: 定投任务配置

**主要功能：**
- 获取 AHR999 指标值
- 根据 AHR999 值计算定投金额
- 执行交易所买入操作
- 发送消息通知
- 支持调试模式（不实际买入）

### 3. AHR999 模块 (ahr999/)

负责 AHR999 指标的获取和计算：

- **GetAhr999()**: 获取当前 AHR999 值
- **GetAhr999At()**: 获取指定日期的 AHR999 值
- **CalculateAmount()**: 根据 AHR999 值计算定投金额
- **GetRecommendedAmount()**: 获取推荐的定投金额

**数据来源：**
- 优先使用本地缓存
- 缓存未命中时调用 API
- 按月存储历史数据

### 4. 交易所 API 模块 (exchange_api/)

基于 Binance API 的交易接口：

- **Client**: 交易所客户端
- **NewClient()**: 创建认证客户端
- **NewClientWithoutAuth()**: 创建无认证客户端（仅公开接口）

**主要功能：**
- 获取 BTC/ETH 价格
- 执行市价买入
- 获取账户余额
- 健康检查

### 5. 消息推送模块 (helper/)

支持多种消息推送方式：

- **TelegramBot**: Telegram 推送
- **LarkBot**: 飞书/Lark 推送
- **WeChatPusher**: 微信推送
- **MessagePusher**: 统一推送接口

## 配置说明

### config.json 配置项

```json
{
  "task_config": {
    "name": "regular-buy",           // 任务名称
    "schedule": "0 0 7 * * *",       // cron 表达式（每天7点）
    "log_level": "info"              // 日志级别
  },
  "params_config": {
    "debug": false,                  // 调试模式
    "base_amount": 100,              // 基础金额（USDT）
    "use_ahr999": true,              // 是否使用 AHR999
    "ahr999_timer_table": {          // AHR999 倍数表
      "<0.45": 8,
      "0.45-0.6": 4,
      "0.6-0.8": 2,
      "0.8-0.9": 1,
      "0.9-1.1": 0.5,
      "1.1-1.2": 0.25,
      "1.2-1.4": 0.125,
      "1.4-1.6": 0,
      "1.6-1.8": 0,
      ">1.8": 0
    }
  },
  "message_config": {
    "enabled": true,                 // 是否启用推送
    "push_method": [                 // 推送方式
      "telegram",
      "lark",
      "wechat"
    ]
  }
}
```

### 环境变量

```bash
# 交易所 API（生产模式需要）
export BINANCE_API_KEY="your_api_key"
export BINANCE_SECRET_KEY="your_secret_key"

# 消息推送
export TELEGRAM_BOT_TOKEN="your_bot_token"
export TELEGRAM_CHAT_ID="your_chat_id"
export LARK_TOKEN="your_lark_token"
export WECHAT_SEND_KEY="your_send_key"

# 网络代理（可选）
export HTTPS_PROXY="http://127.0.0.1:7890"
```

## 使用方法

### 1. 编译程序

```bash
go build -o regular_input main.go config.go
```

### 2. 运行程序

```bash
# 生产模式
./regular_input

# 调试模式（修改 config.json 中 debug 为 true）
./regular_input
```

### 3. 测试功能

```bash
# 测试配置
go run test_debug.go config

# 测试调试模式
go run test_debug.go debug

# 测试调度器（每30秒执行一次）
go run test_debug.go scheduler
```

## 工作流程

1. **启动阶段**：
   - 加载配置文件
   - 验证配置有效性
   - 初始化交易所客户端
   - 初始化消息推送器
   - 创建调度器并添加任务

2. **定时执行**：
   - 根据 cron 表达式触发任务
   - 获取当前 AHR999 值
   - 根据倍数表计算定投金额
   - 执行买入操作（或调试模式记录）
   - 发送结果通知

3. **错误处理**：
   - 网络错误重试
   - 配置错误提示
   - 交易失败通知

## 安全注意事项

1. **API 密钥安全**：
   - 不要将 API 密钥提交到代码仓库
   - 使用环境变量存储敏感信息
   - 定期轮换 API 密钥

2. **调试模式**：
   - 生产环境务必关闭调试模式
   - 调试模式不会实际执行交易

3. **网络代理**：
   - 确保代理服务器稳定可靠
   - 定期检查网络连接

## 故障排除

### 常见问题

1. **网络连接失败**：
   - 检查代理设置
   - 确认网络连接正常
   - 检查防火墙设置

2. **API 调用失败**：
   - 验证 API 密钥是否正确
   - 检查 API 权限设置
   - 确认账户余额充足

3. **消息推送失败**：
   - 检查推送服务配置
   - 验证 Token 和 Chat ID
   - 确认推送服务可用

### 日志分析

程序会输出详细的日志信息，包括：
- 配置加载状态
- 任务执行过程
- 错误信息和堆栈
- 交易结果详情

## 扩展功能

### 添加新的推送方式

1. 在 `helper/` 目录下创建新的推送器
2. 实现 `MessagePusher` 接口
3. 在 `auto_buy.go` 中添加初始化逻辑

### 添加新的交易所

1. 在 `exchange_api/` 目录下实现新的客户端
2. 实现统一的交易接口
3. 在 `auto_buy.go` 中添加交易所选择逻辑

### 添加新的指标

1. 创建新的指标计算模块
2. 在配置中添加指标相关设置
3. 在 `auto_buy.go` 中集成指标计算
