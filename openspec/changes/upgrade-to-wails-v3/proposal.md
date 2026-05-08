## Why

将桌面框架从 Wails v2（v2.12.0）升级到 Wails v3（alpha），以学习 v3 全新的服务式架构、多窗口管理和原生系统集成能力。本项目为个人学习项目，v3 的 alpha 状态提供了近距离接触前沿 API 设计的机会，迁移过程本身就是学习目标。

## What Changes

- **BREAKING**: 后端模块路径从 `github.com/wailsapp/wails/v2` 切换到 `github.com/wailsapp/wails/v3`
- **BREAKING**: 应用入口从 `wails.Run(&options.App{})` 改为 `application.New(application.Options{})` + `app.NewWindow()`
- **BREAKING**: 绑定机制从 `Bind: []interface{}{}` 改为 `Services: []application.Service{}` 注册模式
- **BREAKING**: 资源服务从 `AssetServer: &assetserver.Options{}` 改为 `Assets: application.AssetOptions{}`
- **BREAKING**: 前端绑定目录从 `frontend/wailsjs/` 迁移到 `frontend/bindings/`
- **BREAKING**: CLI 命令从 `wails` 改为 `wails3`，构建系统引入 `Taskfile.yml`
- **BREAKING**: 平台选项从子包（`windows.Options`/`mac.Options`）移到 `application` 包内
- **NEW**: 系统托盘支持（最小化到托盘）
- **NEW**: 原生菜单系统
- **NEW**: 多窗口能力（如导出进度独立窗口）

## Capabilities

### New Capabilities
- `wails-v3-migration`: Wails v2 到 v3 的整体迁移，包括后端 API、前端绑定、构建系统的全面升级
- `system-tray`: 系统托盘功能，支持最小化到托盘和后台运行
- `native-menu`: 原生应用菜单和右键上下文菜单
- `multi-window`: 多窗口管理，如导出/导入操作的独立进度窗口

### Modified Capabilities
<!-- 无现有 specs 需要修改 -->

## Impact

- `main.go`：完全重写应用创建逻辑（`wails.Run` → `application.New` + `NewWindow`）
- `app.go`：服务注册方式变更，需拆分独立 Service
- `go.mod`：依赖路径从 `v2` 改为 `v3`
- `wails.json`：配置格式可能变化
- `frontend/src/` 下所有引用 wailsjs 的文件需更新 import 路径
- `frontend/wailsjs/` 废弃，迁移到 `frontend/bindings/`
- `build/` 目录结构适配 v3
- `build.bat`/`build.sh`：命令从 `wails` 改为 `wails3`
- 新增 `Taskfile.yml` 构建配置
