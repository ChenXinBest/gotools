package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"gotools/internal/tools"
)

func TestCommonTool_ExecuteCommand(t *testing.T) {
	common := tools.NewCommonTool()

	// 测试执行一个简单的命令
	output, err := common.ExecuteCommand("echo", "Hello, World!")
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v", err)
	}

	t.Logf("Command output: %s", output)
}

func TestCommonTool_NormalizePath(t *testing.T) {
	common := tools.NewCommonTool()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Windows path", "C:\\Users\\test\\file.txt", "C:/Users/test/file.txt"},
		{"Unix path", "/home/user/file.txt", "/home/user/file.txt"},
		{"Mixed path", "C:\\Users/test\\file.txt", "C:/Users/test/file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.NormalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCommonTool_EnsureDirExists(t *testing.T) {
	common := tools.NewCommonTool()

	// 创建临时目录进行测试
	tmpDir := filepath.Join(os.TempDir(), "test_common_tool")
	defer os.RemoveAll(tmpDir)

	err := common.EnsureDirExists(tmpDir)
	if err != nil {
		t.Fatalf("EnsureDirExists failed: %v", err)
	}

	// 验证目录是否存在
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	// 测试目录已存在的情况
	err = common.EnsureDirExists(tmpDir)
	if err != nil {
		t.Fatalf("EnsureDirExists failed for existing directory: %v", err)
	}
}

func TestCommonTool_FileExists(t *testing.T) {
	common := tools.NewCommonTool()

	// 测试不存在的文件
	if common.FileExists("/non/existent/file.txt") {
		t.Error("FileExists should return false for non-existent file")
	}

	// 测试存在的文件
	tmpFile := filepath.Join(os.TempDir(), "test_file.txt")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()
	defer os.Remove(tmpFile)

	if !common.FileExists(tmpFile) {
		t.Error("FileExists should return true for existing file")
	}
}

func TestCommonTool_RemoveIfExists(t *testing.T) {
	common := tools.NewCommonTool()

	// 创建临时文件进行测试
	tmpFile := filepath.Join(os.TempDir(), "test_remove.txt")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()

	// 删除文件
	err = common.RemoveIfExists(tmpFile)
	if err != nil {
		t.Fatalf("RemoveIfExists failed: %v", err)
	}

	// 验证文件已被删除
	if common.FileExists(tmpFile) {
		t.Error("File should have been removed")
	}

	// 测试文件不存在的情况
	err = common.RemoveIfExists(tmpFile)
	if err != nil {
		t.Fatalf("RemoveIfExists failed for non-existent file: %v", err)
	}
}

func TestCommonTool_ValidateDatabaseConnection(t *testing.T) {
	common := tools.NewCommonTool()

	tests := []struct {
		name    string
		conn    tools.DatabaseConnection
		wantErr bool
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
			wantErr: false,
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
			wantErr: true,
		},
		{
			name: "Invalid port",
			conn: tools.DatabaseConnection{
				ID:   "test-id",
				Name: "Test Connection",
				Host: "localhost",
				Port: 0,
				User: "root",
			},
			wantErr: true,
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.ValidateDatabaseConnection(tt.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDatabaseConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommonTool_BuildDatabaseURI(t *testing.T) {
	common := tools.NewCommonTool()

	conn := tools.DatabaseConnection{
		User:     "root",
		Password: "password",
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
	}

	// 测试包含数据库的URI
	uri := common.BuildDatabaseURI(conn, true)
	expected := "root:password@localhost:3306/testdb"
	if uri != expected {
		t.Errorf("BuildDatabaseURI() = %q, want %q", uri, expected)
	}

	// 测试不包含数据库的URI
	uri = common.BuildDatabaseURI(conn, false)
	expected = "root:password@localhost:3306"
	if uri != expected {
		t.Errorf("BuildDatabaseURI() = %q, want %q", uri, expected)
	}
}

func TestCommonTool_CheckCommandExists(t *testing.T) {
	common := tools.NewCommonTool()

	// 测试存在的命令
	err := common.CheckCommandExists("go")
	if err != nil {
		t.Logf("Note: 'go' command not found (this might be expected): %v", err)
	}

	// 测试不存在的命令
	err = common.CheckCommandExists("non_existent_command_12345")
	if err == nil {
		t.Error("CheckCommandExists should return error for non-existent command")
	}
}

func TestCommonTool_ValidateCompressionType(t *testing.T) {
	common := tools.NewCommonTool()

	tests := []struct {
		name        string
		compression string
		expected    bool
	}{
		{"gzip", "gzip", true},
		{"gz", "gz", true},
		{"zstd", "zstd", true},
		{"none", "none", true},
		{"empty", "", true},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.ValidateCompressionType(tt.compression)
			if result != tt.expected {
				t.Errorf("ValidateCompressionType(%q) = %v, want %v", tt.compression, result, tt.expected)
			}
		})
	}
}

func TestCommonTool_GetDefaultExportConfig(t *testing.T) {
	common := tools.NewCommonTool()

	config := common.GetDefaultExportConfig()

	// 验证默认值
	if threads, ok := config["threads"].(int); !ok || threads != 4 {
		t.Errorf("Default threads should be 4, got %v", config["threads"])
	}

	if compression, ok := config["compression"].(string); !ok || compression != "gzip" {
		t.Errorf("Default compression should be 'gzip', got %v", config["compression"])
	}

	if overwrite, ok := config["overwrite"].(bool); !ok || !overwrite {
		t.Errorf("Default overwrite should be true, got %v", config["overwrite"])
	}

	if skipDefiner, ok := config["skipDefiner"].(bool); !ok || !skipDefiner {
		t.Errorf("Default skipDefiner should be true, got %v", config["skipDefiner"])
	}
}
