## ADDED Requirements

### Requirement: 系统托盘支持

应用 SHALL 在关闭窗口时最小化到系统托盘而非退出（可选择行为）。

系统托盘 SHALL 显示应用图标。

系统托盘菜单 SHALL 包含「显示窗口」和「退出」选项。

点击系统托盘图标 SHALL 切换主窗口的显示/隐藏。

#### Scenario: 关闭窗口最小化到托盘

- **WHEN** 用户点击窗口关闭按钮
- **THEN** 窗口隐藏
- **AND** 应用保持运行
- **AND** 系统托盘图标可见

#### Scenario: 从托盘恢复窗口

- **WHEN** 用户双击或点击托盘图标
- **THEN** 主窗口恢复显示
- **AND** 窗口置于最前

#### Scenario: 通过托盘退出应用

- **WHEN** 用户在托盘菜单中点击「退出」
- **THEN** 应用完全退出
- **AND** 托盘图标消失

#### Scenario: 托盘图标

- **WHEN** 应用启动
- **THEN** 系统托盘区域显示应用图标
- **AND** 图标使用 `build/appicon.png`
