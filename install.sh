#!/bin/bash

# GoWeb CLI 安装脚本
set -e

INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="goweb"
REPO_URL="https://github.com/Programmer-Bugs-Bunny/GreenFishCli"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 输出函数
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 检查操作系统
detect_os() {
    case "$(uname -s)" in
        Darwin)
            OS="darwin"
            ;;
        Linux)
            OS="linux"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            OS="windows"
            ;;
        *)
            error "不支持的操作系统: $(uname -s)"
            ;;
    esac
    
    info "检测到操作系统: $OS"
}

# 检查架构
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        *)
            error "不支持的架构: $(uname -m)"
            ;;
    esac
    
    info "检测到架构: $ARCH"
}

# 创建安装目录
create_install_dir() {
    if [ ! -d "$INSTALL_DIR" ]; then
        info "创建安装目录: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi
}

# 获取最新版本号
get_latest_version() {
    info "获取最新版本信息..."
    
    # 尝试从 GitHub API 获取最新版本
    if command -v curl &> /dev/null; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/Programmer-Bugs-Bunny/GreenFishCli/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/Programmer-Bugs-Bunny/GreenFishCli/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        error "需要 curl 或 wget 来下载文件"
    fi
    
    if [ -z "$LATEST_VERSION" ]; then
        return 1
    fi
    
    info "最新版本: $LATEST_VERSION"
}

# 从 GitHub Releases 下载二进制文件
download_binary() {
    info "正在下载预编译的二进制文件..."
    
    # 构建二进制文件名
    if [ "$OS" = "windows" ]; then
        BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}"
    fi
    
    DOWNLOAD_URL="https://github.com/Programmer-Bugs-Bunny/GreenFishCli/releases/download/${LATEST_VERSION}/${BINARY_FILE}"
    
    info "下载 URL: $DOWNLOAD_URL"
    
    # 下载二进制文件
    if command -v curl &> /dev/null; then
        if ! curl -L "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY_NAME" --fail --silent --show-error; then
            return 1
        fi
    elif command -v wget &> /dev/null; then
        if ! wget -q "$DOWNLOAD_URL" -O "$INSTALL_DIR/$BINARY_NAME"; then
            return 1
        fi
    else
        error "需要 curl 或 wget 来下载文件"
    fi
    
    # 设置执行权限
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    success "二进制文件下载完成！"
}

# 检查 Go 环境（仅用于源码安装）
check_go() {
    if ! command -v go &> /dev/null; then
        error "未找到 Go 环境，请先安装 Go (https://golang.org/dl/)"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    info "发现 Go 版本: $GO_VERSION"
}

# 从源码编译安装（备用方案）
install_from_source() {
    warning "正在尝试从源码编译安装（这需要 Go 环境）..."
    
    check_go
    
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    info "克隆仓库..."
    if ! git clone "$REPO_URL" goweb-cli; then
        error "克隆仓库失败"
    fi
    
    cd goweb-cli
    
    info "编译二进制文件..."
    if ! go build -o "$BINARY_NAME" .; then
        error "编译失败"
    fi
    
    info "安装二进制文件..."
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # 清理临时目录
    cd /
    rm -rf "$TEMP_DIR"
    
    success "GoWeb CLI 源码安装完成！"
}

# 主要安装函数
install_goweb() {
    # 首先尝试从 Releases 下载
    if get_latest_version && download_binary; then
        success "从 GitHub Releases 安装成功！"
        return 0
    fi
    
    # 如果下载失败，提供选择
    warning "从 Releases 下载失败"
    echo ""
    echo "可选的安装方式:"
    echo "1. 从源码编译安装（需要 Go 环境）"
    echo "2. 退出安装"
    echo ""
    read -p "请选择 (1-2): " choice
    
    case $choice in
        1)
            install_from_source
            ;;
        2)
            info "安装已取消"
            exit 0
            ;;
        *)
            error "无效选择"
            ;;
    esac
}

# 检查 PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warning "安装目录 $INSTALL_DIR 不在 PATH 中"
        echo ""
        echo "请将以下行添加到您的 shell 配置文件中："
        echo ""
        echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
        echo "支持的配置文件:"
        echo "  - Bash: ~/.bashrc 或 ~/.bash_profile"
        echo "  - Zsh: ~/.zshrc"
        echo "  - Fish: ~/.config/fish/config.fish"
        echo ""
        echo "然后执行: source ~/.bashrc (或对应的配置文件)"
    else
        success "PATH 配置正确"
    fi
}

# 验证安装
verify_installation() {
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        success "安装成功！"
        echo ""
        info "验证安装:"
        "$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null || true
        echo ""
        info "使用方法:"
        echo "    goweb create my-project    # 创建新项目"
        echo "    goweb --help              # 查看帮助"
    else
        error "安装失败，未找到二进制文件"
    fi
}

# 主函数
main() {
    echo "🚀 GoWeb CLI 安装程序"
    echo "===================="
    echo ""
    
    detect_os
    detect_arch
    create_install_dir
    install_goweb
    check_path
    verify_installation
    
    echo ""
    success "安装完成！开始使用 GoWeb CLI 创建您的第一个项目吧！"
}

# 运行主函数
main "$@"