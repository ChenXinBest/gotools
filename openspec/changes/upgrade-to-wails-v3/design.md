## Context

当前项目基于 Wails v2.12.0，使用 `wails.Run()` 单窗口模式，所有业务方法绑定在 `App` 结构体上。前后端通过 `frontend/wailsjs/` 下的自动生成绑定通信。

Wails v3 引入了服务式架构（Service Registration）、多窗口管理、系统托盘等能力，同时 CLI 从 `wails` 改为 `wails3`，构建系统新增 `Taskfile.yml`。v3 当前为 alpha 阶段，API 仍在演进。

## Goals / Non-Goals

**Goals:**
- 将整个项目的后端从 Wails v2 API 迁移到 v3 API
- 将前端绑定从 `wailsjs` 迁移到 v3 的 `bindings` 体系
- 更新构建系统以适配 `wails3` CLI 和 `Taskfile.yml`
- 在迁移基础上利用 v3 新特性：系统托盘、原生菜单
- 保持现有功能不变（进程管理、数据库备份）
- 确保 Windows 平台正常运行

**Non-Goals:**
- 不改变 Vue 3 前端框架版本或架构
- 不重构业务逻辑（process service, database service 的 Go 代码不变）
- 不引入新的数据库功能
- 不做跨平台（macOS/Linux）的深度调优

## Decisions

### 1. 分层迁移策略：先做核心迁移，再追加新特性

将升级分两层：第一层是纯 API 迁移（功能等价），第二层是利用 v3 新特性做增强。

**理由：** 分离风险。如果 v3 alpha 有 bug，可以快速定位是迁移问题还是新特性问题。

### 2. 后端架构拆分

```
v2:                          v3:
App (one struct)             ├── ProcessService  (独立)
  ├── processService         ├── DatabaseService (独立)
  ├── databaseService        ├── DialogService   (独立)
  └── dialogService          └── 用 application.NewService() 注册
```

将 `App` 结构体拆分为独立的 Service 结构体：
- `ProcessService` — 进程管理相关方法
- `DatabaseService` — 数据库备份/导入/导出相关方法  
- `DialogService` — 文件对话框相关方法

**理由：** v3 的设计哲学就是服务化注册，每个 Service 是独立可测试的 Go 结构体。拆分后的代码更清晰，也方便后续单独抽离为独立模块。

### 3. 对话框方案

v2 中 `DialogService` 使用 `runtime.OpenFileDialog()` 等 API。v3 中将使用 `application.OpenFileDialog()` 系列（位于 `wailsapp/wails/v3/pkg/application` 包）。

**需要确认：** v3 alpha 中对话框 API 的签名是否与 v2 兼容。如果 API 有变化，封装一层适配。

### 4. 资源服务（Assets）方案

```
v2:                          v3:
AssetServer:                 Assets: application.AssetOptions{
  &assetserver.Options{          Handler: application.BundledAssetFileServer(assets),
    Assets: assets,           }
  }
```

使用 `application.BundledAssetFileServer()` 或 `application.AssetFileServerFS()` 代替 v2 的 `assetserver.Options`。

### 5. 平台选项迁移

```
v2:                                              v3:
Windows: &windows.Options{...}                   Windows: application.WindowsOptions{...}
Mac: &mac.Options{...}                           Mac: application.MacOptions{...}
```

所有平台选项从子包移到 `application` 包内，选项字段名和类型可能变化，需对照 v3 API 逐一映射。

### 6. 前端绑定迁移策略

```
v2:                                        v3:
import { ... } from                        import { ... } from
  '../wailsjs/go/main/App'                  '../bindings/gotools/main/App.js'
window.runtime.WindowXXX()                  // 通过 @wailsapp/runtime 或 bindings 调用
```

v3 的绑定生成目录从 `wailsjs/go/` 改为 `bindings/`，模块名基于 `wails.json` 中的 `name` 字段。

**迁移方式：** 
1. 先用 `wails3 dev` 生成新的绑定
2. 批量更新前端 import 路径
3. 替换 runtime 调用为 v3 等效 API

### 7. 构建系统

```
v2:        wails.json 管理配置
           wails dev / wails build

v3:        wails.json (格式变化) + Taskfile.yml
           wails3 dev / wails3 build
```

`Taskfile.yml` 是 v3 推荐的构建编排工具，类似 Makefile 但更简洁。v3 也支持直接用 `wails3` 命令。

### 8. 系统托盘实现方案

利用 v3 的 `application.NewSystemTray()` API 实现：
- 窗口关闭时最小化到托盘而不是退出
- 托盘菜单：显示/隐藏窗口、退出
- 托盘图标使用现有 `build/appicon.png`

### 9. 原生菜单实现方案

利用 v3 的 `application.NewMenu()` / `app.NewMenuItem()` API：
- 应用菜单（文件、编辑、视图、帮助）
- 替代部分现有 Vue 组件中的按钮操作
- 添加键盘快捷键（Ctrl+Q 退出等）

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| v3 alpha API 不稳定，后续版本可能 breaking | 锁定 `v3.0.0-alpha.72` 版本，不追最新 |
| 对话框 API 签名不兼容 | 先在测试项目中验证，封装适配层 |
| 前端绑定生成路径/格式与预期不符 | 先做小范围 POC：只绑定一个 Service，验证绑定可用 |
| `Taskfile.yml` 学习成本 | v3 项目模板自带 Taskfile.yml，参考模板即可 |
| Windows 平台特定功能（如 WebView2）行为差异 | 在 Windows 上开发测试，及时发现问题 |
| 没有官方 v2→v3 迁移文档 | 参考 v3 示例项目和 GitHub 上的真实项目（MrRSS, tingly-box 等）源码 |

## Open Questions

- v3 alpha 中 `go-webview2` 版本是否与当前 Windows 系统兼容？
- `DialogService` 的 `SelectFolder()` 等是否完全对应 v3 的 API？
- 前端 `window.runtime.XXX` 调用在 v3 中的精确替代方案是什么——是通过 `@wailsapp/runtime` npm 包还是绑定生成？
- 后端方法返回值中的自定义 struct（如 `ExportResponse`、`ImportResponse`）是否需要调整以适配 v3 的序列化？
