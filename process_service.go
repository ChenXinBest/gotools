package main

import (
	"context"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/services"
	"gotools/internal/tools"
)

// ProcessService 进程管理 Wails 服务
type ProcessService struct {
	ctx    context.Context
	inner *services.ProcessService
}

// NewProcessService 创建进程管理服务
func NewProcessService() *ProcessService {
	return &ProcessService{
		inner: services.NewProcessService(),
	}
}

// Startup 应用启动时由 Wails 调用
func (s *ProcessService) Startup(ctx context.Context) {
	s.ctx = ctx
	log.Info("ProcessService started")
}

// GetSystemProcessInfos 获取系统进程信息列表
func (s *ProcessService) GetSystemProcessInfos() ([]tools.ProcessInfo, error) {
	return s.inner.GetSystemProcessInfos()
}

// SearchPidByKeyWord 搜索进程
func (s *ProcessService) SearchPidByKeyWord(keyword string) (tools.ProcessInfo, error) {
	return s.inner.SearchPidByKeyWord(keyword)
}

// KillProcessByPID 终止进程
func (s *ProcessService) KillProcessByPID(pid int32) error {
	return s.inner.KillProcessByPID(pid)
}

// Greet 问候用户（示例方法，保留兼容性）
func (s *ProcessService) Greet(name string) string {
	log.Info("Greeting user", "name", name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
