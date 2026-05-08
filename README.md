# GoTools

## 项目简介

GoTools是一个基于Go语言和Wails框架开发的工具集合，旨在汇总平时开发中用到的各种小工具，提高开发效率。

### 主要功能

- **进程管理器**：实时监控系统进程，查看CPU、内存使用情况，支持进程搜索、批量操作和进程终止
- **数据库备份**：
  - 支持MySQL数据库连接管理
  - 支持MySQL Shell和mysqldump两种导出工具
  - 支持数据库/表导出和导入功能
  - 支持导入前冲突检测和自动处理
- **更多工具**：持续更新中...

## 技术栈

- **后端**：Go 1.26+
- **前端**：Vue 3 + Vite 5
- **框架**：Wails v3 (v3.0.0-alpha.87)
- **依赖库**：
  - `github.com/shirou/gopsutil/v3` - 系统进程信息获取
  - `vue-icons-plus/fa` - Font Awesome 图标库
  - `@wailsio/runtime` - Wails v3 前端运行时

## 快速开始

### 环境要求

- **Go 1.26+**（本项目使用Go 1.26最新语法特性）
- Node.js 18+
- **Wails v3 CLI**（需要从源码编译安装，见下方说明）

### Windows环境配置说明

本项目已针对Windows环境进行路径处理优化：

1. **路径分隔符**：所有文件路径操作使用 `filepath` 包确保跨平台兼容
2. **MySQL Shell路径**：导出功能自动将Windows路径转换为正斜杠格式（`filepath.ToSlash`）
3. **特殊字符处理**：路径中包含空格或特殊字符时，程序会正确处理

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
wails3 dev
```

这将启动一个Vite开发服务器，提供前端热重载功能。

### 构建

要构建可分发的生产模式包，使用：

```bash
wails3 build
```

## 项目结构

```
gotools/
├── app.go                    # Wails应用入口
├── main.go                   # Wails v3 应用入口
├── models.go                 # DTO 数据结构
├── process_service.go        # 进程管理 Wails 服务
├── database_service.go       # 数据库 Wails 服务
├── dialog_service.go         # 对话框 Wails 服务
├── tray_menu.go              # 系统托盘 + 应用菜单
├── multi_window.go           # 多窗口管理
├── Taskfile.yml              # v3 构建任务
├── go.mod                    # Go模块文件
├── build/                    # 构建配置
│   ├── config.yml            # v3 构建配置
│   ├── Taskfile.yml          # 构建子任务
│   ├── windows/              # Windows 构建脚本
│   ├── darwin/               # macOS 构建脚本
│   └── linux/                # Linux 构建脚本
├── frontend/                 # 前端代码
│   ├── src/                  # 源代码
│   │   ├── components/       # Vue组件
│   │   ├── views/            # 页面视图
│   │   ├── stores/           # Pinia状态管理
│   │   ├── router/           # Vue Router
│   │   ├── bindings/         # Wails v3 自动生成绑定
│   │   └── style.css         # 全局样式
│   ├── node_modules/         # 前端依赖
│   └── package.json          # 前端依赖
├── internal/                 # 内部包
│   ├── config/               # 配置管理
│   ├── log/                  # 日志系统
│   ├── services/             # 业务服务层
│   ├── tools/                # 工具实现
│   └── version/              # 版本信息
├── test/                     # 测试文件
│   └── tools/                # 工具测试
└── openspec/                 # 变更提案和规格文档
```

## 功能说明

### 进程管理器

#### 功能特性

- 实时显示系统所有进程信息（PID、进程名、命令行、CPU、内存）
- 按应用名分组聚合显示，同一应用的多个进程归类在一起
- CPU/内存使用率可视化显示（进度条）
- 显示进程网络端口监听信息
- 支持按CPU、内存、进程名等排序

#### 操作说明

- **搜索进程**：支持按进程名、PID、命令行、端口等关键词搜索
- **批量选择**：
  - `Ctrl + 点击`：多选/取消选择
  - `Shift + 点击`：范围选择
  - 拖拽选择：按住左键拖动选择多个进程
- **终止进程**：右键菜单选择"终止进程"

### 数据库备份

#### 功能概述

数据库备份工具提供MySQL数据库的导入导出功能，支持两种导出工具：MySQL Shell 和 mysqldump。

#### 连接管理

- 添加、编辑、删除数据库连接配置
- 连接配置自动持久化存储
- 支持管理多个数据库连接

#### 导出功能

**MySQL Shell 方式**（推荐，功能更强大）：

- 支持导出整个数据库实例、指定数据库或指定表
- 多线程并行导出，性能更优
- 支持压缩格式：gzip、zstd、无压缩
- 支持设置分块大小（Chunk Size）
- 支持选择性包含/排除特定数据库或表
- 支持跳过Definer、跳过Binlog选项

**mysqldump 方式**（兼容性好）：

- 支持导出数据库和表
- 支持 single-transaction（一致性快照）
- 支持导出存储过程和事件
- 支持压缩输出

#### 导入功能

- 导入已导出的数据库文件
- 多线程并行导入（MySQL Shell方式）
- **导入冲突检测**：
  - 导入前自动检测目标数据库是否存在同名表
  - 显示冲突对象列表
  - 支持自动删除冲突表后继续导入

## Go 1.26 新特性使用

本项目已采用Go 1.26的最新语法特性：

- 使用 `map[string]any` 替代 `map[string]interface{}`
- 使用 `strings.SplitSeq` 进行高效的字符串分割迭代
- 使用 `log/slog` 结构化日志
- 使用 `filepath.ToSlash` 进行跨平台路径转换

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

本项目采用MIT许可证。详见LICENSE文件。
