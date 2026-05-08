package main

import (
	"os"
	"runtime"

	"gotools/internal/log"
	"gotools/internal/version"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// setupSystemTray 创建并配置系统托盘
func setupSystemTray(app *application.App, mainWindow application.Window) {
	iconPath := "build/appicon.png"
	if runtime.GOOS == "windows" {
		iconPath = "build/windows/icon.png"
		// Fallback to appicon.png if windows icon doesn't exist
		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			iconPath = "build/appicon.png"
		}
	}

	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		log.Warn("Could not load tray icon", "path", iconPath, "error", err)
		iconData = nil
	}

	// 创建托盘菜单
	trayMenu := app.NewMenu()
	trayMenu.Add("显示窗口").OnClick(func(ctx *application.Context) {
		mainWindow.Show()
	})
	trayMenu.AddSeparator()
	trayMenu.Add("退出").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// 创建系统托盘
	tray := app.SystemTray.New()
	tray.SetLabel("GoTools")
	tray.SetTooltip("GoTools - " + version.Version)
	if iconData != nil {
		tray.SetIcon(iconData)
	}
	tray.SetMenu(trayMenu)
	tray.AttachWindow(mainWindow)

	log.Info("System tray created")
}

// setupAppMenu 创建原生应用菜单
func setupAppMenu(app *application.App) *application.Menu {
	menu := app.NewMenu()

	// 文件菜单
	fileMenu := menu.AddSubmenu("文件")
	fileMenu.Add("退出").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// 视图菜单
	viewMenu := menu.AddSubmenu("视图")
	viewMenu.Add("进程管理器").OnClick(func(ctx *application.Context) {
		app.Event.Emit("navigate", "/process-manager")
	})
	viewMenu.Add("数据库备份").OnClick(func(ctx *application.Context) {
		app.Event.Emit("navigate", "/database-backup")
	})
	viewMenu.Add("设置").OnClick(func(ctx *application.Context) {
		app.Event.Emit("navigate", "/settings")
	})
	viewMenu.AddSeparator()
	viewMenu.Add("刷新").OnClick(func(ctx *application.Context) {
		app.Event.Emit("refresh")
	})

	// 帮助菜单
	helpMenu := menu.AddSubmenu("帮助")
	helpMenu.Add("关于 GoTools").OnClick(func(ctx *application.Context) {
		app.Event.Emit("show-about", version.GetVersionString())
	})

	log.Info("Application menu created")
	return menu
}

// createContextMenu 为进程管理器创建右键上下文菜单
func createContextMenu(app *application.App) *application.Menu {
	contextMenu := app.NewMenu()
	contextMenu.Add("终止进程").OnClick(func(ctx *application.Context) {
		// 通过事件通知前端终止选中的进程
		app.Event.Emit("kill-selected")
	})

	return contextMenu
}
