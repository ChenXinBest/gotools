package main

import (
	"context"

	"gotools/internal/log"
	"gotools/internal/services"
)

// DialogService 对话框 Wails 服务
type DialogService struct {
	ctx    context.Context
	inner *services.DialogService
}

// NewDialogService 创建对话框服务
func NewDialogService() *DialogService {
	return &DialogService{
		inner: services.NewDialogService(),
	}
}

// Startup 应用启动时由 Wails 调用
func (s *DialogService) Startup(ctx context.Context) {
	s.ctx = ctx
	log.Info("DialogService started")
}

// SelectFolder 打开系统文件夹选择对话框
func (s *DialogService) SelectFolder() (string, error) {
	return s.inner.SelectFolder()
}

// SelectFile 打开系统文件选择对话框
func (s *DialogService) SelectFile() (string, error) {
	return s.inner.SelectFile()
}

// SelectSaveFile 打开系统保存文件对话框
func (s *DialogService) SelectSaveFile(defaultName string) (string, error) {
	return s.inner.SelectSaveFile(defaultName)
}
