package tools_test

import (
	"testing"

	"gotools/internal/tools"
)

func TestShowSystemInfo(t *testing.T) {
	// 测试函数能够正常执行而不崩溃
	tools.ShowSystemInfo()
}

func TestGetSystemProcessInfos(t *testing.T) {
	// 测试获取进程信息列表
	infos, err := tools.GetSystemProcessInfos()
	if err != nil {
		t.Fatalf("GetSystemProcessInfos failed: %v", err)
	}

	if len(infos) == 0 {
		t.Log("No processes found")
	}

	// 验证返回的进程信息结构
	for _, info := range infos {
		t.Logf("PID: %d, Name: %s, Laddr: %s, Raddr: %s, Status: %s",
			info.PID, info.Name, info.Laddr, info.Raddr, info.Status)
	}
}

func TestSearchPidByKeyWord(t *testing.T) {
	// 先获取系统进程信息
	infos, err := tools.GetSystemProcessInfos()
	if err != nil {
		t.Fatalf("GetSystemProcessInfos failed: %v", err)
	}

	// 测试1: 搜索不存在的关键词，应该返回错误
	_, err = tools.SearchPidByKeyWord("this_process_definitely_not_exist_12345")
	if err == nil {
		t.Error("SearchPidByKeyWord should return error for non-existent keyword")
	} else {
		t.Logf("Expected error for non-existent keyword: %v", err)
	}

	// 测试2: 如果有进程信息，尝试搜索真实进程名
	if len(infos) > 0 {
		// 找到第一个有名称的进程
		var realProcessName string
		for _, info := range infos {
			if info.Name != "" {
				realProcessName = info.Name
				break
			}
		}

		if realProcessName != "" {
			// 搜索真实进程名
			info, err := tools.SearchPidByKeyWord(realProcessName)
			if err != nil {
				t.Errorf("SearchPidByKeyWord failed for existing process %s: %v", realProcessName, err)
			} else {
				t.Logf("Found process: PID=%d, Name=%s", info.PID, info.Name)
			}
		}
	}

	// 测试3: 搜索空字符串（应该返回第一个匹配的进程或错误）
	_, err = tools.SearchPidByKeyWord("")
	if err != nil {
		t.Logf("Expected behavior for empty string: %v", err)
	}
}
