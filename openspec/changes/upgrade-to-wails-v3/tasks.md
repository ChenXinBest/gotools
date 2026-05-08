## 1. 环境准备

- [x] 1.1 安装 Wails v3 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
- [x] 1.2 用 `wails3 init -n test -t vue` 创建测试项目，理解 v3 项目结构
- [x] 1.3 确认 `wails3 dev` / `wails3 build` 在 Windows 上正常工作
- [x] 1.4 创建 git 分支 `wails-v3-upgrade`

## 2. Go 依赖升级

- [x] 2.1 更新 `go.mod`：将 `github.com/wailsapp/wails/v2 v2.12.0` 改为 `github.com/wailsapp/wails/v3`
- [x] 2.2 更新源代码移除 v2 引用后运行 `go mod tidy` 解决依赖冲突
- [x] 2.3 确认编译通过（`go build ./...`）

## 3. main.go — 应用入口重写

- [x] 3.1 导入路径更新：从 `github.com/wailsapp/wails/v2` 改为 `github.com/wailsapp/wails/v3/pkg/application`
- [x] 3.2 移除 v2 子包导入：`options`, `assetserver`, `mac`, `windows`
- [x] 3.3 用 `application.New(application.Options{...})` 替代 `wails.Run(&options.App{...})`
- [x] 3.4 用 `app.Window.NewWithOptions(...)` 显式创建主窗口
- [x] 3.5 配置 `Assets: application.AssetOptions{Handler: application.BundledAssetFileServer(assets)}`
- [x] 3.6 平台选项迁移：`application.WindowsOptions{...}` / `application.MacOptions{...}`
- [x] 3.7 生命周期回调迁移：各 Service 的 `Startup()` 方法自动处理
- [x] 3.8 确认 `go vet ./...` 通过

## 4. app.go — Service 拆分与注册

- [x] 4.1 将 `App` 结构体的对话框相关方法拆分为独立的 `DialogService` 结构体
- [x] 4.2 将 `App` 结构体的进程管理相关方法拆分为独立的 `ProcessService` 结构体
- [x] 4.3 将 `App` 结构体的数据库相关方法拆分为独立的 `DatabaseService` 结构体
- [x] 4.4 各 Service 保留 `ctx context.Context` 和 v3 兼容的 startup 模式
- [x] 4.5 在 `main.go` 中用 `application.NewService()` 注册所有 Service
- [x] 4.6 移除旧的 `App` 聚合结构体和 `NewApp()` 工厂方法
- [ ] 4.7 验证 `wails3 dev` 启动后前端能正常调用后端方法

## 5. 构建配置适配

- [x] 5.1 移除 `wails.json`（v3 不使用），创建 `build/config.yml`
- [x] 5.2 参考 v3 模板创建 `Taskfile.yml` + `build/Taskfile.yml`
- [ ] 5.3 验证 `wails3 dev` 正常运行

## 6. 前端绑定路径迁移

- [x] 6.1 运行 `wails3 generate bindings` 生成 v3 绑定文件到 `frontend/bindings/`
- [x] 6.2 更新 `frontend/src/stores/app.js`（无后端调用，无需修改）
- [x] 6.3 更新 `frontend/src/stores/process.js` → `bindings/gotools/processservice.js`
- [x] 6.4 更新 `frontend/src/stores/database.js` → `bindings/gotools/databaseservice.js`
- [x] 6.5 更新 `DatabaseBackup.vue` → `bindings/gotools/databaseservice.js` + `@wailsio/runtime`
- [x] 6.6 更新 `HelloWorld.vue` → `bindings/gotools/processservice.js`
- [x] 6.7 更新 `ProcessManager.vue` / `Export.vue` / `Settings.vue` 导入路径
- [x] 6.8 移除 `frontend/wailsjs/` 目录
- [x] 6.9 验证前端 `npm run build` 通过（87 modules, 无错误）

## 7. 对话框服务迁移

- [x] 7.1 调研确定：现有对话框基于 PowerShell（`tools/dialog_tool.go`），独立于 Wails 版本，无需迁移
- [ ] 7.2 可选：切换到 v3 原生对话框 API（`application.OpenFileDialog()` 等）
- [ ] 7.3 验证 `SelectFolder()`、`SelectFile()`、`SelectSaveFile()` 正常工作

## 8. 构建脚本更新

- [x] 8.1 更新 `build.bat`：将 `wails` 命令替换为 `wails3`
- [x] 8.2 更新 `build.sh`：将 `wails` 命令替换为 `wails3`
- [x] 8.3 验证 `wails3 build` 生成可执行文件

## 9. 系统托盘功能

- [x] 9.1 创建 `tray_menu.go`，集成系统托盘 + 应用菜单
- [x] 9.2 配置托盘图标（`build/appicon.png`）
- [x] 9.3 实现托盘菜单：「显示窗口」「退出」
- [x] 9.4 关闭窗口时最小化到托盘（`DisableQuitOnLastWindowClosed: true`）
- [ ] 9.5 运行时验证托盘功能在 Windows 上正常

## 10. 原生菜单功能

- [x] 10.1 创建应用菜单（`tray_menu.go`）
- [x] 10.2 实现「文件」菜单（退出）
- [x] 10.3 实现「视图」菜单（页面切换选项 + 刷新）
- [x] 10.4 实现「帮助」菜单（关于 GoTools）
- [x] 10.5 实现右键上下文菜单框架（`createContextMenu`）
- [ ] 10.6 运行时验证菜单在 Windows 上正常

## 11. 多窗口功能

- [x] 11.1 创建 `multi_window.go` — 导出进度窗口 + 冲突检测窗口
- [x] 11.2 事件系统通信（`app.Event.On/Emit`）
- [x] 11.3 Vue 组件 `ExportProgress.vue` + `ImportConflicts.vue` + 路由注册
- [ ] 11.4 运行时验证多窗口在 Windows 上正常

## 12. 验证与清理

- [x] 12.1 `go build ./...` 通过
- [x] 12.2 `wails3 build` 完整构建通过（frontend 91 modules + go binary）
- [x] 12.3 `go vet ./...` 无错误
- [x] 12.4 `go test ./...` 54 passed in 8 packages
- [x] 12.5 前端 `npm run build` 通过
- [x] 12.6 更新 `AGENTS.md` 和 `README.md` 中 Wails 版本信息
- [ ] 12.7 运行时功能验证（需要 GUI 环境）
