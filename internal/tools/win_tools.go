// win_tools windows系统操作相关工具
package tools

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	"gotools/internal/log"
)

// ProcessInfo 进程信息结构体
type ProcessInfo struct {
	PID        int32
	Name       string
	Cmdline    string
	CPUPercent float64
	MemoryMB   uint64
	ListenAddr string
	ListenPort uint32
	Status     string
}

// ShowSystemInfo 展示系统进程信息
func ShowSystemInfo() {
	infos, _ := GetSystemProcessInfos()

	for _, info := range infos {
		log.Info("Process info", "pid", info.PID, "name", info.Name, "cmdline", info.Cmdline, "cpu", info.CPUPercent, "memory", info.MemoryMB)
	}
}

// GetSystemProcessInfos 获取系统进程信息列表
func GetSystemProcessInfos() ([]ProcessInfo, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	// 获取CPU核心数
	cpuCount := float64(runtime.NumCPU())

	var infos []ProcessInfo
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, _ := p.Name()
		cmdline, _ := p.Cmdline()
		cpuPercent, _ := p.CPUPercent()
		// 归一化CPU使用率，除以核心数
		normalizedCPU := cpuPercent / cpuCount
		memInfo, _ := p.MemoryInfo()

		var memMB uint64
		if memInfo != nil {
			memMB = memInfo.RSS / 1024 / 1024
		}

		// 获取该进程的端口占用信息
		ports, _ := getProcessPorts(pid)

		infos = append(infos, ProcessInfo{
			PID:        pid,
			Name:       name,
			Cmdline:    cmdline,
			CPUPercent: normalizedCPU,
			MemoryMB:   memMB,
			ListenAddr: ports.listenAddr,
			ListenPort: ports.listenPort,
			Status:     ports.status,
		})
	}

	return infos, nil
}

// processPorts 用于存储进程的端口信息
type processPorts struct {
	listenAddr string
	listenPort uint32
	status     string
}

// getProcessPorts 获取指定进程的端口占用信息
func getProcessPorts(pid int32) (processPorts, error) {
	connections, err := net.ConnectionsPid("tcp", pid)
	if err != nil || len(connections) == 0 {
		return processPorts{}, err
	}

	conn := connections[0]
	return processPorts{
		listenAddr: conn.Laddr.IP,
		listenPort: conn.Laddr.Port,
		status:     conn.Status,
	}, nil
}

// SearchPidByKeyWord 通过进程名、进程命令、PID、端口号等关键信息搜索进程信息
func SearchPidByKeyWord(keyword string) (ProcessInfo, error) {
	infos, err := GetSystemProcessInfos()
	if err != nil {
		return ProcessInfo{}, err
	}

	// 尝试将 keyword 转换为 PID
	var searchPID int32
	n, _ := fmt.Sscanf(keyword, "%d", &searchPID)
	hasValidPID := n == 1

	for _, info := range infos {
		// 匹配进程名
		if keyword != "" && contains(info.Name, keyword) {
			return info, nil
		}
		// 匹配命令行
		if keyword != "" && contains(info.Cmdline, keyword) {
			return info, nil
		}
		// 匹配 PID（仅当 keyword 是有效数字时）
		if hasValidPID && info.PID == searchPID {
			return info, nil
		}
		// 匹配监听地址
		if keyword != "" && contains(info.ListenAddr, keyword) {
			return info, nil
		}
		// 匹配监听端口
		if keyword != "" && fmt.Sprintf("%d", info.ListenPort) == keyword {
			return info, nil
		}
	}

	return ProcessInfo{}, fmt.Errorf("未找到匹配进程: %s", keyword)
}

// KillProcessByPID 根据PID终止指定进程
func KillProcessByPID(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("无法找到进程 %d: %v", pid, err)
	}

	err = p.Kill()
	if err != nil {
		return fmt.Errorf("终止进程 %d 失败: %v", pid, err)
	}

	return nil
}
