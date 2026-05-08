package main

import (
	"gotools/internal/log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// showExportProgressWindow 创建并显示导出进度独立窗口
func showExportProgressWindow(app *application.App) *application.WebviewWindow {
	progressWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "导出进度",
		Width:  500,
		Height: 300,
		// 固定窗口大小，避免用户调整
		DisableResize: true,
		MinWidth:      400,
		MinHeight:     250,
		MaxWidth:      600,
		MaxHeight:     400,
		URL:           "/export-progress",
	})
	progressWin.Show()

	log.Info("Export progress window created")
	return progressWin
}

// showImportConflictWindow 创建并显示导入冲突检测独立窗口
func showImportConflictWindow(app *application.App) *application.WebviewWindow {
	conflictWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "导入冲突检测",
		Width:  650,
		Height: 450,
		MinWidth:  500,
		MinHeight: 300,
		URL:   "/import-conflicts",
	})
	conflictWin.Show()

	log.Info("Import conflict window created")
	return conflictWin
}

// setupMultiWindowEvents 注册多窗口相关事件处理
func setupMultiWindowEvents(app *application.App) {
	// 监听前端请求打开导出进度窗口
	app.Event.On("open-export-progress", func(event *application.CustomEvent) {
		showExportProgressWindow(app)
	})

	// 监听前端请求打开导入冲突窗口
	app.Event.On("open-import-conflicts", func(event *application.CustomEvent) {
		showImportConflictWindow(app)
	})

	// 监听导出进度更新事件，转发到前端
	app.Event.On("export-progress-update", func(event *application.CustomEvent) {
		log.Info("Export progress update", "data", event.Data)
	})

	// 监听关闭所有辅助窗口事件
	app.Event.On("close-child-windows", func(event *application.CustomEvent) {
		for _, win := range app.Window.GetAll() {
			// 保留主窗口（通过 Name 判断）
			if win.Name() != "" && win.Name() != "main" {
				win.Close()
			}
		}
	})
}
