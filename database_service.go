package main

import (
	"context"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/services"
	"gotools/internal/tools"
)

// DatabaseService 数据库操作 Wails 服务
type DatabaseService struct {
	ctx    context.Context
	inner *services.DatabaseService
}

// NewDatabaseService 创建数据库服务
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{
		inner: services.NewDatabaseService(),
	}
}

// Startup 应用启动时由 Wails 调用
func (s *DatabaseService) Startup(ctx context.Context) {
	s.ctx = ctx
	log.Info("DatabaseService started")
}

// GetDatabaseConnections 获取所有数据库连接信息
func (s *DatabaseService) GetDatabaseConnections() ([]tools.DatabaseConnection, error) {
	return s.inner.GetDatabaseConnections()
}

// GetDatabaseConnection 根据 ID 获取数据库连接信息
func (s *DatabaseService) GetDatabaseConnection(id string) (tools.DatabaseConnection, error) {
	return s.inner.GetDatabaseConnection(id)
}

// AddDatabaseConnection 添加数据库连接信息
func (s *DatabaseService) AddDatabaseConnection(conn tools.DatabaseConnection) error {
	return s.inner.AddDatabaseConnection(conn)
}

// UpdateDatabaseConnection 更新数据库连接信息
func (s *DatabaseService) UpdateDatabaseConnection(conn tools.DatabaseConnection) error {
	return s.inner.UpdateDatabaseConnection(conn)
}

// DeleteDatabaseConnection 删除数据库连接信息
func (s *DatabaseService) DeleteDatabaseConnection(id string) error {
	return s.inner.DeleteDatabaseConnection(id)
}

// ConnectDatabase 使用 mysqlshell 连接数据库
func (s *DatabaseService) ConnectDatabase(conn tools.DatabaseConnection) error {
	return s.inner.ConnectDatabase(conn)
}

// ListDatabases 查询数据库列表
func (s *DatabaseService) ListDatabases(conn tools.DatabaseConnection) ([]string, error) {
	return s.inner.ListDatabases(conn)
}

// ListTables 查询指定数据库下的表列表
func (s *DatabaseService) ListTables(conn tools.DatabaseConnection) ([]string, error) {
	return s.inner.ListTables(conn)
}

// GetExportSettings 获取导出设置
func (s *DatabaseService) GetExportSettings() (tools.ExportSettings, error) {
	return s.inner.GetExportSettings()
}

// SaveExportSettings 保存导出设置
func (s *DatabaseService) SaveExportSettings(settings tools.ExportSettings) error {
	return s.inner.SaveExportSettings(settings)
}

// ExportDatabases 导出数据库（使用MySQL Shell）
func (s *DatabaseService) ExportDatabases(req ExportRequest) (ExportResponse, error) {
	log.Info("Exporting databases", "databases", req.Databases, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
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

	err = s.inner.ExportDatabases(conn, req.Databases, config)
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
func (s *DatabaseService) ExportTables(req ExportRequest) (ExportResponse, error) {
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

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
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

	err = s.inner.ExportTables(conn, req.Database, req.Tables, config)
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
func (s *DatabaseService) ImportDatabases(req ImportRequest) (ImportResponse, error) {
	log.Info("Importing databases", "input", req.InputDir)

	if req.ConnectionID == "" {
		return ImportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.InputDir == "" {
		return ImportResponse{Success: false, Message: "导入目录不能为空"}, nil
	}

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
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

	err = s.inner.ImportDatabases(conn, config)
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
func (s *DatabaseService) ImportTables(req ImportRequest) (ImportResponse, error) {
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

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
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

	err = s.inner.ImportTables(conn, req.Database, config)
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
func (s *DatabaseService) CheckImportConflicts(req CheckImportConflictsRequest) (tools.ImportConflictCheckResult, error) {
	log.Info("Checking import conflicts", "input", req.InputDir)

	if req.ConnectionID == "" {
		return tools.ImportConflictCheckResult{}, fmt.Errorf("连接ID不能为空")
	}

	if req.InputDir == "" {
		return tools.ImportConflictCheckResult{}, fmt.Errorf("导入目录不能为空")
	}

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return tools.ImportConflictCheckResult{}, fmt.Errorf("获取连接信息失败: %s", err.Error())
	}

	result, err := s.inner.CheckImportConflicts(conn, req.InputDir)
	if err != nil {
		log.Error("Error checking import conflicts", "input", req.InputDir, "error", err)
		return tools.ImportConflictCheckResult{}, err
	}

	log.Info("Checked import conflicts", "has_conflicts", result.HasConflicts)
	return *result, nil
}

// DropConflictingTables 删除冲突表
func (s *DatabaseService) DropConflictingTables(req DropConflictingTablesRequest) error {
	log.Info("Dropping conflicting tables", "conflict_count", len(req.Conflicts))

	if req.ConnectionID == "" {
		return fmt.Errorf("连接ID不能为空")
	}

	if len(req.Conflicts) == 0 {
		return nil
	}

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return fmt.Errorf("获取连接信息失败: %s", err.Error())
	}

	err = s.inner.DropConflictingTables(conn, req.Conflicts)
	if err != nil {
		log.Error("Error dropping conflicting objects", "error", err)
		return err
	}

	log.Info("Dropped conflicting tables successfully")
	return nil
}

// ==================== mysqldump 相关方法 ====================

// ListDatabasesMySQLDump 使用 mysqldump 工具查询数据库列表
func (s *DatabaseService) ListDatabasesMySQLDump(conn tools.DatabaseConnection) ([]string, error) {
	return s.inner.ListDatabasesMySQLDump(conn)
}

// ListTablesMySQLDump 使用 mysqldump 工具查询表列表
func (s *DatabaseService) ListTablesMySQLDump(conn tools.DatabaseConnection) ([]string, error) {
	return s.inner.ListTablesMySQLDump(conn)
}

// ExportDatabasesMySQLDump 使用 mysqldump 导出多个数据库
func (s *DatabaseService) ExportDatabasesMySQLDump(req MySQLDumpExportRequest) (ExportResponse, error) {
	log.Info("Exporting databases using mysqldump", "databases", req.Databases, "output", req.OutputDir)

	if req.ConnectionID == "" {
		return ExportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if len(req.Databases) == 0 {
		return ExportResponse{Success: false, Message: "请选择要导出的数据库"}, nil
	}

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
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
		err = s.inner.ExportDatabaseMySQLDump(conn, req.Databases[0], config)
	} else {
		err = s.inner.ExportDatabasesMySQLDump(conn, req.Databases, config)
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
func (s *DatabaseService) ExportTablesMySQLDump(req MySQLDumpExportRequest) (ExportResponse, error) {
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

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
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

	err = s.inner.ExportTablesMySQLDump(conn, req.Database, req.Tables, config)
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
func (s *DatabaseService) ImportDumpMySQLDump(req MySQLDumpImportRequest) (ImportResponse, error) {
	log.Info("Importing dump file using mysql", "input", req.InputFile, "database", req.Database)

	if req.ConnectionID == "" {
		return ImportResponse{Success: false, Message: "连接ID不能为空"}, nil
	}

	if req.InputFile == "" {
		return ImportResponse{Success: false, Message: "输入文件不能为空"}, nil
	}

	conn, err := s.inner.GetDatabaseConnection(req.ConnectionID)
	if err != nil {
		log.Error("Error getting database connection", "id", req.ConnectionID, "error", err)
		return ImportResponse{Success: false, Message: "获取连接信息失败: " + err.Error()}, nil
	}

	err = s.inner.ImportDumpMySQLDump(conn, req.InputFile, req.Database)
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
