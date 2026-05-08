@echo off
chcp 65001 >nul

:: GoTools 构建脚本 - 用于交叉编译和版本注入

setlocal enabledelayedexpansion

:: 版本信息
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set GIT_TAG=%%i
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
for /f "tokens=*" %%i in ('git rev-parse --abbrev-ref HEAD 2^>nul') do set GIT_BRANCH=%%i
for /f "tokens=*" %%i in ('date /t') do set BUILD_DATE=%%i
for /f "tokens=*" %%i in ('time /t') do set BUILD_TIME=%%i

:: 如果获取不到Git信息，使用默认值
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown
if "%GIT_BRANCH%"=="" set GIT_BRANCH=unknown
if "%GIT_TAG%"=="" set GIT_TAG=dev

set VERSION=%GIT_TAG%
set BUILD_TIMESTAMP=%BUILD_DATE% %BUILD_TIME%

echo ========================================
echo GoTools 构建脚本
echo ========================================
echo 版本: %VERSION%
echo 提交: %GIT_COMMIT%
echo 分支: %GIT_BRANCH%
echo 时间: %BUILD_TIMESTAMP%
echo ========================================

:: 创建输出目录
if not exist build mkdir build

:: 构建标志
set LDFLAGS=-X "gotools/internal/version.Version=%VERSION%"
set LDFLAGS=%LDFLAGS% -X "gotools/internal/version.BuildTime=%BUILD_TIMESTAMP%"
set LDFLAGS=%LDFLAGS% -X "gotools/internal/version.GitCommit=%GIT_COMMIT%"
set LDFLAGS=%LDFLAGS% -X "gotools/internal/version.GitBranch=%GIT_BRANCH%"

:: 构建Windows/amd64版本
echo.
echo [1/4] 构建 windows/amd64 版本...
set GOOS=windows
set GOARCH=amd64
wails3 build -ldflags "%LDFLAGS%"
if errorlevel 1 (
    echo 构建失败！
    goto error
)
echo 构建完成: build\gotools-windows-amd64.exe

:: 构建Linux/amd64版本
echo.
echo [2/4] 构建 linux/amd64 版本...
set GOOS=linux
set GOARCH=amd64
wails build -o "build\gotools-linux-amd64" -ldflags "%LDFLAGS%"
if errorlevel 1 (
    echo 构建失败！
    goto error
)
echo 构建完成: build\gotools-linux-amd64

:: 构建macOS/amd64版本
echo.
echo [3/4] 构建 darwin/amd64 版本...
set GOOS=darwin
set GOARCH=amd64
wails build -o "build\gotools-darwin-amd64" -ldflags "%LDFLAGS%"
if errorlevel 1 (
    echo 构建失败！
    goto error
)
echo 构建完成: build\gotools-darwin-amd64

:: 构建macOS/arm64版本
echo.
echo [4/4] 构建 darwin/arm64 版本...
set GOOS=darwin
set GOARCH=arm64
wails build -o "build\gotools-darwin-arm64" -ldflags "%LDFLAGS%"
if errorlevel 1 (
    echo 构建失败！
    goto error
)
echo 构建完成: build\gotools-darwin-arm64

:: 重置环境变量
set GOOS=
set GOARCH=

:: 生成构建信息文件
echo.
echo 生成构建信息...
echo GoTools 构建信息 > build\build-info.txt
echo =============== >> build\build-info.txt
echo 版本: %VERSION% >> build\build-info.txt
echo 提交: %GIT_COMMIT% >> build\build-info.txt
echo 分支: %GIT_BRANCH% >> build\build-info.txt
echo 时间: %BUILD_TIMESTAMP% >> build\build-info.txt
echo 平台: windows/amd64, linux/amd64, darwin/amd64, darwin/arm64 >> build\build-info.txt
echo Go版本: %GO_VERSION% >> build\build-info.txt

echo.
echo ========================================
echo 所有构建任务完成！
echo 二进制文件已保存到 build/ 目录
echo 构建信息已保存到 build/build-info.txt
echo ========================================
pause
exit /b 0

:error
echo.
echo ========================================
echo 构建过程中出现错误！
echo ========================================
pause
exit /b 1
