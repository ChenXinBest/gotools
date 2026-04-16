## Why

本项目（GoTools）已实现进程管理器和数据库备份两大核心功能，但缺乏正式的规格说明文档。通过创建项目状态快照，建立规格说明体系，为后续功能迭代和代码维护提供基准参考。

## What Changes

本次变更不涉及功能修改，而是建立项目规格文档体系：

- 创建项目规格说明文档（specs），记录当前已实现功能的详细规格
- 创建架构设计文档（design），描述系统整体架构和技术选型
- 创建任务清单（tasks），整理待改进项和优化方向

## Capabilities

### New Capabilities

- `process-manager`: 进程管理器功能，包含进程列表展示、搜索、排序、批量选择、进程终止
- `database-backup`: 数据库备份功能，包含连接管理、数据导出（MySQL Shell/mysqldump）、数据导入、冲突检测
- `database-connections`: 数据库连接配置管理（增删改查）
- `database-export`: 数据导出功能，支持多种导出工具和压缩格式
- `database-import`: 数据导入功能，支持导入前冲突检测

### Modified Capabilities

（无 - 本次为状态快照，暂无需求变更）

## Impact

- 文档体系：`openspec/specs/` 目录下创建规格文档
- 开发规范：后续功能开发需遵循规格驱动开发流程
- 技术栈确认：Go 1.26 + Wails v2 + Vue 3 + Pinia
