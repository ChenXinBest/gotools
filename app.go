package main

import (
	"context"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/tools"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	log.Info("Creating new app instance")
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Info("App startup")
}

func (a *App) Greet(name string) string {
	log.Info("Greeting user", "name", name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetSystemProcessInfos() ([]tools.ProcessInfo, error) {
	log.Info("Getting system process infos")
	processes, err := tools.GetSystemProcessInfos()
	if err != nil {
		log.Error("Error getting system process infos", "error", err)
		return nil, err
	}
	log.Info("Got system process infos", "count", len(processes))
	return processes, nil
}

func (a *App) SearchPidByKeyWord(keyword string) (tools.ProcessInfo, error) {
	log.Info("Searching process by keyword", "keyword", keyword)
	process, err := tools.SearchPidByKeyWord(keyword)
	if err != nil {
		log.Error("Error searching process by keyword", "keyword", keyword, "error", err)
		return tools.ProcessInfo{}, err
	}
	log.Info("Found process by keyword", "keyword", keyword, "pid", process.PID)
	return process, nil
}

func (a *App) KillProcessByPID(pid int32) error {
	log.Info("Killing process by PID", "pid", pid)
	err := tools.KillProcessByPID(pid)
	if err != nil {
		log.Error("Error killing process", "pid", pid, "error", err)
		return err
	}
	log.Info("Killed process successfully", "pid", pid)
	return nil
}

// GetDatabaseConnections 获取所有数据库连接信息
func (a *App) GetDatabaseConnections() ([]tools.DatabaseConnection, error) {
	log.Info("Getting all database connections")
	connections, err := tools.GetDatabaseConnections()
	if err != nil {
		log.Error("Error getting database connections", "error", err)
		return nil, err
	}
	log.Info("Got database connections", "count", len(connections))
	return connections, nil
}

// GetDatabaseConnection 根据 ID 获取数据库连接信息
func (a *App) GetDatabaseConnection(id string) (tools.DatabaseConnection, error) {
	log.Info("Getting database connection by ID", "id", id)
	connection, err := tools.GetDatabaseConnection(id)
	if err != nil {
		log.Error("Error getting database connection", "id", id, "error", err)
		return tools.DatabaseConnection{}, err
	}
	log.Info("Got database connection", "id", id, "name", connection.Name)
	return connection, nil
}

// AddDatabaseConnection 添加数据库连接信息
func (a *App) AddDatabaseConnection(conn tools.DatabaseConnection) error {
	log.Info("Adding database connection", "name", conn.Name)
	err := tools.AddDatabaseConnection(conn)
	if err != nil {
		log.Error("Error adding database connection", "name", conn.Name, "error", err)
		return err
	}
	log.Info("Added database connection successfully", "name", conn.Name)
	return nil
}

// UpdateDatabaseConnection 更新数据库连接信息
func (a *App) UpdateDatabaseConnection(conn tools.DatabaseConnection) error {
	log.Info("Updating database connection", "id", conn.ID, "name", conn.Name)
	err := tools.UpdateDatabaseConnection(conn)
	if err != nil {
		log.Error("Error updating database connection", "id", conn.ID, "name", conn.Name, "error", err)
		return err
	}
	log.Info("Updated database connection successfully", "id", conn.ID, "name", conn.Name)
	return nil
}

// DeleteDatabaseConnection 删除数据库连接信息
func (a *App) DeleteDatabaseConnection(id string) error {
	log.Info("Deleting database connection", "id", id)
	err := tools.DeleteDatabaseConnection(id)
	if err != nil {
		log.Error("Error deleting database connection", "id", id, "error", err)
		return err
	}
	log.Info("Deleted database connection successfully", "id", id)
	return nil
}

// ConnectDatabase 使用 mysqlshell 连接数据库
func (a *App) ConnectDatabase(conn tools.DatabaseConnection) error {
	log.Info("Connecting to database", "name", conn.Name, "host", conn.Host)
	mysqlShellTool := tools.NewMySQLShellTool()
	err := mysqlShellTool.ConnectDatabase(conn)
	if err != nil {
		log.Error("Error connecting to database", "name", conn.Name, "host", conn.Host, "error", err)
		return err
	}
	log.Info("Connected to database successfully", "name", conn.Name, "host", conn.Host)
	return nil
}

// ListDatabases 查询数据库列表
func (a *App) ListDatabases(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing databases", "host", conn.Host)
	mysqlShellTool := tools.NewMySQLShellTool()
	databases, err := mysqlShellTool.ListDatabases(conn)
	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "error", err)
		return nil, err
	}
	log.Info("Listed databases successfully", "host", conn.Host, "count", len(databases))
	return databases, nil
}

// ListTables 查询指定数据库下的表列表
func (a *App) ListTables(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing tables", "host", conn.Host, "database", conn.Database)
	mysqlShellTool := tools.NewMySQLShellTool()
	tables, err := mysqlShellTool.ListTables(conn)
	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "database", conn.Database, "error", err)
		return nil, err
	}
	log.Info("Listed tables successfully", "host", conn.Host, "database", conn.Database, "count", len(tables))
	return tables, nil
}

// SelectFolder 打开系统文件夹选择对话框
func (a *App) SelectFolder() (string, error) {
	return tools.SelectFolder()
}

func (a *App) GetExportSettings() (tools.ExportSettings, error) {
	log.Info("Getting export settings")
	settings, err := tools.GetExportSettings()
	if err != nil {
		log.Error("Error getting export settings", "error", err)
		return tools.ExportSettings{}, err
	}
	log.Info("Got export settings successfully")
	return settings, nil
}

func (a *App) SaveExportSettings(settings tools.ExportSettings) error {
	log.Info("Saving export settings")
	err := tools.SaveExportSettings(settings)
	if err != nil {
		log.Error("Error saving export settings", "error", err)
		return err
	}
	log.Info("Saved export settings successfully")
	return nil
}

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

type ExportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

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

type ImportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

type CheckImportConflictsRequest struct {
	ConnectionID string `json:"connection_id"`
	InputDir     string `json:"input_dir"`
}

type DropConflictingTablesRequest struct {
	ConnectionID string                 `json:"connection_id"`
	Conflicts    []tools.ImportConflict `json:"conflicts"`
}

func (a *App) ExportDatabases(req ExportRequest) (ExportResponse, error) {
	log.Info("Exporting databases", "databases", req.Databases, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	conn, err := tools.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ExportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	mysqlShellTool := tools.NewMySQLShellTool()
	config := tools.ExportConfig{
		OutputDir:      req.OutputDir,
		Threads:        req.Threads,
		Compression:    req.Compression,
		ChunkSize:      req.ChunkSize,
		IncludeSchemas: req.IncludeSchemas,
		ExcludeSchemas: req.ExcludeSchemas,
		IncludeTables:  req.IncludeTables,
		ExcludeTables:  req.ExcludeTables,
		Overwrite:      req.Overwrite,
	}

	err = mysqlShellTool.ExportDatabases(conn, req.Databases, config)
	if err != nil {
		log.Error("Error exporting databases", "databases", req.Databases, "error", err)
		return ExportResponse{Success: false, Message: "导出失败: " + err.Error()}, nil
	}

	log.Info("Exported databases successfully", "databases", req.Databases, "output", req.OutputDir)
	return ExportResponse{
		Success: true,
		Message: "导出成功",
		Path:    req.OutputDir,
	}, nil
}

func (a *App) ExportTables(req ExportRequest) (ExportResponse, error) {
	log.Info("Exporting tables", "database", req.Database, "tables", req.Tables, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.Database == "" {
		return ExportResponse{Success: false, Message: "数据库名不能为空"}, nil
	}

	if len(req.Tables) == 0 {
		return ExportResponse{Success: false, Message: "请选择要导出的表"}, nil
	}

	conn, err := tools.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ExportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	conn.Database = req.Database

	mysqlShellTool := tools.NewMySQLShellTool()
	config := tools.ExportConfig{
		OutputDir:      req.OutputDir,
		Threads:        req.Threads,
		Compression:    req.Compression,
		ChunkSize:      req.ChunkSize,
		IncludeSchemas: req.IncludeSchemas,
		ExcludeSchemas: req.ExcludeSchemas,
		IncludeTables:  req.IncludeTables,
		ExcludeTables:  req.ExcludeTables,
		Overwrite:      req.Overwrite,
	}

	err = mysqlShellTool.ExportTables(conn, req.Database, req.Tables, config)
	if err != nil {
		log.Error("Error exporting tables", "database", req.Database, "tables", req.Tables, "error", err)
		return ExportResponse{Success: false, Message: "导出失败: " + err.Error()}, nil
	}

	log.Info("Exported tables successfully", "database", req.Database, "tables", req.Tables, "output", req.OutputDir)
	return ExportResponse{
		Success: true,
		Message: "导出成功",
		Path:    req.OutputDir,
	}, nil
}

func (a *App) ImportDatabases(req ImportRequest) (ImportResponse, error) {
	log.Info("Importing databases", "input", req.InputDir)

	if req.ConnectionID == "" {
		return ImportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.InputDir == "" {
		return ImportResponse{Success: false, Message: "导入目录不能为空"}, nil
	}

	conn, err := tools.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ImportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	mysqlShellTool := tools.NewMySQLShellTool()
	config := tools.ImportConfig{
		InputDir:       req.InputDir,
		Threads:        req.Threads,
		Schema:         req.Schema,
		IncludeSchemas: req.IncludeSchemas,
		ExcludeSchemas: req.ExcludeSchemas,
		IncludeTables:  req.IncludeTables,
		ExcludeTables:  req.ExcludeTables,
		ResetProgress:  req.ResetProgress,
		WaitTimeout:    req.WaitTimeout,
	}

	err = mysqlShellTool.ImportDatabases(conn, config)
	if err != nil {
		log.Error("Error importing databases", "input", req.InputDir, "error", err)
		return ImportResponse{Success: false, Message: "导入失败: " + err.Error()}, nil
	}

	log.Info("Imported databases successfully", "input", req.InputDir)
	return ImportResponse{
		Success: true,
		Message: "导入成功",
		Path:    req.InputDir,
	}, nil
}

func (a *App) ImportTables(req ImportRequest) (ImportResponse, error) {
	log.Info("Importing tables", "database", req.Database, "input", req.InputDir)

	if req.ConnectionID == "" {
		return ImportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.InputDir == "" {
		return ImportResponse{Success: false, Message: "导入目录不能为空"}, nil
	}

	if req.Database == "" {
		return ImportResponse{Success: false, Message: "目标数据库不能为空"}, nil
	}

	conn, err := tools.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ImportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	mysqlShellTool := tools.NewMySQLShellTool()
	config := tools.ImportConfig{
		InputDir:      req.InputDir,
		Threads:       req.Threads,
		IncludeTables: req.IncludeTables,
		ExcludeTables: req.ExcludeTables,
		ResetProgress: req.ResetProgress,
		WaitTimeout:   req.WaitTimeout,
	}

	err = mysqlShellTool.ImportTables(conn, req.Database, config)
	if err != nil {
		log.Error("Error importing tables", "database", req.Database, "input", req.InputDir, "error", err)
		return ImportResponse{Success: false, Message: "导入失败: " + err.Error()}, nil
	}

	log.Info("Imported tables successfully", "database", req.Database, "input", req.InputDir)
	return ImportResponse{
		Success: true,
		Message: "导入成功",
		Path:    req.InputDir,
	}, nil
}

func (a *App) CheckImportConflicts(req CheckImportConflictsRequest) (tools.ImportConflictCheckResult, error) {
	log.Info("Checking import conflicts", "input", req.InputDir)

	if req.ConnectionID == "" {
		return tools.ImportConflictCheckResult{}, fmt.Errorf("连接ID不能为空")
	}

	if req.InputDir == "" {
		return tools.ImportConflictCheckResult{}, fmt.Errorf("导入目录不能为空")
	}

	conn, err := tools.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return tools.ImportConflictCheckResult{}, fmt.Errorf("获取连接信息失败: %s", err.Error())
	}

	mysqlShellTool := tools.NewMySQLShellTool()

	_, schemaMetas, err := mysqlShellTool.ParseDumpMetadata(req.InputDir)
	if err != nil {
		log.Error("Error parsing dump metadata", "input", req.InputDir, "error", err)
		return tools.ImportConflictCheckResult{}, err
	}

	result, err := mysqlShellTool.CheckExistingObjects(conn, schemaMetas)
	if err != nil {
		log.Error("Error checking existing objects", "input", req.InputDir, "error", err)
		return tools.ImportConflictCheckResult{}, err
	}

	log.Info("Checked import conflicts", "has_conflicts", result.HasConflicts)
	return *result, nil
}

func (a *App) DropConflictingTables(req DropConflictingTablesRequest) error {
	log.Info("Dropping conflicting tables", "conflict_count", len(req.Conflicts))

	if req.ConnectionID == "" {
		return fmt.Errorf("连接ID不能为空")
	}

	if len(req.Conflicts) == 0 {
		return nil
	}

	conn, err := tools.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return fmt.Errorf("获取连接信息失败: %s", err.Error())
	}

	mysqlShellTool := tools.NewMySQLShellTool()
	err = mysqlShellTool.DropObjects(conn, req.Conflicts)
	if err != nil {
		log.Error("Error dropping conflicting objects", "error", err)
		return err
	}

	log.Info("Dropped conflicting tables successfully")
	return nil
}
