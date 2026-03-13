# GoTools

## 项目简介

GoTools是一个基于Go语言和Wails框架开发的工具集合，旨在汇总平时开发中用到的各种小工具，提高开发效率。

### 主要功能

- **进程管理器**：实时监控系统进程，查看CPU、内存使用情况，支持进程搜索和终止
- **更多工具**：持续更新中...

## 技术栈

- **后端**：Go语言
- **前端**：Vue 3 + Vite
- **框架**：Wails

## 快速开始

### 环境要求

- Go 1.20+
- Node.js 16+
- Wails CLI

### 安装

1. 克隆项目

```bash
git clone https://github.com/yourusername/gotools.git
cd gotools
```

2. 安装依赖

```bash
# 安装前端依赖
npm install

# 安装Go依赖
go mod tidy
```

### 开发模式

运行以下命令启动开发服务器：

```bash
wails dev
```

这将启动一个Vite开发服务器，提供前端热重载功能。同时，还会启动一个开发服务器，运行在 http://localhost:34115，您可以在浏览器中访问并通过devtools调用Go代码。

### 构建

要构建可分发的生产模式包，使用：

```bash
wails build
```

## 项目结构

```
gotools/
├── app.go              # 应用入口
├── frontend/           # 前端代码
│   ├── src/            # 源代码
│   │   ├── components/ # 组件
│   │   └── wailsjs/    # Wails生成的绑定
│   └── package.json    # 前端依赖
├── internal/           # 内部包
│   └── tools/          # 工具实现
├── main.go             # 主入口
├── go.mod              # Go模块文件
└── wails.json          # Wails配置
```

## 功能说明

### 进程管理器

- 实时显示系统进程列表
- 按CPU、内存、进程名等排序
- 搜索进程（支持进程名、PID、命令行等）
- 查看进程详细信息
- 终止选中进程

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

本项目采用MIT许可证。详见LICENSE文件。