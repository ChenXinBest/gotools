// win_tools windows系统操作相关工具
package tools

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 进程信息结构体
type ProcessInfo struct {
	PID        int32
	Name       string
	Cmdline    string
	CPUPercent float64
	MemoryMB   uint64
	Laddr      string
	Raddr      string
	Status     string
}

// ShowSystemInfo 展示系统进程信息
func ShowSystemInfo() {
	infos, _ := GetSystemProcessInfos()

	for _, info := range infos {
		fmt.Printf("PID: %d, 名称: %s, 命令行: %s, CPU: %.2f%%, 内存: %v MB\n",
			info.PID, info.Name, info.Cmdline, info.CPUPercent, info.MemoryMB)
	}
}

// GetSystemProcessInfos 获取系统进程信息列表
func GetSystemProcessInfos() ([]ProcessInfo, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	var infos []ProcessInfo
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, _ := p.Name()
		cmdline, _ := p.Cmdline()
		cpuPercent, _ := p.CPUPercent()
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
			CPUPercent: cpuPercent,
			MemoryMB:   memMB,
			Laddr:      ports.laddr,
			Raddr:      ports.raddr,
			Status:     ports.status,
		})
	}

	return infos, nil
}

// processPorts 用于存储进程的端口信息
type processPorts struct {
	laddr  string
	raddr  string
	status string
}

// getProcessPorts 获取指定进程的端口占用信息
func getProcessPorts(pid int32) (processPorts, error) {
	connections, err := net.ConnectionsPid("tcp", pid)
	if err != nil || len(connections) == 0 {
		return processPorts{}, err
	}

	// 只取第一个连接的信息
	conn := connections[0]
	return processPorts{
		laddr:  conn.Laddr.String(),
		raddr:  conn.Raddr.String(),
		status: conn.Status,
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
		// 匹配本地端口
		if keyword != "" && contains(info.Laddr, keyword) {
			return info, nil
		}
		// 匹配远程端口
		if keyword != "" && contains(info.Raddr, keyword) {
			return info, nil
		}
	}

	return ProcessInfo{}, fmt.Errorf("未找到匹配进程: %s", keyword)
}
