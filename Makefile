# TRX委托服务 Makefile

# 变量定义
BINARY_NAME=lending-trx
BUILD_DIR=build
DOCKER_IMAGE=lending-trx

# 默认目标
.PHONY: all
all: build

# 构建所有服务
.PHONY: build
build:
	@echo "🔨 构建TRX委托服务..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/root/*.go
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 构建传统版本（根目录main.go）
.PHONY: build-legacy
build-legacy:
	@echo "🔨 构建传统版本..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/lending-trx-legacy main.go
	@echo "✅ 传统版本构建完成: $(BUILD_DIR)/lending-trx-legacy"

# 清理构建文件
.PHONY: clean
clean:
	@echo "🧹 清理构建文件..."
	rm -rf $(BUILD_DIR)
	@echo "✅ 清理完成"

# 运行完整服务（新架构）
.PHONY: server
server: build
	@echo "🚀 启动完整TRX委托服务..."
	./$(BUILD_DIR)/$(BINARY_NAME) server

# 运行Telegram Bot（新架构）
.PHONY: bot
bot: build
	@echo "🤖 启动Telegram Bot..."
	./$(BUILD_DIR)/$(BINARY_NAME) bot

# 运行传统版本
.PHONY: legacy
legacy: build-legacy
	@echo "🚀 启动传统版本服务..."
	./$(BUILD_DIR)/lending-trx-legacy

# 开发模式运行服务
.PHONY: dev
dev:
	@echo "🔧 开发模式运行服务..."
	go run main.go

# 开发模式运行新架构
.PHONY: dev-server
dev-server:
	@echo "🔧 开发模式运行新架构服务..."
	go run cmd/root/*.go server

# 开发模式运行Bot
.PHONY: dev-bot
dev-bot:
	@echo "🤖 开发模式运行Bot..."
	go run cmd/root/*.go bot

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	go test -v ./...
	@echo "✅ 测试完成"

# 运行Bot测试
.PHONY: test-bot
test-bot:
	@echo "🧪 运行Bot测试..."
	cd cmd/telegram_bot && go test -v
	@echo "✅ Bot测试完成"

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 安装依赖..."
	go mod download
	go mod tidy
	@echo "✅ 依赖安装完成"

# Docker构建
.PHONY: docker-build
docker-build:
	@echo "🐳 构建Docker镜像..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "✅ Docker镜像构建完成"

# Docker运行
.PHONY: docker-run
docker-run: docker-build
	@echo "🐳 运行Docker容器..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker Compose启动所有服务
.PHONY: docker-compose-up
docker-compose-up:
	@echo "🐳 启动所有服务 (数据库 + Server + Bot)..."
	docker-compose up -d

# Docker Compose停止所有服务
.PHONY: docker-compose-down
docker-compose-down:
	@echo "🐳 停止所有服务..."
	docker-compose down

# Docker Compose查看日志
.PHONY: docker-compose-logs
docker-compose-logs:
	@echo "📋 查看服务日志..."
	docker-compose logs -f

# Docker Compose重启服务
.PHONY: docker-compose-restart
docker-compose-restart:
	@echo "🔄 重启所有服务..."
	docker-compose restart

# 格式化代码
.PHONY: fmt
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	@echo "✅ 代码格式化完成"

# 代码检查
.PHONY: lint
lint:
	@echo "🔍 代码检查..."
	golangci-lint run
	@echo "✅ 代码检查完成"

# 显示帮助
.PHONY: help
help:
	@echo "🔧 TRX委托服务 Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  build        构建新架构服务"
	@echo "  build-legacy 构建传统版本"
	@echo "  clean        清理构建文件"
	@echo "  server       运行新架构完整服务"
	@echo "  bot          运行新架构Telegram Bot"
	@echo "  legacy       运行传统版本"
	@echo "  dev          开发模式运行传统版本"
	@echo "  dev-server   开发模式运行新架构服务"
	@echo "  dev-bot      开发模式运行Bot"
	@echo "  test         运行所有测试"
	@echo "  test-bot     运行Bot测试"
	@echo "  deps         安装依赖"
	@echo "  docker-build 构建Docker镜像"
	@echo "  docker-run   运行Docker容器"
	@echo "  docker-compose-up    启动所有服务 (数据库+Server+Bot)"
	@echo "  docker-compose-down  停止所有服务"
	@echo "  docker-compose-logs  查看服务日志"
	@echo "  docker-compose-restart 重启所有服务"
	@echo "  fmt          格式化代码"
	@echo "  lint         代码检查"
	@echo "  help         显示此帮助"
	@echo ""
	@echo "新架构命令示例:"
	@echo "  make server    # 启动完整服务"
	@echo "  make bot       # 启动Telegram Bot"
	@echo "  ./build/lending-trx server --port 9090  # 指定端口"
	@echo "  ./build/lending-trx bot --token YOUR_TOKEN  # 指定Token"
	@echo ""
	@echo "传统版本命令示例:"
	@echo "  make legacy    # 运行传统版本"
	@echo "  make dev       # 开发模式运行" 