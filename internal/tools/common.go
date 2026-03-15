package tools

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gotools/internal/log"
)

// CommonTool 包含所有工具共享的通用方法
type CommonTool struct{}

// NewCommonTool 创建公共工具实例
func NewCommonTool() *CommonTool {
	return &CommonTool{}
}

// ExecuteCommand 执行命令并返回输出，统一错误处理
func (c *CommonTool) ExecuteCommand(name string, args ...string) (string, error) {
	log.Info("Executing command", "command", name, "args", args)
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		log.Error("Command execution failed", "command", name, "args", args, "error", err, "output", outputStr)
		return "", fmt.Errorf("命令执行失败: %s, 错误信息: %s", err.Error(), outputStr)
	}

	log.Info("Command executed successfully", "command", name)
	return string(output), nil
}

// ExecuteCommandWithStdin 执行命令并提供标准输入
func (c *CommonTool) ExecuteCommandWithStdin(name string, stdinContent string, args ...string) (string, error) {
	log.Info("Executing command with stdin", "command", name)
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(stdinContent)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		log.Error("Command execution with stdin failed", "command", name, "error", err, "output", outputStr)
		return "", fmt.Errorf("命令执行失败: %s, 错误信息: %s", err.Error(), outputStr)
	}

	return string(output), nil
}

// NormalizePath 将路径转换为跨平台兼容格式
func (c *CommonTool) NormalizePath(path string) string {
	// 将反斜杠转换为正斜杠，确保在Windows上与mysqlsh兼容
	return filepath.ToSlash(path)
}

// EnsureDirExists 确保目录存在，不存在则创建
func (c *CommonTool) EnsureDirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Info("Creating directory", "path", dirPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Error("Failed to create directory", "path", dirPath, "error", err)
			return fmt.Errorf("创建目录失败: %s", err.Error())
		}
	}
	return nil
}

// FileExists 检查文件是否存在
func (c *CommonTool) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// RemoveIfExists 如果文件或目录存在则删除
func (c *CommonTool) RemoveIfExists(path string) error {
	if c.FileExists(path) {
		log.Info("Removing existing path", "path", path)
		if err := os.RemoveAll(path); err != nil {
			log.Error("Failed to remove path", "path", path, "error", err)
			return fmt.Errorf("删除路径失败: %s", err.Error())
		}
	}
	return nil
}

// FormatError 格式化错误信息
func (c *CommonTool) FormatError(baseErr error, context string) error {
	if baseErr == nil {
		return nil
	}
	if context != "" {
		return fmt.Errorf("%s: %s", context, baseErr.Error())
	}
	return baseErr
}

// ValidateDatabaseConnection 验证数据库连接参数
func (c *CommonTool) ValidateDatabaseConnection(conn DatabaseConnection) error {
	result := ValidateDatabaseConnection(conn)
	if !result.IsValid {
		return fmt.Errorf("%s", strings.Join(result.GetErrorMessages(), "; "))
	}
	return nil
}

// ValidateExportConfig 验证导出配置
func (c *CommonTool) ValidateExportConfig(outputDir string) error {
	if outputDir == "" {
		return fmt.Errorf("输出目录不能为空")
	}
	return nil
}

// BuildDatabaseURI 构建数据库连接URI
func (c *CommonTool) BuildDatabaseURI(conn DatabaseConnection, includeDatabase bool) string {
	uri := fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port)
	if includeDatabase && conn.Database != "" {
		uri += "/" + conn.Database
	}
	return uri
}

// ParseConnectionString 解析连接字符串
func (c *CommonTool) ParseConnectionString(connectionString string) (DatabaseConnection, error) {
	// 简单的连接字符串解析，格式: user:password@host:port/database
	// 实际项目中可能需要更复杂的解析逻辑
	return DatabaseConnection{}, fmt.Errorf("连接字符串解析功能待实现")
}

// BuildCommandArgs 构建命令行参数
func (c *CommonTool) BuildCommandArgs(baseArgs []string, options map[string]string) []string {
	args := make([]string, len(baseArgs))
	copy(args, baseArgs)

	for key, value := range options {
		if value != "" {
			args = append(args, fmt.Sprintf("--%s=%s", key, value))
		}
	}

	return args
}

// JoinStrings 安全连接字符串
func (c *CommonTool) JoinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	return strings.Join(strs, sep)
}

// ContainsString 检查字符串切片是否包含指定字符串
func (c *CommonTool) ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// RemoveDuplicates 移除字符串切片中的重复项
func (c *CommonTool) RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, str := range slice {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

// ValidateCompressionType 验证压缩类型
func (c *CommonTool) ValidateCompressionType(compression string) bool {
	return containsInSlice(ValidCompressionTypes, strings.ToLower(compression))
}

// GetDefaultExportConfig 获取默认导出配置
func (c *CommonTool) GetDefaultExportConfig() map[string]interface{} {
	return map[string]interface{}{
		"threads":     4,
		"compression": "gzip",
		"overwrite":   true,
		"skipDefiner": true,
	}
}

// ExecuteCommandWithOutputFile 执行命令并重定向输出到文件
func (c *CommonTool) ExecuteCommandWithOutputFile(command string, outputFile string, compression string, args ...string) error {
	log.Info("Executing command with output file", "command", command, "output", outputFile, "compression", compression)

	// 创建输出文件目录
	outputDir := filepath.Dir(outputFile)
	if err := c.EnsureDirExists(outputDir); err != nil {
		return err
	}

	cmd := exec.Command(command, args...)

	var outFile *os.File
	var err error

	if strings.ToLower(compression) == "gzip" {
		outFile, err = os.Create(outputFile)
		if err != nil {
			log.Error("Failed to create output file", "path", outputFile, "error", err)
			return fmt.Errorf("创建输出文件失败: %s", err.Error())
		}
		defer outFile.Close()

		gzipWriter := gzip.NewWriter(outFile)
		defer gzipWriter.Close()

		cmd.Stdout = gzipWriter
	} else {
		cmd.Stdout, err = os.Create(outputFile)
		if err != nil {
			log.Error("Failed to create output file", "path", outputFile, "error", err)
			return fmt.Errorf("创建输出文件失败: %s", err.Error())
		}
		defer cmd.Stdout.(*os.File).Close()
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Error("Failed to get stderr pipe", "error", err)
		return fmt.Errorf("获取标准错误输出失败: %s", err.Error())
	}

	if err := cmd.Start(); err != nil {
		log.Error("Failed to start command", "command", command, "error", err)
		return fmt.Errorf("启动命令失败: %s", err.Error())
	}

	stderrOutput, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		log.Error("Error executing command", "command", command, "error", err, "stderr", string(stderrOutput))
		return fmt.Errorf("命令执行失败: %s, 错误信息: %s", err.Error(), string(stderrOutput))
	}

	log.Info("Command executed successfully", "command", command, "output", outputFile)
	return nil
}

// ExecuteCommandWithInputFile 执行命令并从文件读取输入
func (c *CommonTool) ExecuteCommandWithInputFile(command string, inputFile string, args ...string) (string, error) {
	log.Info("Executing command with input file", "command", command, "input", inputFile)

	if !c.FileExists(inputFile) {
		return "", fmt.Errorf("输入文件不存在: %s", inputFile)
	}

	cmd := exec.Command(command, args...)

	inputFileHandle, err := os.Open(inputFile)
	if err != nil {
		log.Error("Failed to open input file", "path", inputFile, "error", err)
		return "", fmt.Errorf("打开输入文件失败: %s", err.Error())
	}
	defer inputFileHandle.Close()

	isGzip := strings.HasSuffix(strings.ToLower(inputFile), ".gz")

	if isGzip {
		gzipReader, err := gzip.NewReader(inputFileHandle)
		if err != nil {
			log.Error("Failed to create gzip reader", "path", inputFile, "error", err)
			return "", fmt.Errorf("创建解压读取器失败: %s", err.Error())
		}
		defer gzipReader.Close()

		cmd.Stdin = gzipReader
	} else {
		cmd.Stdin = inputFileHandle
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		log.Error("Error executing command with input file", "command", command, "input", inputFile, "error", err, "output", outputStr)
		return "", fmt.Errorf("命令执行失败: %s, 错误信息: %s", err.Error(), outputStr)
	}

	log.Info("Command executed successfully with input file", "command", command, "input", inputFile)
	return string(output), nil
}

// CheckCommandExists 检查命令是否存在
func (c *CommonTool) CheckCommandExists(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		log.Error("Command not found", "command", command, "error", err)
		return fmt.Errorf("未找到命令: %s", command)
	}
	return nil
}
