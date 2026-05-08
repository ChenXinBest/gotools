package main

import (
	"embed"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/version"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化日志系统
	if err := log.Init(); err != nil {
		println("Error initializing log:", err.Error())
		return
	}

	// 记录启动信息
	log.Info("Starting GoTools", "version", version.Version, "platform", version.Platform)

	// 创建服务实例
	processService := NewProcessService()
	databaseService := NewDatabaseService()
	dialogService := NewDialogService()

	// 构建应用标题
	appTitle := fmt.Sprintf("GoTools %s", version.Version)

	// Create a new Wails v3 application
	app := application.New(application.Options{
		Name:        "GoTools",
		Description: "A collection of useful tools for developers",
		Services: []application.Service{
			application.NewService(processService),
			application.NewService(databaseService),
			application.NewService(dialogService),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false, // 允许最小化到托盘
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true, // 关闭窗口时不退出应用
		},
	})

	// Create the main window
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "main",
		Title:  appTitle,
		Width:  1200,
		Height: 800,
		BackgroundColour: application.NewRGB(10, 10, 10),
		URL:    "/",
	})

	// 创建系统托盘
	setupSystemTray(app, mainWindow)

	// 创建并设置应用菜单
	appMenu := setupAppMenu(app)
	app.Menu.SetApplicationMenu(appMenu)

	// 注册多窗口事件
	setupMultiWindowEvents(app)

	log.Info("GoTools started successfully")

	// Run the application
	err := app.Run()

	if err != nil {
		log.Error("Error running app", "error", err)
	}
}
