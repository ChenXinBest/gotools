package tools

import (
	"os/exec"
	"strings"
)

// SelectFolder 打开系统文件夹选择对话框并返回选择的路径
func SelectFolder() (string, error) {
	// 使用 PowerShell 调用 .NET 的 FolderBrowserDialog
	script := `
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.FolderBrowserDialog
$dialog.Description = "选择导出目录"
$dialog.ShowNewFolderButton = $true
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $dialog.SelectedPath
} else {
    Write-Output ""
}
`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", nil
	}

	return path, nil
}
