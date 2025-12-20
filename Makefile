# ============================================
# Wails 应用 Makefile
# 项目：Media App
# ============================================

.PHONY: help dev build build-darwin build-windows build-linux clean install

# 默认目标
.DEFAULT_GOAL := help

# 变量
VERSION ?= 1.0.0
ARCH ?= $(shell go env GOARCH)

# 帮助信息
help:
	@echo "Media App 构建命令"
	@echo ""
	@echo "使用方法:"
	@echo "  make dev              启动开发模式"
	@echo "  make build            构建当前平台"
	@echo "  make build-darwin     构建 macOS 版本"
	@echo "  make build-windows    构建 Windows 版本"
	@echo "  make build-linux      构建 Linux 版本"
	@echo "  make build-all        构建所有平台"
	@echo "  make clean            清理构建目录"
	@echo "  make install          安装依赖"
	@echo ""
	@echo "变量:"
	@echo "  VERSION=x.x.x         设置版本号 (默认: 1.0.0)"
	@echo "  ARCH=amd64|arm64      设置目标架构 (默认: 当前架构)"
	@echo ""
	@echo "示例:"
	@echo "  make build-darwin ARCH=arm64"
	@echo "  make build-windows VERSION=2.0.0"

# 安装依赖
install:
	@echo "安装前端依赖..."
	@cd frontend && npm install

# 开发模式
dev:
	@wails dev

# 构建当前平台
build: install
	@wails build -clean

# 构建 macOS
build-darwin: install
	@echo "构建 macOS $(ARCH)..."
	@wails build -clean -platform darwin/$(ARCH)

# 构建 macOS 通用二进制
build-darwin-universal: install
	@echo "构建 macOS amd64..."
	@wails build -clean -platform darwin/amd64
	@echo "构建 macOS arm64..."
	@wails build -clean -platform darwin/arm64

# 构建 Windows
build-windows: install
	@echo "构建 Windows $(ARCH)..."
	@wails build -clean -platform windows/$(ARCH)

# 构建 Windows 带 NSIS 安装包
build-windows-nsis: install
	@echo "构建 Windows $(ARCH) (NSIS)..."
	@wails build -clean -platform windows/$(ARCH) -nsis

# 构建 Linux
build-linux: install
	@echo "构建 Linux $(ARCH)..."
	@wails build -clean -platform linux/$(ARCH)

# 构建所有平台
build-all: install
	@echo "构建所有平台..."
	@wails build -clean -platform darwin/amd64
	@wails build -clean -platform darwin/arm64
	@wails build -clean -platform windows/amd64
	@wails build -clean -platform linux/amd64

# 清理
clean:
	@echo "清理构建目录..."
	@rm -rf build/bin
	@rm -rf frontend/dist
	@echo "清理完成"


