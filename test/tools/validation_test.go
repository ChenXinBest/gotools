package tools_test

import (
	"testing"

	"gotools/internal/config"
	"gotools/internal/tools"
)

func TestValidateDatabaseConnection(t *testing.T) {
	tests := []struct {
		name       string
		conn       tools.DatabaseConnection
		wantValid  bool
		errorCount int
	}{
		{
			name: "Valid connection",
			conn: tools.DatabaseConnection{
				ID:       "test-id",
				Name:     "Test Connection",
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "password",
				Database: "testdb",
			},
			wantValid:  true,
			errorCount: 0,
		},
		{
			name: "Empty ID",
			conn: tools.DatabaseConnection{
				ID:   "",
				Name: "Test Connection",
				Host: "localhost",
				Port: 3306,
				User: "root",
			},
			wantValid:  false,
			errorCount: 1,
		},
		{
			name: "Empty name",
			conn: tools.DatabaseConnection{
				ID:   "test-id",
				Name: "",
				Host: "localhost",
				Port: 3306,
				User: "root",
			},
			wantValid:  false,
			errorCount: 1,
		},
		{
			name: "Empty host",
			conn: tools.DatabaseConnection{
				ID:   "test-id",
				Name: "Test Connection",
				Host: "",
				Port: 3306,
				User: "root",
			},
			wantValid:  false,
			errorCount: 1,
		},
		{
			name: "Invalid port (0)",
			conn: tools.DatabaseConnection{
				ID:   "test-id",
				Name: "Test Connection",
				Host: "localhost",
				Port: 0,
				User: "root",
			},
			wantValid:  false,
			errorCount: 1,
		},
		{
			name: "Invalid port (70000)",
			conn: tools.DatabaseConnection{
				ID:   "test-id",
				Name: "Test Connection",
				Host: "localhost",
				Port: 70000,
				User: "root",
			},
			wantValid:  false,
			errorCount: 1,
		},
		{
			name: "Empty user",
			conn: tools.DatabaseConnection{
				ID:   "test-id",
				Name: "Test Connection",
				Host: "localhost",
				Port: 3306,
				User: "",
			},
			wantValid:  false,
			errorCount: 1,
		},
		{
			name: "Multiple errors",
			conn: tools.DatabaseConnection{
				ID:   "",
				Name: "",
				Host: "",
				Port: 0,
				User: "",
			},
			wantValid:  false,
			errorCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.ValidateDatabaseConnection(tt.conn)
			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateDatabaseConnection().IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
			if len(result.Errors) != tt.errorCount {
				t.Errorf("ValidateDatabaseConnection().Errors count = %v, want %v", len(result.Errors), tt.errorCount)
			}
		})
	}
}

func TestValidateExportSettings(t *testing.T) {
	tests := []struct {
		name      string
		settings  tools.ExportSettings
		wantValid bool
	}{
		{
			name: "Valid settings",
			settings: tools.ExportSettings{
				ExportTool: "mysql-shell",
				ExportPath: "/tmp/export",
				MySQLShell: config.MySQLShellConfig{
					Threads:     4,
					Compression: "gzip",
					ChunkSize:   "1G",
				},
			},
			wantValid: true,
		},
		{
			name: "Invalid export tool",
			settings: tools.ExportSettings{
				ExportTool: "invalid-tool",
			},
			wantValid: false,
		},
		{
			name: "Invalid threads (negative)",
			settings: tools.ExportSettings{
				MySQLShell: config.MySQLShellConfig{
					Threads: -1,
				},
			},
			wantValid: false,
		},
		{
			name: "Invalid threads (too high)",
			settings: tools.ExportSettings{
				MySQLShell: config.MySQLShellConfig{
					Threads: 100,
				},
			},
			wantValid: false,
		},
		{
			name: "Invalid compression",
			settings: tools.ExportSettings{
				MySQLShell: config.MySQLShellConfig{
					Compression: "invalid",
				},
			},
			wantValid: false,
		},
		{
			name: "Invalid chunk size",
			settings: tools.ExportSettings{
				MySQLShell: config.MySQLShellConfig{
					ChunkSize: "invalid",
				},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.ValidateExportSettings(tt.settings)
			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateExportSettings().IsValid = %v, want %v", result.IsValid, tt.wantValid)
				if !result.IsValid {
					t.Logf("Errors: %v", result.GetErrorMessages())
				}
			}
		})
	}
}

func TestValidateExportConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    tools.ExportConfig
		wantValid bool
	}{
		{
			name: "Valid config",
			config: tools.ExportConfig{
				OutputDir:   "/tmp/export",
				Threads:     4,
				Compression: "gzip",
				ChunkSize:   "1G",
			},
			wantValid: true,
		},
		{
			name: "Empty output dir",
			config: tools.ExportConfig{
				OutputDir: "",
			},
			wantValid: false,
		},
		{
			name: "Invalid threads",
			config: tools.ExportConfig{
				OutputDir: "/tmp/export",
				Threads:   -1,
			},
			wantValid: false,
		},
		{
			name: "Invalid compression",
			config: tools.ExportConfig{
				OutputDir:   "/tmp/export",
				Compression: "invalid",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.ValidateExportConfig(tt.config)
			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateExportConfig().IsValid = %v, want %v", result.IsValid, tt.wantValid)
				if !result.IsValid {
					t.Logf("Errors: %v", result.GetErrorMessages())
				}
			}
		})
	}
}

func TestValidateImportConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    tools.ImportConfig
		wantValid bool
	}{
		{
			name: "Valid config",
			config: tools.ImportConfig{
				InputDir:    "/tmp/import",
				Threads:     4,
				WaitTimeout: 300,
			},
			wantValid: true,
		},
		{
			name: "Empty input dir",
			config: tools.ImportConfig{
				InputDir: "",
			},
			wantValid: false,
		},
		{
			name: "Invalid threads",
			config: tools.ImportConfig{
				InputDir: "/tmp/import",
				Threads:  -1,
			},
			wantValid: false,
		},
		{
			name: "Invalid wait timeout",
			config: tools.ImportConfig{
				InputDir:    "/tmp/import",
				WaitTimeout: 90000,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.ValidateImportConfig(tt.config)
			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateImportConfig().IsValid = %v, want %v", result.IsValid, tt.wantValid)
				if !result.IsValid {
					t.Logf("Errors: %v", result.GetErrorMessages())
				}
			}
		})
	}
}

func TestValidationResult(t *testing.T) {
	result := tools.NewValidationResult()

	// 初始状态应该是有效的
	if !result.IsValid {
		t.Error("NewValidationResult should be valid by default")
	}

	// 添加错误
	result.AddError("field1", "error message 1")

	// 应该无效
	if result.IsValid {
		t.Error("ValidationResult should be invalid after adding error")
	}

	// 应该有错误
	if !result.HasErrors() {
		t.Error("ValidationResult should have errors")
	}

	// 验证错误消息
	messages := result.GetErrorMessages()
	if len(messages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(messages))
	}

	// 添加另一个错误
	result.AddError("field2", "error message 2")

	messages = result.GetErrorMessages()
	if len(messages) != 2 {
		t.Errorf("Expected 2 error messages, got %d", len(messages))
	}
}

func TestValidateAll(t *testing.T) {
	conn := tools.DatabaseConnection{
		ID:       "test-id",
		Name:     "Test Connection",
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		Database: "testdb",
	}

	exportSettings := tools.ExportSettings{
		ExportTool: "mysql-shell",
		ExportPath: "/tmp/export",
		MySQLShell: config.MySQLShellConfig{
			Threads:     4,
			Compression: "gzip",
			ChunkSize:   "1G",
		},
	}

	exportConfig := tools.ExportConfig{
		OutputDir:   "/tmp/export",
		Threads:     4,
		Compression: "gzip",
		ChunkSize:   "1G",
	}

	importConfig := tools.ImportConfig{
		InputDir:    "/tmp/import",
		Threads:     4,
		WaitTimeout: 300,
	}

	// 测试有效的所有配置
	result := tools.ValidateAll(conn, exportSettings, exportConfig, importConfig)
	if !result.IsValid {
		t.Errorf("ValidateAll should return valid result, errors: %v", result.GetErrorMessages())
	}

	// 测试无效的配置
	invalidConn := tools.DatabaseConnection{
		ID:   "",
		Name: "",
		Host: "",
		Port: 0,
		User: "",
	}

	invalidSettings := tools.ExportSettings{
		ExportTool: "invalid",
		MySQLShell: config.MySQLShellConfig{
			Compression: "invalid",
			Threads:     -1,
		},
	}

	invalidExportConfig := tools.ExportConfig{
		OutputDir: "",
		Threads:   -1,
	}

	invalidImportConfig := tools.ImportConfig{
		InputDir:    "",
		WaitTimeout: 90000,
	}

	result = tools.ValidateAll(invalidConn, invalidSettings, invalidExportConfig, invalidImportConfig)
	if result.IsValid {
		t.Error("ValidateAll should return invalid result for invalid configs")
	}

	if len(result.Errors) == 0 {
		t.Error("ValidateAll should return errors for invalid configs")
	}

	t.Logf("Validation errors: %v", result.GetErrorMessages())
}
