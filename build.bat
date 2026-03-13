@echo off
chcp 65001 >nul

:: 构建脚本 - 用于交叉编译GoTools

echo 开始构建GoTools...

:: 创建输出目录
if not exist build mkdir build

:: 构建Windows/amd64版本
echo 构建 windows/amd64 版本...
wails build -o build\gotools-windows-amd64.exe

:: 构建Linux/amd64版本
echo 构建 linux/amd64 版本...
set GOOS=linux
set GOARCH=amd64
wails build -o build\gotools-linux-amd64

:: 构建macOS/amd64版本
echo 构建 darwin/amd64 版本...
set GOOS=darwin
set GOARCH=amd64
wails build -o build\gotools-darwin-amd64

:: 构建macOS/arm64版本
echo 构建 darwin/arm64 版本...
set GOOS=darwin
set GOARCH=arm64
wails build -o build\gotools-darwin-arm64

:: 重置环境变量
set GOOS=
set GOARCH=

echo 所有构建任务完成！
echo 二进制文件已保存到 build/ 目录
pause
