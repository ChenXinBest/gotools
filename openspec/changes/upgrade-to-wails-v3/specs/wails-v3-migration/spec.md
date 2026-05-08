## ADDED Requirements

### Requirement: Wails v2 到 v3 的 API 迁移

应用入口 SHALL 使用 `application.New()` 替代 `wails.Run()`。

应用窗口 SHALL 通过 `app.NewWindow()` 显式创建。

后端服务 SHALL 通过 `application.NewService()` 注册，替代 v2 的 `Bind` 字段。

资源服务 SHALL 使用 `application.AssetOptions` 替代 v2 的 `assetserver.Options`。

平台选项 SHALL 从子包（`windows.Options`/`mac.Options`）迁移到 `application` 包内。

Go 模块依赖 SHALL 从 `github.com/wailsapp/wails/v2` 改为 `github.com/wailsapp/wails/v3`。

#### Scenario: 应用入口使用 application.New()

- **WHEN** 启动应用
- **THEN** 使用 `application.New(application.Options{...})` 创建应用实例
- **AND** 调用 `app.NewWindow()` 创建主窗口
- **AND** 调用 `app.Run()` 启动应用

#### Scenario: Service 注册替代 Bind

- **WHEN** 注册后端服务
- **THEN** 使用 `application.NewService(&ProcessService{})` 注册进程管理服务
- **AND** 使用 `application.NewService(&DatabaseService{})` 注册数据库服务
- **AND** 使用 `application.NewService(&DialogService{})` 注册对话框服务
- **AND** 不使用 v2 的 `Bind` 字段

#### Scenario: 资源服务迁移

- **WHEN** 配置前端资源
- **THEN** 使用 `application.AssetOptions{Handler: application.BundledAssetFileServer(assets)}` 
- **AND** 不使用 v2 的 `assetserver.Options`

#### Scenario: 平台选项迁移

- **WHEN** 配置平台特定选项
- **THEN** 使用 `application.WindowsOptions{...}` 替代 `windows.Options{...}`
- **AND** 使用 `application.MacOptions{...}` 替代 `mac.Options{...}`

### Requirement: 前端绑定路径迁移

前端代码 SHALL 从 v3 生成的 `bindings/` 目录导入 Go 方法，替代 v2 的 `wailsjs/` 目录。

前端 Runtime 调用 SHALL 使用 v3 对应的 API。

#### Scenario: 导入路径更新

- **WHEN** 前端调用后端方法
- **THEN** import 路径从 `../wailsjs/go/main/App` 改为 `../bindings/gotools/main/App.js`

#### Scenario: Runtime 调用更新

- **WHEN** 前端调用窗口操作等 runtime 方法
- **THEN** 使用 v3 生成的 runtime 绑定替代 `window.runtime.XXX`

### Requirement: 构建系统迁移

构建命令 SHALL 从 `wails` 改为 `wails3`。

项目 SHALL 包含 `Taskfile.yml` 以支持 v3 的构建流程。

`wails.json` 配置 SHALL 适配 v3 格式。

#### Scenario: 开发者启动开发服务器

- **WHEN** 运行 `wails3 dev`
- **THEN** 启动带热重载的开发服务器

#### Scenario: 生产构建

- **WHEN** 运行 `wails3 build`
- **THEN** 生成生产可执行文件

#### Scenario: 构建脚本更新

- **WHEN** 执行 `build.bat` 或 `build.sh`
- **THEN** 内部使用 `wails3` 命令而非 `wails`
