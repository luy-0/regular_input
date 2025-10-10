# 多阶段构建 Dockerfile
# 第一阶段：构建阶段
FROM golang:1.21.6-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存层）
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o regular_input .

# 第二阶段：运行阶段
FROM alpine:3.18

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/regular_input .

# 复制配置文件模板
COPY --from=builder /app/config.json.example ./config.json.example

# 创建配置目录
RUN mkdir -p /app/config && chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口（如果需要）
# EXPOSE 8080

# 设置环境变量
ENV TZ=Asia/Shanghai

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep regular_input || exit 1

# 启动命令
CMD ["./regular_input"]
