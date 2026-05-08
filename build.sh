#!/bin/bash

# GoTools 构建脚本 - 用于交叉编译和版本注入

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 版本信息
GIT_TAG=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date '+%Y-%m-%d')
BUILD_TIME=$(date '+%H:%M:%S')
BUILD_TIMESTAMP="$BUILD_DATE $BUILD_TIME"
GO_VERSION=$(go version | awk '{print $3}')

VERSION=$GIT_TAG

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}GoTools 构建脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "版本: ${GREEN}$VERSION${NC}"
echo -e "提交: ${GREEN}$GIT_COMMIT${NC}"
echo -e "分支: ${GREEN}$GIT_BRANCH${NC}"
echo -e "时间: ${GREEN}$BUILD_TIMESTAMP${NC}"
echo -e "Go版本: ${GREEN}$GO_VERSION${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 创建输出目录
mkdir -p build

# 构建标志
LDFLAGS="-X gotools/internal/version.Version=$VERSION"
LDFLAGS="$LDFLAGS -X gotools/internal/version.BuildTime=$BUILD_TIMESTAMP"
LDFLAGS="$LDFLAGS -X gotools/internal/version.GitCommit=$GIT_COMMIT"
LDFLAGS="$LDFLAGS -X gotools/internal/version.GitBranch=$GIT_BRANCH"

# 定义构建目标
targets=("windows/amd64" "linux/amd64" "darwin/amd64" "darwin/arm64")
total=${#targets[@]}
count=0

for target in "${targets[@]}"; do
    IFS="/" read -r os arch <<< "$target"
    count=$((count + 1))
    
    echo -e "${YELLOW}[$count/$total]${NC} 构建 ${os}/${arch} 版本..."
    
    # 设置环境变量
    export GOOS="$os"
    export GOARCH="$arch"
    
    # 构建
    output="build/gotools-${os}-${arch}"
    if [ "$os" == "windows" ]; then
        output="${output}.exe"
    fi
    
    # 使用wails构建
    if command -v wails3 &> /dev/null; then
        wails3 build -ldflags "$LDFLAGS"
    else
        # 如果没有wails，使用go build
        echo -e "${YELLOW}警告: wails 未找到，使用 go build${NC}"
        go build -ldflags "$LDFLAGS" -o "$output" .
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}构建完成: $output${NC}"
    else
        echo -e "${RED}构建失败: $os/$arch${NC}"
        exit 1
    fi
    
    echo ""
done

# 重置环境变量
unset GOOS
unset GOARCH

# 生成构建信息文件
echo "GoTools 构建信息" > build/build-info.txt
echo "=================" >> build/build-info.txt
echo "版本: $VERSION" >> build/build-info.txt
echo "提交: $GIT_COMMIT" >> build/build-info.txt
echo "分支: $GIT_BRANCH" >> build/build-info.txt
echo "时间: $BUILD_TIMESTAMP" >> build/build-info.txt
echo "Go版本: $GO_VERSION" >> build/build-info.txt
echo "平台: ${targets[*]}" >> build/build-info.txt

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}所有构建任务完成！${NC}"
echo -e "二进制文件已保存到 ${GREEN}build/${NC} 目录"
echo -e "构建信息已保存到 ${GREEN}build/build-info.txt${NC}"
echo -e "${BLUE}========================================${NC}"
