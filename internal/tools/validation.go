package tools

import (
	"fmt"
	"regexp"
	"strings"
)

// 有效的压缩类型列表
var ValidCompressionTypes = []string{"gzip", "gz", "zstd", "none", ""}

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors"`
}

// NewValidationResult 创建新的验证结果
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}
}

// AddError 添加验证错误
func (r *ValidationResult) AddError(field, message string) {
	r.IsValid = false
	r.Errors = append(r.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors 检查是否有错误
func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// GetErrorMessages 获取所有错误消息
func (r *ValidationResult) GetErrorMessages() []string {
	messages := make([]string, len(r.Errors))
	for i, err := range r.Errors {
		messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}
	return messages
}

// Validator 验证器接口
type Validator interface {
	Validate() *ValidationResult
}

// ValidateDatabaseConnection 验证数据库连接配置
func ValidateDatabaseConnection(conn DatabaseConnection) *ValidationResult {
	result := NewValidationResult()

	// 验证ID
	if conn.ID == "" {
		result.AddError("ID", "连接ID不能为空")
	} else if len(conn.ID) > 50 {
		result.AddError("ID", "连接ID长度不能超过50个字符")
	}

	// 验证名称
	if conn.Name == "" {
		result.AddError("Name", "连接名称不能为空")
	} else if len(conn.Name) > 100 {
		result.AddError("Name", "连接名称长度不能超过100个字符")
	}

	// 验证主机
	if conn.Host == "" {
		result.AddError("Host", "数据库主机不能为空")
	} else if !isValidHost(conn.Host) {
		result.AddError("Host", "数据库主机格式无效")
	}

	// 验证端口
	if conn.Port <= 0 || conn.Port > 65535 {
		result.AddError("Port", "数据库端口必须在1-65535之间")
	}

	// 验证用户
	if conn.User == "" {
		result.AddError("User", "数据库用户不能为空")
	} else if len(conn.User) > 50 {
		result.AddError("User", "数据库用户长度不能超过50个字符")
	}

	// 验证密码（允许为空，但如果提供则检查长度）
	if len(conn.Password) > 100 {
		result.AddError("Password", "数据库密码长度不能超过100个字符")
	}

	// 验证数据库名（可选，但如果提供则检查格式）
	if conn.Database != "" && !isValidDatabaseName(conn.Database) {
		result.AddError("Database", "数据库名称格式无效")
	}

	return result
}

// ValidateExportSettings 验证导出设置
func ValidateExportSettings(settings ExportSettings) *ValidationResult {
	result := NewValidationResult()

	// 验证导出工具
	validTools := []string{"mysql-shell", "mysqldump"}
	if settings.ExportTool != "" && !containsInSlice(validTools, settings.ExportTool) {
		result.AddError("ExportTool", "导出工具必须是 mysql-shell 或 mysqldump")
	}

	// 验证导出路径
	if settings.ExportPath != "" && len(settings.ExportPath) > 500 {
		result.AddError("ExportPath", "导出路径长度不能超过500个字符")
	}

	// 验证 MySQLShell 配置
	if settings.MySQLShell.Threads < 0 || settings.MySQLShell.Threads > 64 {
		result.AddError("Threads", "线程数必须在0-64之间")
	}

	if settings.MySQLShell.Compression != "" && !containsInSlice(ValidCompressionTypes, settings.MySQLShell.Compression) {
		result.AddError("Compression", "压缩类型必须是 gzip、gz、zstd 或 none")
	}

	if settings.MySQLShell.ChunkSize != "" && !isValidChunkSize(settings.MySQLShell.ChunkSize) {
		result.AddError("ChunkSize", "分块大小格式无效，例如：1G、512M")
	}

	// MySQLDump 配置验证（压缩格式）
	if settings.MySQLDump.Compression != "" && !containsInSlice(ValidCompressionTypes, settings.MySQLDump.Compression) {
		result.AddError("Compression", "压缩类型必须是 gzip、gz、zstd 或 none")
	}

	return result
}

// ValidateExportConfig 验证导出配置
func ValidateExportConfig(config ExportConfig) *ValidationResult {
	result := NewValidationResult()

	// 验证输出目录
	if config.OutputDir == "" {
		result.AddError("OutputDir", "输出目录不能为空")
	} else if len(config.OutputDir) > 500 {
		result.AddError("OutputDir", "输出目录长度不能超过500个字符")
	}

	// 验证线程数
	if config.Threads < 0 || config.Threads > 64 {
		result.AddError("Threads", "线程数必须在0-64之间")
	}

	// 验证压缩类型
	if config.Compression != "" && !containsInSlice(ValidCompressionTypes, strings.ToLower(config.Compression)) {
		result.AddError("Compression", "压缩类型必须是 gzip、gz、zstd 或 none")
	}

	// 验证分块大小格式
	if config.ChunkSize != "" && !isValidChunkSize(config.ChunkSize) {
		result.AddError("ChunkSize", "分块大小格式无效，例如：1G、512M")
	}

	return result
}

// ValidateImportConfig 验证导入配置
func ValidateImportConfig(config ImportConfig) *ValidationResult {
	result := NewValidationResult()

	// 验证输入目录
	if config.InputDir == "" {
		result.AddError("InputDir", "输入目录不能为空")
	} else if len(config.InputDir) > 500 {
		result.AddError("InputDir", "输入目录长度不能超过500个字符")
	}

	// 验证线程数
	if config.Threads < 0 || config.Threads > 64 {
		result.AddError("Threads", "线程数必须在0-64之间")
	}

	// 验证等待超时
	if config.WaitTimeout < 0 || config.WaitTimeout > 86400 {
		result.AddError("WaitTimeout", "等待超时必须在0-86400秒之间")
	}

	return result
}

// 辅助函数

// isValidHost 验证主机名格式
func isValidHost(host string) bool {
	// 简单的主机名验证，支持IP地址和域名
	if host == "localhost" || host == "127.0.0.1" {
		return true
	}

	// IP地址格式
	ipPattern := `^(\d{1,3}\.){3}\d{1,3}$`
	if matched, _ := regexp.MatchString(ipPattern, host); matched {
		return true
	}

	// 域名格式
	domainPattern := `^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$`
	if matched, _ := regexp.MatchString(domainPattern, host); matched {
		return true
	}

	return false
}

// isValidDatabaseName 验证数据库名称格式
func isValidDatabaseName(name string) bool {
	// 数据库名称允许字母、数字、下划线和连字符
	pattern := `^[a-zA-Z_][a-zA-Z0-9_-]{0,63}$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}

// isValidChunkSize 验证分块大小格式
func isValidChunkSize(size string) bool {
	// 支持格式：1G、512M、1024K等
	pattern := `^\d+[KMG]?$`
	matched, _ := regexp.MatchString(pattern, strings.ToUpper(size))
	return matched
}

// isValidSchemaList 验证模式列表格式
func isValidSchemaList(schemas string) bool {
	// 逗号分隔的模式列表
	if schemas == "" {
		return true
	}

	parts := strings.Split(schemas, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && !isValidDatabaseName(part) {
			return false
		}
	}
	return true
}

// isValidTableList 验证表列表格式
func isValidTableList(tables string) bool {
	// 逗号分隔的表列表
	if tables == "" {
		return true
	}

	parts := strings.Split(tables, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && !isValidDatabaseName(part) {
			return false
		}
	}
	return true
}

// containsInSlice 检查切片是否包含指定字符串
func containsInSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateAll 验证所有配置
func ValidateAll(conn DatabaseConnection, exportSettings ExportSettings, exportConfig ExportConfig, importConfig ImportConfig) *ValidationResult {
	result := NewValidationResult()

	// 验证数据库连接
	connResult := ValidateDatabaseConnection(conn)
	if !connResult.IsValid {
		for _, err := range connResult.Errors {
			result.AddError("Connection."+err.Field, err.Message)
		}
	}

	// 验证导出设置
	settingsResult := ValidateExportSettings(exportSettings)
	if !settingsResult.IsValid {
		for _, err := range settingsResult.Errors {
			result.AddError("ExportSettings."+err.Field, err.Message)
		}
	}

	// 验证导出配置
	exportConfigResult := ValidateExportConfig(exportConfig)
	if !exportConfigResult.IsValid {
		for _, err := range exportConfigResult.Errors {
			result.AddError("ExportConfig."+err.Field, err.Message)
		}
	}

	// 验证导入配置
	importConfigResult := ValidateImportConfig(importConfig)
	if !importConfigResult.IsValid {
		for _, err := range importConfigResult.Errors {
			result.AddError("ImportConfig."+err.Field, err.Message)
		}
	}

	return result
}
