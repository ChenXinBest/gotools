## Context

GoTools 是一个基于 Go + Wails + Vue 3 开发的桌面工具集，目前实现了两大核心功能：进程管理器和数据库备份导出。项目采用前后分离架构，通过 Wails 框架实现 Go 后端与 Vue 前端的双向绑定。

**技术栈：**
- 后端：Go 1.26+，Wails v2.12.0
- 前端：Vue 3 + Vite + Pinia
- 系统依赖：gopsutil/v3（进程信息）、MySQL Shell、mysqldump

**项目结构：**
- `internal/tools/` - 工具函数封装（MySQL操作、进程操作、验证逻辑）
- `internal/services/` - 业务服务层
- `internal/config/` - 配置管理
- `internal/log/` - 日志系统
- `frontend/src/` - Vue 前端代码

## Goals / Non-Goals

**Goals:**
- 建立完整的规格说明文档体系
- 明确当前已实现功能的详细规格
- 梳理系统架构和技术选型
- 识别待改进项和优化方向

**Non-Goals:**
- 本次变更不修改任何代码
- 不引入新功能
- 不进行任何重构

## Decisions

### 1. 前后端分离架构（通过 Wails 绑定）

**决策**：使用 Wails 框架实现前后端绑定，前端通过自动生成的 JS 绑定调用 Go 函数。

**理由**：
- Wails 提供原生桌面应用体验
- 自动生成类型安全的绑定
- 支持前端热重载开发模式

### 2. Pinia 状态管理

**决策**：前端使用 Pinia 进行状态管理，按功能模块划分 store。

**理由**：
- Pinia 是 Vue 3 官方推荐的状态管理库
- 简洁的 API 和完善的 TypeScript 支持
- 模块化设计便于维护

### 3. 进程管理实现

**决策**：使用 gopsutil/v3 库获取系统进程信息。

**理由**：
- 跨平台支持（Windows/Linux/macOS）
- 成熟的 Go 生态系统
- 提供了丰富的系统信息 API

### 4. 数据库备份工具选择

**决策**：支持 MySQL Shell 和 mysqldump 两种导出工具。

**理由**：
- MySQL Shell 功能更强大，支持多线程并行导出
- mysqldump 兼容性好，广泛使用
- 用户可根据实际环境选择

### 5. 配置持久化

**决策**：使用 JSON 文件存储配置，路径为可执行文件同目录。

**理由**：
- 简单直接，无需额外依赖
- 便于用户查看和修改
- Windows 环境下使用 filepath 确保跨平台兼容

## Risks / Trade-offs

[Risk] Wails 绑定生成依赖开发服务
→ [Mitigation] 修改后端代码后需重启 `wails dev` 以更新绑定

[Risk] 配置文件路径包含敏感信息（数据库密码）
→ [Mitigation] 配置文件中存储密码，但日志中不记录敏感信息

[Risk] MySQL Shell/mysqldump 命令需预装
→ [Mitigation] 导出功能依赖系统上安装的 MySQL 工具，无工具时给出明确错误提示

[Risk] 进程信息实时性
→ [Mitigation] CPU/内存使用率为采样值，存在一定延迟；提供手动刷新和自动刷新功能

## Open Questions

1. 是否需要支持其他数据库（PostgreSQL、SQLite）？
2. 进程管理器的排序和搜索性能在进程数多时可能下降，是否需要优化？
3. 配置文件是否需要加密存储？
4. 是否需要实现导入/导出的进度通知机制？
