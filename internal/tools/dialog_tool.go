package tools

import (
	"os/exec"
	"strings"
)

// SelectFolder 打开系统文件夹选择对话框并返回选择的路径
func SelectFolder() (string, error) {
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

// SelectFile 打开系统文件选择对话框并返回选择的文件路径
func SelectFile() (string, error) {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.OpenFileDialog
$dialog.Title = "选择文件"
$dialog.Filter = "所有文件 (*.*)|*.*|SQL文件 (*.sql)|*.sql|压缩文件 (*.gz)|*.gz"
$dialog.FilterIndex = 1
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $dialog.FileName
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

// SelectSaveFile 打开系统保存文件对话框并返回选择的文件路径
func SelectSaveFile(defaultName string) (string, error) {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.SaveFileDialog
$dialog.Title = "保存文件"
$dialog.Filter = "SQL文件 (*.sql)|*.sql|压缩SQL文件 (*.sql.gz)|*.sql.gz|所有文件 (*.*)|*.*"
$dialog.FilterIndex = 1
$dialog.FileName = "` + defaultName + `"
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $dialog.FileName
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
