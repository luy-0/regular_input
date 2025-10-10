# Makefile for BTC DCA Bot

# 变量定义
BINARY_NAME = regular_input
IMAGE_NAME = btc-dca-bot
TAG = latest
CONTAINER_NAME = btc-dca-bot
GO_VERSION = 1.21.6
BUILD_TIME = $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION = 1.0.0

# 构建标志
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 默认目标
.PHONY: help
help:
	@echo "🤖 BTC 定投机器人 - 可用命令:"
	@echo ""
	@echo "📦 构建相关:"
	@echo "  make build        - 构建 Go 二进制文件"
	@echo "  make build-docker - 构建 Docker 镜像"
	@echo "  make clean        - 清理构建文件"
	@echo ""
	@echo "🐳 Docker 相关:"
	@echo "  make run          - 运行 Docker 容器"
	@echo "  make stop         - 停止 Docker 容器"
	@echo "  make logs         - 查看容器日志"
	@echo "  make shell        - 进入容器 shell"
	@echo "  make status       - 查看容器状态"
	@echo ""
	@echo "🔧 开发相关:"
	@echo "  make dev          - 开发模式运行"
	@echo "  make test         - 运行测试"
	@echo "  make lint         - 代码检查"
	@echo "  make fmt          - 格式化代码"
	@echo "  make deps         - 安装依赖"
	@echo ""
	@echo "⚙️  配置相关:"
	@echo "  make config       - 创建配置文件"
	@echo "  make env          - 创建环境变量文件"
	@echo "  make validate     - 验证配置"
	@echo ""
	@echo "📊 监控相关:"
	@echo "  make health       - 健康检查"
	@echo "  make stats        - 查看统计信息"

# ==================== 构建相关 ====================

# 构建 Go 二进制文件
.PHONY: build
build: deps
	@echo "🔨 构建 Go 二进制文件..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME) .
	@echo "✅ 构建完成: bin/$(BINARY_NAME)"

# 构建 Docker 镜像
.PHONY: build-docker
build-docker:
	@echo "🐳 构建 Docker 镜像..."
	docker build -t $(IMAGE_NAME):$(TAG) .
	@echo "✅ Docker 镜像构建完成: $(IMAGE_NAME):$(TAG)"

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 安装 Go 依赖..."
	go mod tidy
	go mod download
	@echo "✅ 依赖安装完成"

# 清理构建文件
.PHONY: clean
clean:
	@echo "🧹 清理构建文件..."
	rm -rf bin/
	docker stop $(CONTAINER_NAME) 2>/dev/null || true
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	docker rmi $(IMAGE_NAME):$(TAG) 2>/dev/null || true
	@echo "✅ 清理完成"

# ==================== Docker 相关 ====================

# 运行 Docker 容器
.PHONY: run
run: build-docker
	@echo "🚀 运行 Docker 容器..."
	@if [ ! -f .env ]; then \
		echo "⚠️  .env 文件不存在，正在创建..."; \
		make env; \
	fi
	@if [ ! -f config.json ]; then \
		echo "⚠️  config.json 文件不存在，正在创建..."; \
		cp config.json.example config.json; \
	fi
	docker run -d \
		--name $(CONTAINER_NAME) \
		--restart unless-stopped \
		-v "$$(pwd)/config.json:/app/config.json:ro" \
		-v "$$(pwd)/.env:/app/.env:ro" \
		-e TZ=Asia/Shanghai \
		$(IMAGE_NAME):$(TAG)
	@echo "✅ 容器已启动: $(CONTAINER_NAME)"

# 停止 Docker 容器
.PHONY: stop
stop:
	@echo "🛑 停止 Docker 容器..."
	docker stop $(CONTAINER_NAME) 2>/dev/null || true
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@echo "✅ 容器已停止并删除"

# 查看容器日志
.PHONY: logs
logs:
	@echo "📋 查看容器日志..."
	docker logs -f $(CONTAINER_NAME)

# 进入容器 shell
.PHONY: shell
shell:
	@echo "🐚 进入容器 shell..."
	docker exec -it $(CONTAINER_NAME) sh

# 查看容器状态
.PHONY: status
status:
	@echo "📊 容器状态:"
	docker ps -a --filter name=$(CONTAINER_NAME)

# ==================== 开发相关 ====================

# 开发模式运行
.PHONY: dev
dev: config
	@echo "🔧 开发模式运行..."
	@echo "⚠️  请确保已设置环境变量或创建 .env 文件"
	go run main.go

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 代码检查
.PHONY: lint
lint:
	@echo "🔍 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过代码检查"; \
		echo "   安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 格式化代码
.PHONY: fmt
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	@echo "✅ 代码格式化完成"

# ==================== 配置相关 ====================

# 创建配置文件
.PHONY: config
config:
	@echo "⚙️  创建配置文件..."
	@if [ ! -f config.json ]; then \
		cp config.json.example config.json; \
		echo "✅ 已创建 config.json 文件"; \
		echo "📝 请编辑 config.json 文件配置您的参数"; \
	else \
		echo "✅ config.json 文件已存在"; \
	fi

# 创建环境变量文件
.PHONY: env
env:
	@echo "🔐 创建环境变量文件..."
	@if [ ! -f .env ]; then \
		echo "# BTC 定投机器人环境变量配置" > .env; \
		echo "" >> .env; \
		echo "# 币安 API 配置（生产环境必需）" >> .env; \
		echo "BINANCE_API_KEY=your_api_key" >> .env; \
		echo "BINANCE_SECRET_KEY=your_secret_key" >> .env; \
		echo "" >> .env; \
		echo "# 消息推送配置（可选）" >> .env; \
		echo "TELEGRAM_BOT_TOKEN=your_telegram_bot_token" >> .env; \
		echo "TELEGRAM_CHAT_ID=your_chat_id" >> .env; \
		echo "LARK_TOKEN=your_lark_token" >> .env; \
		echo "FEISHU_TOKEN=your_feishu_token" >> .env; \
		echo "WEIXIN_FT_TOKEN=your_wechat_token" >> .env; \
		echo "" >> .env; \
		echo "# 代理配置（可选）" >> .env; \
		echo "HTTPS_PROXY=http://proxy:port" >> .env; \
		echo "" >> .env; \
		echo "# 调试模式" >> .env; \
		echo "DEBUG=true" >> .env; \
		echo "✅ 已创建 .env 文件"; \
		echo "📝 请编辑 .env 文件填入您的实际配置"; \
	else \
		echo "✅ .env 文件已存在"; \
	fi

# 验证配置
.PHONY: validate
validate:
	@echo "🔍 验证配置..."
	@if [ -f config.json ]; then \
		echo "✅ config.json 文件存在"; \
		python3 -m json.tool config.json > /dev/null && echo "✅ config.json 格式正确" || echo "❌ config.json 格式错误"; \
	else \
		echo "❌ config.json 文件不存在，请运行 make config"; \
	fi
	@if [ -f .env ]; then \
		echo "✅ .env 文件存在"; \
	else \
		echo "⚠️  .env 文件不存在，请运行 make env"; \
	fi

# ==================== 监控相关 ====================

# 健康检查
.PHONY: health
health:
	@echo "🏥 健康检查..."
	@if docker ps --filter name=$(CONTAINER_NAME) --filter status=running | grep -q $(CONTAINER_NAME); then \
		echo "✅ 容器运行正常"; \
		docker exec $(CONTAINER_NAME) ps aux | grep regular_input || echo "⚠️  进程可能未运行"; \
	else \
		echo "❌ 容器未运行"; \
	fi

# 查看统计信息
.PHONY: stats
stats:
	@echo "📊 统计信息..."
	@echo "容器状态:"
	docker stats --no-stream $(CONTAINER_NAME) 2>/dev/null || echo "容器未运行"
	@echo ""
	@echo "镜像信息:"
	docker images $(IMAGE_NAME) 2>/dev/null || echo "镜像不存在"

# ==================== 其他命令 ====================

# 重新构建并启动
.PHONY: rebuild
rebuild: clean build-docker run
	@echo "🔄 重新构建并启动完成"

# 快速重启
.PHONY: restart
restart: stop run
	@echo "🔄 快速重启完成"

# 查看镜像
.PHONY: images
images:
	@echo "🐳 Docker 镜像:"
	docker images | grep $(IMAGE_NAME) || echo "镜像不存在"

# 查看所有容器
.PHONY: ps
ps:
	@echo "📋 所有容器:"
	docker ps -a

# 显示版本信息
.PHONY: version
version:
	@echo "📋 版本信息:"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Git 提交: $(GIT_COMMIT)"
	@echo "  Go 版本: $(GO_VERSION)"

# 修复配置文件问题
.PHONY: fix-config
fix-config:
	@echo "🔧 修复配置文件问题..."
	@if [ -d config.json ]; then \
		echo "❌ 发现 config.json 是目录，正在删除..."; \
		rm -rf config.json; \
	fi
	@if [ -d .env ]; then \
		echo "❌ 发现 .env 是目录，正在删除..."; \
		rm -rf .env; \
	fi
	@echo "📝 重新创建配置文件..."
	@make config
	@make env
	@echo "✅ 配置文件修复完成"

# 检查配置文件状态
.PHONY: check-config
check-config:
	@echo "🔍 检查配置文件状态..."
	@echo "config.json:"
	@if [ -f config.json ]; then \
		echo "  ✅ 文件存在"; \
		ls -la config.json; \
	elif [ -d config.json ]; then \
		echo "  ❌ 是目录，需要修复"; \
		ls -la config.json; \
	else \
		echo "  ⚠️  文件不存在"; \
	fi
	@echo ""
	@echo ".env:"
	@if [ -f .env ]; then \
		echo "  ✅ 文件存在"; \
		ls -la .env; \
	elif [ -d .env ]; then \
		echo "  ❌ 是目录，需要修复"; \
		ls -la .env; \
	else \
		echo "  ⚠️  文件不存在"; \
	fi

# 调试容器内文件状态
.PHONY: debug-container
debug-container:
	@echo "🐛 调试容器内文件状态..."
	@if docker ps --filter name=$(CONTAINER_NAME) --filter status=running | grep -q $(CONTAINER_NAME); then \
		echo "容器运行中，检查文件状态..."; \
		echo "工作目录:"; \
		docker exec $(CONTAINER_NAME) pwd; \
		echo ""; \
		echo "文件列表:"; \
		docker exec $(CONTAINER_NAME) ls -la /app/; \
		echo ""; \
		echo "config.json 内容:"; \
		docker exec $(CONTAINER_NAME) cat /app/config.json 2>/dev/null || echo "无法读取 config.json"; \
		echo ""; \
		echo ".env 内容:"; \
		docker exec $(CONTAINER_NAME) cat /app/.env 2>/dev/null || echo "无法读取 .env"; \
	else \
		echo "❌ 容器未运行，请先运行 make run"; \
	fi 