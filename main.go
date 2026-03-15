package main

import (
	"embed"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/version"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
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

	// Create an instance of the app structure
	app := NewApp()

	// 构建应用标题
	appTitle := fmt.Sprintf("GoTools %s", version.Version)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  appTitle,
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 10, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		// Windows 特定选项
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		// macOS 特定选项
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarDefault(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		log.Error("Error running app", "error", err)
	}
}
