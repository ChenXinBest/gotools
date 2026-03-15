package version

import (
	"fmt"
	"runtime"
)

// 版本信息 - 通过构建时注入
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
	GoVersion = runtime.Version()
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// Info 版本信息结构
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GitBranch string `json:"git_branch"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// GetInfo 获取版本信息
func GetInfo() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		GoVersion: GoVersion,
		Platform:  Platform,
	}
}

// GetVersionString 获取版本字符串
func GetVersionString() string {
	commit := GitCommit
	if len(commit) > 8 {
		commit = commit[:8]
	}
	return fmt.Sprintf("GoTools %s (%s) - %s - Built with %s",
		Version, commit, Platform, GoVersion)
}

// PrintVersion 打印版本信息
func PrintVersion() {
	info := GetInfo()
	fmt.Printf("GoTools %s\n", info.Version)
	fmt.Printf("  Build Time: %s\n", info.BuildTime)
	fmt.Printf("  Git Commit: %s\n", info.GitCommit)
	fmt.Printf("  Git Branch: %s\n", info.GitBranch)
	fmt.Printf("  Go Version: %s\n", info.GoVersion)
	fmt.Printf("  Platform:   %s\n", info.Platform)
}
