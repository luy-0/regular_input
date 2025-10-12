# 基于 AHR999 指标的 BTC 定投机器人

基于 AHR999 指标的智能 BTC 定投系统，支持多种消息推送方式和灵活的配置管理。

## 🚀 功能特性

- 🤖 **智能定投策略**: 基于 AHR999 指标自动调整定投金额
- ⏰ **定时调度**: 支持 cron 表达式的秒级定时任务
- 📱 **多平台推送**: 支持 Telegram、飞书、Lark、微信等多种消息推送
- 🔧 **灵活配置**: JSON 配置文件 + 环境变量双重配置管理
- 🛡️ **安全模式**: 调试模式支持，避免误操作
- 🚀 **容器化部署**: 完整的 Docker 支持
- 📊 **实时监控**: 订单状态跟踪和消息通知
- 🔄 **自动重试**: 订单状态检查和自动重试机制

## 📁 项目结构

```
regular_input/
├── ahr999/                # AHR999 指标计算模块
│   ├── ahr999.go         # AHR999 数据获取和计算
│   └── calculate_amount.go # 定投金额计算逻辑
├── auto_buy/              # 定投任务核心逻辑
│   └── auto_buy.go       # 定投任务执行器
├── config/                # 配置管理模块
│   ├── config.go         # 配置结构和验证
│   └── env.go            # 环境变量管理
├── exchange_api/          # 交易所 API 客户端
│   ├── client.go         # 币安 API 封装
│   └── example.go        # API 使用示例
├── helper/                # 消息推送助手
│   ├── lark_bot.go       # Lark/飞书推送
│   ├── telegram_bot.go   # Telegram 推送
│   ├── wechat_bot.go     # 微信推送
│   ├── messages.go       # 消息格式定义
│   ├── message_formatter.go # 消息格式化
│   └── test_pushers.go   # 推送测试工具
├── scheduler/             # 定时任务调度器
│   └── scheduler.go      # Cron 任务调度
├── scripts/               # 部署脚本
│   ├── docker-build.sh   # Docker 构建脚本
│   ├── setup-env.sh      # 环境设置脚本
│   └── validate-dockerfile.sh # Dockerfile 验证
├── main.go                # 程序入口
├── config.json.example    # 配置文件模板
├── Dockerfile             # Docker 构建文件
├── Makefile              # 构建和部署脚本
└── README.md             # 项目说明
```

## 🚀 快速开始

### 1. 环境准备

确保已安装以下环境：
- Go 1.21+ 
- Docker（可选）
- Git

### 2. 克隆项目

```bash
git clone <repository-url>
cd regular_input
```

### 3. 配置设置

#### 方式一：使用 Makefile（推荐）

```bash
# 创建配置文件
make env

# 编辑配置文件
vim config.json
vim .env
```

#### 方式二：手动配置

```bash
# 复制配置文件模板
cp config.json.example config.json

# 创建环境变量文件
cp .env.example .env  # 如果存在
# 或手动创建 .env 文件
```

### 4. 配置文件说明

#### config.json 配置

```json
{
  "task_config": {
    "name": "regular-buy",
    "schedule": "0 0 7 * * *",
    "log_level": "info"
  },
  "params_config": {
    "debug": false,
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

#### .env 环境变量

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
```

### 5. 运行程序

#### 方式一：开发模式

```bash
# 安装依赖
make deps

# 开发模式运行
make dev
```

#### 方式二：Docker 部署

```bash
# 构建并运行
make build-docker
make run

# 查看日志
make logs
```

#### 方式三：直接运行

```bash
# 安装依赖
go mod tidy

# 运行程序
go run main.go
```

## 📊 AHR999 定投策略

AHR999 是一个用于判断比特币投资时机的指标，系统根据该指标自动调整定投金额：

| AHR999 区间 | 投资建议 | 定投倍数 | 说明 |
|------------|---------|---------|------|
| < 0.45 | 极度低估 | 8x | 大幅加仓 |
| 0.45-0.6 | 低估 | 4x | 适度加仓 |
| 0.6-0.8 | 较低估 | 2x | 正常定投 |
| 0.8-0.9 | 略低估 | 1x | 标准定投 |
| 0.9-1.1 | 正常 | 0.5x | 减半定投 |
| 1.1-1.2 | 略高估 | 0.25x | 少量定投 |
| 1.2-1.4 | 高估 | 0.125x | 微量定投 |
| 1.4-1.6 | 较高估 | 0x | 暂停定投 |
| 1.6-1.8 | 高估 | 0x | 暂停定投 |
| > 1.8 | 极度高估 | 0x | 暂停定投 |

## 📱 消息推送

支持多种消息推送方式，可同时配置多个：

### Telegram
- 需要：`TELEGRAM_BOT_TOKEN` 和 `TELEGRAM_CHAT_ID`
- 特点：实时推送，支持富文本格式

### 飞书/Lark
- 需要：`FEISHU_TOKEN` 或 `LARK_TOKEN`
- 特点：企业级消息推送

### 微信
- 需要：`WEIXIN_FT_TOKEN`
- 特点：通过 Server 酱推送

## 🛠️ Makefile 命令

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

## 🔒 安全注意事项

1. **API 密钥安全**
   - 妥善保管币安 API 密钥
   - 建议设置 IP 白名单
   - 只开启现货交易权限，不要开启提币权限

2. **测试模式**
   - 建议先在调试模式下测试
   - 确认无误后再切换到生产模式
   - 使用小额资金进行测试

3. **配置验证**
   - 使用 `make validate` 验证配置
   - 定期检查配置文件格式

## 🏗️ 技术架构

### 核心模块

- **调度器** (`scheduler/`): 基于 robfig/cron 的定时任务调度
- **定投引擎** (`auto_buy/`): 核心定投逻辑和订单管理
- **AHR999 计算** (`ahr999/`): 指标获取和金额计算
- **交易所 API** (`exchange_api/`): 币安 API 封装
- **消息推送** (`helper/`): 多平台消息推送支持
- **配置管理** (`config/`): 配置加载和验证

### 技术栈

- **语言**: Go 1.21+
- **调度引擎**: robfig/cron/v3
- **交易所 API**: go-binance
- **消息推送**: 多平台支持
- **容器化**: Docker
- **构建工具**: Make

## 📈 监控和日志

### 日志级别
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 监控指标
- 定投执行状态
- 订单成交情况
- AHR999 指标变化
- 系统健康状态

## 🐛 故障排除

### 常见问题

1. **配置文件错误**
   ```bash
   make validate  # 验证配置格式
   make fix-config  # 修复配置文件问题
   ```

2. **API 连接失败**
   - 检查网络连接
   - 验证 API 密钥
   - 检查代理设置

3. **消息推送失败**
   - 验证推送配置
   - 检查网络连接
   - 查看日志信息

### 调试模式

```bash
# 启用调试模式
export DEBUG=true

# 或修改 config.json
{
  "params_config": {
    "debug": true
  }
}
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如有问题，请通过以下方式联系：
- 提交 GitHub Issue
- 查看项目文档
- 检查日志信息