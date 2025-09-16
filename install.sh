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

# 检查 Go 环境
check_go() {
    if ! command -v go &> /dev/null; then
        error "未找到 Go 环境，请先安装 Go (https://golang.org/dl/)"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    info "发现 Go 版本: $GO_VERSION"
}

# 从源码编译安装
install_from_source() {
    info "正在从源码编译安装 GoWeb CLI..."
    
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    info "克隆仓库..."
    git clone "$REPO_URL" goweb-cli
    cd goweb-cli
    
    info "编译二进制文件..."
    go build -o "$BINARY_NAME" .
    
    info "安装二进制文件..."
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # 清理临时目录
    cd /
    rm -rf "$TEMP_DIR"
    
    success "GoWeb CLI 安装完成！"
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
    check_go
    create_install_dir
    install_from_source
    check_path
    verify_installation
    
    echo ""
    success "安装完成！开始使用 GoWeb CLI 创建您的第一个项目吧！"
}

# 运行主函数
main "$@"