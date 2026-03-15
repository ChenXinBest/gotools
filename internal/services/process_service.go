package services

import (
	"gotools/internal/log"
	"gotools/internal/tools"
)

// ProcessService 进程管理服务
type ProcessService struct{}

// NewProcessService 创建进程管理服务实例
func NewProcessService() *ProcessService {
	log.Info("Creating new ProcessService instance")
	return &ProcessService{}
}

// GetSystemProcessInfos 获取系统进程信息列表
func (s *ProcessService) GetSystemProcessInfos() ([]tools.ProcessInfo, error) {
	log.Info("Getting system process infos")
	processes, err := tools.GetSystemProcessInfos()
	if err != nil {
		log.Error("Error getting system process infos", "error", err)
		return nil, err
	}
	log.Info("Got system process infos", "count", len(processes))
	return processes, nil
}

// SearchPidByKeyWord 通过进程名、进程命令、PID、端口号等关键信息搜索进程信息
func (s *ProcessService) SearchPidByKeyWord(keyword string) (tools.ProcessInfo, error) {
	log.Info("Searching process by keyword", "keyword", keyword)
	process, err := tools.SearchPidByKeyWord(keyword)
	if err != nil {
		log.Error("Error searching process by keyword", "keyword", keyword, "error", err)
		return tools.ProcessInfo{}, err
	}
	log.Info("Found process by keyword", "keyword", keyword, "pid", process.PID)
	return process, nil
}

// KillProcessByPID 根据PID终止指定进程
func (s *ProcessService) KillProcessByPID(pid int32) error {
	log.Info("Killing process by PID", "pid", pid)
	err := tools.KillProcessByPID(pid)
	if err != nil {
		log.Error("Error killing process", "pid", pid, "error", err)
		return err
	}
	log.Info("Killed process successfully", "pid", pid)
	return nil
}
