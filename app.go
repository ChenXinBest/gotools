package main

import (
	"context"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/services"
	"gotools/internal/tools"
	"gotools/internal/version"
)

// App 应用主结构，负责Wails绑定
type App struct {
	ctx             context.Context
	processService  *services.ProcessService
	databaseService *services.DatabaseService
	dialogService   *services.DialogService
}

// NewApp 创建应用实例
func NewApp() *App {
	log.Info("Creating new app instance")
	return &App{
		processService:  services.NewProcessService(),
		databaseService: services.NewDatabaseService(),
		dialogService:   services.NewDialogService(),
	}
}

// startup 应用启动回调
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Info("App startup")
}

// Greet 问候用户（示例方法）
func (a *App) Greet(name string) string {
	log.Info("Greeting user", "name", name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetSystemProcessInfos 获取系统进程信息列表
func (a *App) GetSystemProcessInfos() ([]tools.ProcessInfo, error) {
	return a.processService.GetSystemProcessInfos()
}

// SearchPidByKeyWord 搜索进程
func (a *App) SearchPidByKeyWord(keyword string) (tools.ProcessInfo, error) {
	return a.processService.SearchPidByKeyWord(keyword)
}

// KillProcessByPID 终止进程
func (a *App) KillProcessByPID(pid int32) error {
	return a.processService.KillProcessByPID(pid)
}

// GetDatabaseConnections 获取所有数据库连接信息
func (a *App) GetDatabaseConnections() ([]tools.DatabaseConnection, error) {
	return a.databaseService.GetDatabaseConnections()
}

// GetDatabaseConnection 根据 ID 获取数据库连接信息
func (a *App) GetDatabaseConnection(id string) (tools.DatabaseConnection, error) {
	return a.databaseService.GetDatabaseConnection(id)
}

// AddDatabaseConnection 添加数据库连接信息
func (a *App) AddDatabaseConnection(conn tools.DatabaseConnection) error {
	return a.databaseService.AddDatabaseConnection(conn)
}

// UpdateDatabaseConnection 更新数据库连接信息
func (a *App) UpdateDatabaseConnection(conn tools.DatabaseConnection) error {
	return a.databaseService.UpdateDatabaseConnection(conn)
}

// DeleteDatabaseConnection 删除数据库连接信息
func (a *App) DeleteDatabaseConnection(id string) error {
	return a.databaseService.DeleteDatabaseConnection(id)
}

// ConnectDatabase 使用 mysqlshell 连接数据库
func (a *App) ConnectDatabase(conn tools.DatabaseConnection) error {
	return a.databaseService.ConnectDatabase(conn)
}

// ListDatabases 查询数据库列表
func (a *App) ListDatabases(conn tools.DatabaseConnection) ([]string, error) {
	return a.databaseService.ListDatabases(conn)
}

// ListTables 查询指定数据库下的表列表
func (a *App) ListTables(conn tools.DatabaseConnection) ([]string, error) {
	return a.databaseService.ListTables(conn)
}

// SelectFolder 打开系统文件夹选择对话框
func (a *App) SelectFolder() (string, error) {
	return a.dialogService.SelectFolder()
}

// SelectFile 打开系统文件选择对话框
func (a *App) SelectFile() (string, error) {
	return a.dialogService.SelectFile()
}

// SelectSaveFile 打开系统保存文件对话框
func (a *App) SelectSaveFile(defaultName string) (string, error) {
	return a.dialogService.SelectSaveFile(defaultName)
}

// GetExportSettings 获取导出设置
func (a *App) GetExportSettings() (tools.ExportSettings, error) {
	return a.databaseService.GetExportSettings()
}

// SaveExportSettings 保存导出设置
func (a *App) SaveExportSettings(settings tools.ExportSettings) error {
	return a.databaseService.SaveExportSettings(settings)
}

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

// ExportDatabases 导出数据库（使用MySQL Shell）
func (a *App) ExportDatabases(req ExportRequest) (ExportResponse, error) {
	log.Info("Exporting databases", "databases", req.Databases, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ExportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

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

	err = a.databaseService.ExportDatabases(conn, req.Databases, config)
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

// ExportTables 导出表（使用MySQL Shell）
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

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ExportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	conn.Database = req.Database

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

	err = a.databaseService.ExportTables(conn, req.Database, req.Tables, config)
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

// ImportDatabases 导入数据库（使用MySQL Shell）
func (a *App) ImportDatabases(req ImportRequest) (ImportResponse, error) {
	log.Info("Importing databases", "input", req.InputDir)

	if req.ConnectionID == "" {
		return ImportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.InputDir == "" {
		return ImportResponse{Success: false, Message: "导入目录不能为空"}, nil
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ImportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

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

	err = a.databaseService.ImportDatabases(conn, config)
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

// ImportTables 导入表（使用MySQL Shell）
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

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ImportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	config := tools.ImportConfig{
		InputDir:      req.InputDir,
		Threads:       req.Threads,
		IncludeTables: req.IncludeTables,
		ExcludeTables: req.ExcludeTables,
		ResetProgress: req.ResetProgress,
		WaitTimeout:   req.WaitTimeout,
	}

	err = a.databaseService.ImportTables(conn, req.Database, config)
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

// CheckImportConflicts 检查导入冲突
func (a *App) CheckImportConflicts(req CheckImportConflictsRequest) (tools.ImportConflictCheckResult, error) {
	log.Info("Checking import conflicts", "input", req.InputDir)

	if req.ConnectionID == "" {
		return tools.ImportConflictCheckResult{}, fmt.Errorf("连接ID不能为空")
	}

	if req.InputDir == "" {
		return tools.ImportConflictCheckResult{}, fmt.Errorf("导入目录不能为空")
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return tools.ImportConflictCheckResult{}, fmt.Errorf("获取连接信息失败: %s", err.Error())
	}

	result, err := a.databaseService.CheckImportConflicts(conn, req.InputDir)
	if err != nil {
		log.Error("Error checking import conflicts", "input", req.InputDir, "error", err)
		return tools.ImportConflictCheckResult{}, err
	}

	log.Info("Checked import conflicts", "has_conflicts", result.HasConflicts)
	return *result, nil
}

// DropConflictingTables 删除冲突表
func (a *App) DropConflictingTables(req DropConflictingTablesRequest) error {
	log.Info("Dropping conflicting tables", "conflict_count", len(req.Conflicts))

	if req.ConnectionID == "" {
		return fmt.Errorf("连接ID不能为空")
	}

	if len(req.Conflicts) == 0 {
		return nil
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return fmt.Errorf("获取连接信息失败: %s", err.Error())
	}

	err = a.databaseService.DropConflictingTables(conn, req.Conflicts)
	if err != nil {
		log.Error("Error dropping conflicting objects", "error", err)
		return err
	}

	log.Info("Dropped conflicting tables successfully")
	return nil
}

// ==================== mysqldump 相关方法 ====================

// ListDatabasesMySQLDump 使用 mysqldump 工具查询数据库列表
func (a *App) ListDatabasesMySQLDump(conn tools.DatabaseConnection) ([]string, error) {
	return a.databaseService.ListDatabasesMySQLDump(conn)
}

// ListTablesMySQLDump 使用 mysqldump 工具查询表列表
func (a *App) ListTablesMySQLDump(conn tools.DatabaseConnection) ([]string, error) {
	return a.databaseService.ListTablesMySQLDump(conn)
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

// ExportDatabasesMySQLDump 使用 mysqldump 导出多个数据库
func (a *App) ExportDatabasesMySQLDump(req MySQLDumpExportRequest) (ExportResponse, error) {
	log.Info("Exporting databases using mysqldump", "databases", req.Databases, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if len(req.Databases) == 0 {
		return ExportResponse{Success: false, Message: "请选择要导出的数据库"}, nil
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ExportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	config := tools.MySQLDumpConfig{
		OutputDir:         req.OutputDir,
		Compression:       req.Compression,
		SingleTransaction: req.SingleTransaction,
		Routines:          req.Routines,
		Events:            req.Events,
		Overwrite:         req.Overwrite,
	}

	if len(req.Databases) == 1 {
		err = a.databaseService.ExportDatabaseMySQLDump(conn, req.Databases[0], config)
	} else {
		err = a.databaseService.ExportDatabasesMySQLDump(conn, req.Databases, config)
	}

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

// ExportTablesMySQLDump 使用 mysqldump 导出指定表
func (a *App) ExportTablesMySQLDump(req MySQLDumpExportRequest) (ExportResponse, error) {
	log.Info("Exporting tables using mysqldump", "database", req.Database, "tables", req.Tables, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.Database == "" {
		return ExportResponse{Success: false, Message: "数据库名不能为空"}, nil
	}

	if len(req.Tables) == 0 {
		return ExportResponse{Success: false, Message: "请选择要导出的表"}, nil
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ExportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	conn.Database = req.Database

	config := tools.MySQLDumpConfig{
		OutputDir:         req.OutputDir,
		Compression:       req.Compression,
		SingleTransaction: req.SingleTransaction,
		Routines:          req.Routines,
		Events:            req.Events,
		Overwrite:         req.Overwrite,
	}

	err = a.databaseService.ExportTablesMySQLDump(conn, req.Database, req.Tables, config)
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

// ImportDumpMySQLDump 使用 mysql 导入 SQL 文件
func (a *App) ImportDumpMySQLDump(req MySQLDumpImportRequest) (ImportResponse, error) {
	log.Info("Importing dump file using mysql", "input", req.InputFile, "database", req.Database)

	if req.ConnectionID == "" {
		return ImportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.InputFile == "" {
		return ImportResponse{Success: false, Message: "输入文件不能为空"}, nil
	}

	conn, err := a.databaseService.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ImportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	err = a.databaseService.ImportDumpMySQLDump(conn, req.InputFile, req.Database)
	if err != nil {
		log.Error("Error importing dump file", "input", req.InputFile, "error", err)
		return ImportResponse{Success: false, Message: "导入失败: " + err.Error()}, nil
	}

	log.Info("Imported dump file successfully", "input", req.InputFile, "database", req.Database)
	return ImportResponse{
		Success: true,
		Message: "导入成功",
		Path:    req.InputFile,
	}, nil
}

// GetVersion 获取应用版本信息
func (a *App) GetVersion() version.Info {
	log.Info("Getting version info")
	return version.GetInfo()
}

// GetVersionString 获取版本字符串
func (a *App) GetVersionString() string {
	return version.GetVersionString()
}
