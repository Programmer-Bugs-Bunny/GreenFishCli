# GoWeb CLI Makefile

.PHONY: build install clean test run help

BINARY_NAME=goweb
BUILD_DIR=build
VERSION?=dev
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# 默认目标
all: build

# 构建二进制文件
build:
	@echo "🔨 构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 跨平台构建
build-all:
	@echo "🔨 构建多平台版本..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (M1/M2)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ 多平台构建完成"
	@ls -la $(BUILD_DIR)/

# 安装到系统路径
install: build
	@echo "📦 安装 $(BINARY_NAME)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ 安装完成"

# 本地安装到用户目录
install-local: build
	@echo "📦 本地安装 $(BINARY_NAME)..."
	@mkdir -p $(HOME)/.local/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/
	@echo "✅ 本地安装完成"
	@echo "💡 确保 $(HOME)/.local/bin 在您的 PATH 中"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 运行开发版本
run:
	@echo "🚀 运行开发版本..."
	go run . $(ARGS)

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "✅ 清理完成"

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	@echo "✅ 代码格式化完成"

# 代码检查
lint:
	@echo "🔍 代码检查..."
	go vet ./...
	@echo "✅ 代码检查完成"

# 更新依赖
tidy:
	@echo "📋 整理依赖..."
	go mod tidy
	@echo "✅ 依赖整理完成"

# 创建发布版本
release: clean fmt lint test build-all
	@echo "🚀 准备发布版本..."
	@echo "✅ 发布版本准备完成"

# 开发流程
dev: fmt lint test run

# 显示帮助信息
help:
	@echo "GoWeb CLI 构建命令："
	@echo ""
	@echo "  build        构建二进制文件"
	@echo "  build-all    构建多平台版本"
	@echo "  install      安装到系统路径 (/usr/local/bin)"
	@echo "  install-local 安装到用户目录 (~/.local/bin)"
	@echo "  test         运行测试"
	@echo "  run          运行开发版本"
	@echo "  clean        清理构建文件"
	@echo "  fmt          格式化代码"
	@echo "  lint         代码检查"
	@echo "  tidy         整理依赖"
	@echo "  release      创建发布版本"
	@echo "  dev          开发流程 (fmt + lint + test + run)"
	@echo "  help         显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build                    # 构建"
	@echo "  make run ARGS='create test'   # 运行并创建测试项目"
	@echo "  make install-local            # 本地安装"