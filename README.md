# BTC 定投机器人

基于 AHR999 指标的智能 BTC 定投系统，支持多种消息推送方式。

## 功能特性

- 🤖 基于 AHR999 指标的智能定投策略
- ⏰ 支持 cron 表达式的定时调度
- 📱 多平台消息推送（Telegram、飞书、微信、Lark）
- 🔧 灵活的配置管理
- 🛡️ 调试模式支持
- 🚀 Docker 容器化部署

## 项目结构

```
regular_input/
├── ahr999/                # AHR999 指标计算
│   ├── ahr999.go
│   └── calculate_amount.go
├── auto_buy/              # 定投任务核心逻辑
│   └── auto_buy.go
├── config/                # 配置管理
│   ├── config.go
│   └── env.go
├── exchange_api/          # 交易所 API 客户端
│   ├── client.go
│   └── example.go
├── helper/                # 消息推送助手
│   ├── lark_bot.go
│   ├── messages.go
│   ├── telegram_bot.go
│   ├── test_pushers.go
│   └── wechat_bot.go
├── scheduler/             # 定时任务调度器
│   └── scheduler.go
├── main.go                # 程序入口
├── config.json.example    # 配置文件模板
├── Dockerfile             # Docker 构建文件
└── README.md              # 项目说明
```

## 快速开始

### 1. 环境准备

确保已安装 Go 1.21+ 和 Docker（可选）。

### 2. 配置设置

复制配置文件模板：

```bash
cp config.json.example config.json
```

编辑 `config.json` 文件：

```json
{
  "task_config": {
    "name": "regular-buy",
    "schedule": "0 0 7 * * *",
    "log_level": "info"
  },
  "params_config": {
    "debug": true,
    "base_amount": 100,
    "use_ahr999": true,
    "ahr999_timer_table": {
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
    "enabled": true,
    "push_method": ["telegram", "lark", "feishu", "wechat"]
  }
}
```

### 3. 环境变量配置

创建 `.env` 文件并配置必要的环境变量：

```bash
# 币安 API 配置（生产环境必需）
BINANCE_API_KEY=your_api_key
BINANCE_SECRET_KEY=your_secret_key

# 消息推送配置（可选）
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_CHAT_ID=your_chat_id
LARK_TOKEN=your_lark_token
FEISHU_TOKEN=your_feishu_token
WEIXIN_FT_TOKEN=your_wechat_token

# 代理配置（可选）
HTTPS_PROXY=http://proxy:port

# 调试模式
DEBUG=true
```

### 4. 运行程序

#### 方式一：使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 快速开始
make config    # 创建配置文件
make env       # 创建环境变量文件
make dev       # 开发模式运行

# Docker 部署
make build-docker  # 构建镜像
make run          # 运行容器
make logs         # 查看日志
```

#### 方式二：直接运行

```bash
# 安装依赖
go mod tidy

# 运行程序
go run main.go
```

#### 方式三：Docker 手动运行

```bash
# 构建镜像
docker build -t btc-dca-bot .

# 运行容器
docker run -d \
  --name btc-dca-bot \
  -v $(pwd)/config.json:/app/config.json \
  -e BINANCE_API_KEY=your_api_key \
  -e BINANCE_SECRET_KEY=your_secret_key \
  -e TELEGRAM_BOT_TOKEN=your_bot_token \
  -e TELEGRAM_CHAT_ID=your_chat_id \
  btc-dca-bot
```

## 配置说明

### 任务配置 (task_config)

- `name`: 任务名称
- `schedule`: cron 定时表达式（支持秒级）
- `log_level`: 日志级别

### 参数配置 (params_config)

- `debug`: 调试模式（true=模拟交易，false=真实交易）
- `base_amount`: 基础定投金额（USDT）
- `use_ahr999`: 是否启用 AHR999 指标
- `ahr999_timer_table`: AHR999 倍数表

### 消息配置 (message_config)

- `enabled`: 是否启用消息推送
- `push_method`: 推送方式列表

## AHR999 定投策略

AHR999 是一个用于判断比特币投资时机的指标：

- **< 0.45**: 极度低估，8倍定投
- **0.45-0.6**: 低估，4倍定投
- **0.6-0.8**: 较低估，2倍定投
- **0.8-0.9**: 略低估，1倍定投
- **0.9-1.1**: 正常，0.5倍定投
- **1.1-1.2**: 略高估，0.25倍定投
- **1.2-1.4**: 高估，0.125倍定投
- **1.4-1.6**: 较高估，暂停定投
- **1.6-1.8**: 高估，暂停定投
- **> 1.8**: 极度高估，暂停定投

## 消息推送

支持多种消息推送方式：

- **Telegram**: 需要 `TELEGRAM_BOT_TOKEN` 和 `TELEGRAM_CHAT_ID`
- **飞书**: 需要 `FEISHU_TOKEN`
- **Lark**: 需要 `LARK_TOKEN`
- **微信**: 需要 `WEIXIN_FT_TOKEN`

## 安全注意事项

1. **API 密钥安全**: 请妥善保管币安 API 密钥，建议设置 IP 白名单
2. **权限控制**: API 密钥只需要现货交易权限，不要开启提币权限
3. **测试模式**: 建议先在调试模式下测试，确认无误后再切换到生产模式
4. **资金安全**: 建议使用小额资金进行测试

## Makefile 命令说明

本项目提供了丰富的 Makefile 命令来简化开发、构建和部署流程：

### 📦 构建相关

```bash
make build        # 构建 Go 二进制文件
make build-docker # 构建 Docker 镜像
make deps         # 安装 Go 依赖
make clean        # 清理构建文件
```

### 🐳 Docker 相关

```bash
make run          # 运行 Docker 容器
make stop         # 停止 Docker 容器
make logs         # 查看容器日志
make shell        # 进入容器 shell
make status       # 查看容器状态
```

### 🔧 开发相关

```bash
make dev          # 开发模式运行
make test         # 运行测试
make lint         # 代码检查
make fmt          # 格式化代码
```

### ⚙️ 配置相关

```bash
make config       # 创建配置文件
make env          # 创建环境变量文件
make validate     # 验证配置
```

### 📊 监控相关

```bash
make health       # 健康检查
make stats        # 查看统计信息
```

### 🚀 快速开始示例

```bash
# 1. 克隆项目
git clone <repository-url>
cd regular_input

# 2. 创建配置文件
make config
make env

# 3. 编辑配置文件
vim config.json
vim .env

# 4. 开发模式运行
make dev

# 或者 Docker 部署
make build-docker
make run
make logs
```

## 技术栈

- **语言**: Go 1.21+
- **调度引擎**: robfig/cron/v3
- **交易所 API**: go-binance
- **消息推送**: 多平台支持
- **容器化**: Docker
- **构建工具**: Make

## 许可证

MIT License 