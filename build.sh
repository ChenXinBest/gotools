#!/bin/bash

# 构建脚本 - 用于交叉编译GoTools

set -e

echo "开始构建GoTools..."

# 创建输出目录
mkdir -p build

# 定义构建目标
targets=("windows/amd64" "linux/amd64" "darwin/amd64" "darwin/arm64")

for target in "${targets[@]}"; do
    IFS="/" read -r os arch <<< "$target"
    echo "构建 $os/$arch 版本..."
    
    # 设置环境变量
    export GOOS="$os"
    export GOARCH="$arch"
    
    # 构建
    output="build/gotools-${os}-${arch}"
    if [ "$os" == "windows" ]; then
        output="${output}.exe"
    fi
    
    wails build -o "$output"
    
    echo "构建完成: $output"
done

echo "所有构建任务完成！"
echo "二进制文件已保存到 build/ 目录"
