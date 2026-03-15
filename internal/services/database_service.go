package services

import (
	"gotools/internal/log"
	"gotools/internal/tools"
)

// DatabaseService 数据库操作服务
type DatabaseService struct {
	mysqlShellTool *tools.MySQLShellTool
	mysqlDumpTool  *tools.MySQLDumpTool
}

// NewDatabaseService 创建数据库服务实例
func NewDatabaseService() *DatabaseService {
	log.Info("Creating new DatabaseService instance")
	return &DatabaseService{
		mysqlShellTool: tools.NewMySQLShellTool(),
		mysqlDumpTool:  tools.NewMySQLDumpTool(),
	}
}

// GetDatabaseConnections 获取所有数据库连接信息
func (s *DatabaseService) GetDatabaseConnections() ([]tools.DatabaseConnection, error) {
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
func (s *DatabaseService) GetDatabaseConnection(id string) (tools.DatabaseConnection, error) {
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
func (s *DatabaseService) AddDatabaseConnection(conn tools.DatabaseConnection) error {
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
func (s *DatabaseService) UpdateDatabaseConnection(conn tools.DatabaseConnection) error {
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
func (s *DatabaseService) DeleteDatabaseConnection(id string) error {
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
func (s *DatabaseService) ConnectDatabase(conn tools.DatabaseConnection) error {
	log.Info("Connecting to database", "name", conn.Name, "host", conn.Host)
	err := s.mysqlShellTool.ConnectDatabase(conn)
	if err != nil {
		log.Error("Error connecting to database", "name", conn.Name, "host", conn.Host, "error", err)
		return err
	}
	log.Info("Connected to database successfully", "name", conn.Name, "host", conn.Host)
	return nil
}

// ListDatabases 查询数据库列表
func (s *DatabaseService) ListDatabases(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing databases", "host", conn.Host)
	databases, err := s.mysqlShellTool.ListDatabases(conn)
	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "error", err)
		return nil, err
	}
	log.Info("Listed databases successfully", "host", conn.Host, "count", len(databases))
	return databases, nil
}

// ListTables 查询指定数据库下的表列表
func (s *DatabaseService) ListTables(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing tables", "host", conn.Host, "database", conn.Database)
	tables, err := s.mysqlShellTool.ListTables(conn)
	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "database", conn.Database, "error", err)
		return nil, err
	}
	log.Info("Listed tables successfully", "host", conn.Host, "database", conn.Database, "count", len(tables))
	return tables, nil
}

// GetExportSettings 获取导出设置
func (s *DatabaseService) GetExportSettings() (tools.ExportSettings, error) {
	log.Info("Getting export settings")
	settings, err := tools.GetExportSettings()
	if err != nil {
		log.Error("Error getting export settings", "error", err)
		return tools.ExportSettings{}, err
	}
	log.Info("Got export settings successfully")
	return settings, nil
}

// SaveExportSettings 保存导出设置
func (s *DatabaseService) SaveExportSettings(settings tools.ExportSettings) error {
	log.Info("Saving export settings")
	err := tools.SaveExportSettings(settings)
	if err != nil {
		log.Error("Error saving export settings", "error", err)
		return err
	}
	log.Info("Saved export settings successfully")
	return nil
}

// ExportDatabase 导出单个数据库
func (s *DatabaseService) ExportDatabase(conn tools.DatabaseConnection, database string, config tools.ExportConfig) error {
	log.Info("Exporting database", "database", database, "host", conn.Host)
	err := s.mysqlShellTool.ExportDatabase(conn, database, config)
	if err != nil {
		log.Error("Error exporting database", "database", database, "error", err)
		return err
	}
	log.Info("Exported database successfully", "database", database)
	return nil
}

// ExportDatabases 导出多个数据库
func (s *DatabaseService) ExportDatabases(conn tools.DatabaseConnection, databases []string, config tools.ExportConfig) error {
	log.Info("Exporting databases", "databases", databases, "host", conn.Host)
	err := s.mysqlShellTool.ExportDatabases(conn, databases, config)
	if err != nil {
		log.Error("Error exporting databases", "databases", databases, "error", err)
		return err
	}
	log.Info("Exported databases successfully", "databases", databases)
	return nil
}

// ExportTables 导出指定表
func (s *DatabaseService) ExportTables(conn tools.DatabaseConnection, database string, tables []string, config tools.ExportConfig) error {
	log.Info("Exporting tables", "database", database, "tables", tables, "host", conn.Host)
	err := s.mysqlShellTool.ExportTables(conn, database, tables, config)
	if err != nil {
		log.Error("Error exporting tables", "database", database, "tables", tables, "error", err)
		return err
	}
	log.Info("Exported tables successfully", "database", database, "tables", tables)
	return nil
}

// ImportDatabases 导入数据库
func (s *DatabaseService) ImportDatabases(conn tools.DatabaseConnection, config tools.ImportConfig) error {
	log.Info("Importing databases", "input", config.InputDir, "host", conn.Host)
	err := s.mysqlShellTool.ImportDatabases(conn, config)
	if err != nil {
		log.Error("Error importing databases", "input", config.InputDir, "error", err)
		return err
	}
	log.Info("Imported databases successfully", "input", config.InputDir)
	return nil
}

// ImportTables 导入表
func (s *DatabaseService) ImportTables(conn tools.DatabaseConnection, database string, config tools.ImportConfig) error {
	log.Info("Importing tables", "database", database, "input", config.InputDir, "host", conn.Host)
	err := s.mysqlShellTool.ImportTables(conn, database, config)
	if err != nil {
		log.Error("Error importing tables", "database", database, "input", config.InputDir, "error", err)
		return err
	}
	log.Info("Imported tables successfully", "database", database, "input", config.InputDir)
	return nil
}

// CheckImportConflicts 检查导入冲突
func (s *DatabaseService) CheckImportConflicts(conn tools.DatabaseConnection, inputDir string) (*tools.ImportConflictCheckResult, error) {
	log.Info("Checking import conflicts", "input", inputDir, "host", conn.Host)

	dumpMeta, schemaMetas, err := s.mysqlShellTool.ParseDumpMetadata(inputDir)
	if err != nil {
		log.Error("Error parsing dump metadata", "input", inputDir, "error", err)
		return nil, err
	}

	result, err := s.mysqlShellTool.CheckExistingObjects(conn, schemaMetas)
	if err != nil {
		log.Error("Error checking existing objects", "input", inputDir, "error", err)
		return nil, err
	}

	log.Info("Checked import conflicts", "has_conflicts", result.HasConflicts, "schemas", dumpMeta.Schemas)
	return result, nil
}

// DropConflictingTables 删除冲突对象
func (s *DatabaseService) DropConflictingTables(conn tools.DatabaseConnection, conflicts []tools.ImportConflict) error {
	log.Info("Dropping conflicting tables", "conflict_count", len(conflicts))
	err := s.mysqlShellTool.DropObjects(conn, conflicts)
	if err != nil {
		log.Error("Error dropping conflicting objects", "error", err)
		return err
	}
	log.Info("Dropped conflicting tables successfully")
	return nil
}

// ListDatabasesMySQLDump 使用 mysqldump 查询数据库列表
func (s *DatabaseService) ListDatabasesMySQLDump(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing databases using mysqldump", "host", conn.Host)
	databases, err := s.mysqlDumpTool.ListDatabases(conn)
	if err != nil {
		log.Error("Error listing databases using mysqldump", "host", conn.Host, "error", err)
		return nil, err
	}
	log.Info("Listed databases successfully using mysqldump", "host", conn.Host, "count", len(databases))
	return databases, nil
}

// ListTablesMySQLDump 使用 mysqldump 查询表列表
func (s *DatabaseService) ListTablesMySQLDump(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing tables using mysqldump", "host", conn.Host, "database", conn.Database)
	tables, err := s.mysqlDumpTool.ListTables(conn)
	if err != nil {
		log.Error("Error listing tables using mysqldump", "host", conn.Host, "database", conn.Database, "error", err)
		return nil, err
	}
	log.Info("Listed tables successfully using mysqldump", "host", conn.Host, "database", conn.Database, "count", len(tables))
	return tables, nil
}

// ExportDatabaseMySQLDump 使用 mysqldump 导出数据库
func (s *DatabaseService) ExportDatabaseMySQLDump(conn tools.DatabaseConnection, database string, config tools.MySQLDumpConfig) error {
	log.Info("Exporting database using mysqldump", "database", database, "host", conn.Host)
	err := s.mysqlDumpTool.ExportDatabase(conn, database, config)
	if err != nil {
		log.Error("Error exporting database using mysqldump", "database", database, "error", err)
		return err
	}
	log.Info("Exported database successfully using mysqldump", "database", database)
	return nil
}

// ExportDatabasesMySQLDump 使用 mysqldump 导出多个数据库
func (s *DatabaseService) ExportDatabasesMySQLDump(conn tools.DatabaseConnection, databases []string, config tools.MySQLDumpConfig) error {
	log.Info("Exporting databases using mysqldump", "databases", databases, "host", conn.Host)
	err := s.mysqlDumpTool.ExportDatabases(conn, databases, config)
	if err != nil {
		log.Error("Error exporting databases using mysqldump", "databases", databases, "error", err)
		return err
	}
	log.Info("Exported databases successfully using mysqldump", "databases", databases)
	return nil
}

// ExportTablesMySQLDump 使用 mysqldump 导出表
func (s *DatabaseService) ExportTablesMySQLDump(conn tools.DatabaseConnection, database string, tables []string, config tools.MySQLDumpConfig) error {
	log.Info("Exporting tables using mysqldump", "database", database, "tables", tables, "host", conn.Host)
	err := s.mysqlDumpTool.ExportTables(conn, database, tables, config)
	if err != nil {
		log.Error("Error exporting tables using mysqldump", "database", database, "tables", tables, "error", err)
		return err
	}
	log.Info("Exported tables successfully using mysqldump", "database", database, "tables", tables)
	return nil
}

// ImportDumpMySQLDump 使用 mysql 导入 SQL 文件
func (s *DatabaseService) ImportDumpMySQLDump(conn tools.DatabaseConnection, inputFile string, database string) error {
	log.Info("Importing dump file using mysql", "input", inputFile, "database", database)
	config := tools.MySQLImportConfig{
		InputFile: inputFile,
		Database:  database,
	}
	err := s.mysqlDumpTool.ImportDump(conn, config)
	if err != nil {
		log.Error("Error importing dump file", "input", inputFile, "error", err)
		return err
	}
	log.Info("Imported dump file successfully", "input", inputFile, "database", database)
	return nil
}
