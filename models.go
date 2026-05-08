package main

import "gotools/internal/tools"

// ExportRequest 导出请求结构
type ExportRequest struct {
	ConnectionID   string   `json:"connection_id"`
	Databases      []string `json:"databases"`
	Database       string   `json:"database"`
	Tables         []string `json:"tables"`
	OutputDir      string   `json:"output_dir"`
	Threads        int      `json:"threads"`
	Compression    string   `json:"compression"`
	ChunkSize      string   `json:"chunk_size"`
	SkipDefiner    bool     `json:"skip_definer"`
	SkipBinlog     bool     `json:"skip_binlog"`
	IncludeSchemas []string `json:"include_schemas"`
	ExcludeSchemas []string `json:"exclude_schemas"`
	IncludeTables  []string `json:"include_tables"`
	ExcludeTables  []string `json:"exclude_tables"`
	Overwrite      bool     `json:"overwrite"`
}

// ExportResponse 导出响应结构
type ExportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

// ImportRequest 导入请求结构
type ImportRequest struct {
	ConnectionID   string   `json:"connection_id"`
	Database       string   `json:"database"`
	InputDir       string   `json:"input_dir"`
	Threads        int      `json:"threads"`
	Schema         string   `json:"schema"`
	IncludeSchemas []string `json:"include_schemas"`
	ExcludeSchemas []string `json:"exclude_schemas"`
	IncludeTables  []string `json:"include_tables"`
	ExcludeTables  []string `json:"exclude_tables"`
	ResetProgress  bool     `json:"reset_progress"`
	WaitTimeout    int      `json:"wait_timeout"`
}

// ImportResponse 导入响应结构
type ImportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

// CheckImportConflictsRequest 检查导入冲突请求结构
type CheckImportConflictsRequest struct {
	ConnectionID string `json:"connection_id"`
	InputDir     string `json:"input_dir"`
}

// DropConflictingTablesRequest 删除冲突表请求结构
type DropConflictingTablesRequest struct {
	ConnectionID string                 `json:"connection_id"`
	Conflicts    []tools.ImportConflict `json:"conflicts"`
}

// MySQLDumpExportRequest mysqldump导出请求结构
type MySQLDumpExportRequest struct {
	ConnectionID      string   `json:"connection_id"`
	Databases         []string `json:"databases"`
	Database          string   `json:"database"`
	Tables            []string `json:"tables"`
	OutputDir         string   `json:"output_dir"`
	Compression       string   `json:"compression"`
	SingleTransaction bool     `json:"single_transaction"`
	Routines          bool     `json:"routines"`
	Events            bool     `json:"events"`
	Overwrite         bool     `json:"overwrite"`
}

// MySQLDumpImportRequest mysqldump导入请求结构
type MySQLDumpImportRequest struct {
	ConnectionID string `json:"connection_id"`
	InputFile    string `json:"input_file"`
	Database     string `json:"database"`
}
