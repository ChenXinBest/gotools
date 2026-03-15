package tools

import (
	"fmt"
	"gotools/internal/log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ErrMySQLDumpNotFound = "MYSQLDUMP_NOT_FOUND"
	ErrMySQLNotFound     = "MYSQL_NOT_FOUND"
	MySQLDownloadURL     = "https://dev.mysql.com/downloads/"
)

type MySQLDumpTool struct {
	common *CommonTool
}

type MySQLDumpConfig struct {
	OutputDir         string
	Compression       string
	SingleTransaction bool
	Routines          bool
	Events            bool
	Overwrite         bool
}

type MySQLImportConfig struct {
	InputFile string
	Database  string
}

func NewMySQLDumpTool() *MySQLDumpTool {
	log.Info("Creating new MySQLDumpTool instance")
	return &MySQLDumpTool{
		common: NewCommonTool(),
	}
}

func (m *MySQLDumpTool) checkMySQLDumpExists() error {
	err := m.common.CheckCommandExists("mysqldump")
	if err != nil {
		return fmt.Errorf("%s: 未找到 mysqldump 命令，请先安装 MySQL 客户端工具", ErrMySQLDumpNotFound)
	}
	return nil
}

func (m *MySQLDumpTool) checkMySQLExists() error {
	err := m.common.CheckCommandExists("mysql")
	if err != nil {
		return fmt.Errorf("%s: 未找到 mysql 命令，请先安装 MySQL 客户端工具", ErrMySQLNotFound)
	}
	return nil
}

func (m *MySQLDumpTool) ConnectDatabase(conn DatabaseConnection) error {
	if err := m.checkMySQLExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	log.Info("Connecting to database using mysql client", "host", conn.Host, "port", conn.Port, "database", conn.Database)

	args := []string{
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.User,
		fmt.Sprintf("-p%s", conn.Password),
		"-e", "SELECT 1;",
	}

	if conn.Database != "" {
		args = append(args, conn.Database)
	}

	output, err := m.common.ExecuteCommand("mysql", args...)
	if err != nil {
		log.Error("Error connecting to database", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", output)
		return fmt.Errorf("连接失败: %s, 错误信息: %s", err.Error(), output)
	}

	log.Info("Connected to database successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	return nil
}

func (m *MySQLDumpTool) ListDatabases(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLExists(); err != nil {
		return nil, err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return nil, err
	}

	log.Info("Listing databases using mysql client", "host", conn.Host, "port", conn.Port)

	args := []string{
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.User,
		fmt.Sprintf("-p%s", conn.Password),
		"-e", "SHOW DATABASES;",
	}

	output, err := m.common.ExecuteCommand("mysql", args...)
	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "port", conn.Port, "error", err, "output", output)
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), output)
	}

	databases := m.parseListOutput(output)
	log.Info("Listed databases successfully", "host", conn.Host, "port", conn.Port, "count", len(databases))
	return databases, nil
}

func (m *MySQLDumpTool) ListTables(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLExists(); err != nil {
		return nil, err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return nil, err
	}

	if conn.Database == "" {
		return nil, fmt.Errorf("数据库名称不能为空")
	}

	log.Info("Listing tables using mysql client", "host", conn.Host, "port", conn.Port, "database", conn.Database)

	args := []string{
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.User,
		fmt.Sprintf("-p%s", conn.Password),
		"-e", "SHOW TABLES;",
		conn.Database,
	}

	output, err := m.common.ExecuteCommand("mysql", args...)
	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", output)
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), output)
	}

	tables := m.parseListOutput(output)
	log.Info("Listed tables successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database, "count", len(tables))
	return tables, nil
}

func (m *MySQLDumpTool) parseListOutput(output string) []string {
	var results []string
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		if i == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Database") && !strings.HasPrefix(line, "Tables_in_") {
			results = append(results, line)
		}
	}

	return results
}

func (m *MySQLDumpTool) buildOutputPath(config MySQLDumpConfig, filename string) string {
	ext := ".sql"
	if strings.ToLower(config.Compression) == "gzip" {
		ext = ".sql.gz"
	}
	return filepath.Join(config.OutputDir, filename+ext)
}

func (m *MySQLDumpTool) ExportDatabase(conn DatabaseConnection, database string, config MySQLDumpConfig) error {
	if err := m.checkMySQLDumpExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	if err := m.common.ValidateExportConfig(config.OutputDir); err != nil {
		return err
	}

	if database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if err := m.common.EnsureDirExists(config.OutputDir); err != nil {
		return err
	}

	outputFile := m.buildOutputPath(config, database)

	if config.Overwrite {
		if err := m.common.RemoveIfExists(outputFile); err != nil {
			return err
		}
	}

	log.Info("Exporting database", "database", database, "output", outputFile)

	args := m.buildExportArgs(conn, config)
	args = append(args, "--databases", database)

	return m.executeExport(args, outputFile, config.Compression)
}

func (m *MySQLDumpTool) ExportDatabases(conn DatabaseConnection, databases []string, config MySQLDumpConfig) error {
	if err := m.checkMySQLDumpExists(); err != nil {
		return err
	}

	if len(databases) == 0 {
		return fmt.Errorf("数据库列表不能为空")
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	filename := "databases"
	if len(databases) == 1 {
		filename = databases[0]
	}
	outputFile := m.buildOutputPath(config, filename)

	if config.Overwrite {
		if _, err := os.Stat(outputFile); err == nil {
			log.Info("Removing existing output file for overwrite", "path", outputFile)
			if err := os.Remove(outputFile); err != nil {
				log.Error("Failed to remove existing file", "path", outputFile, "error", err)
				return fmt.Errorf("删除已存在文件失败: %s", err.Error())
			}
		}
	}

	log.Info("Exporting databases", "databases", databases, "output", outputFile)

	args := m.buildExportArgs(conn, config)
	args = append(args, "--databases")
	args = append(args, databases...)

	return m.executeExport(args, outputFile, config.Compression)
}

func (m *MySQLDumpTool) ExportTables(conn DatabaseConnection, database string, tables []string, config MySQLDumpConfig) error {
	if err := m.checkMySQLDumpExists(); err != nil {
		return err
	}

	if database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if len(tables) == 0 {
		return fmt.Errorf("表列表不能为空")
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	filename := database
	if len(tables) == 1 {
		filename = fmt.Sprintf("%s_%s", database, tables[0])
	} else {
		filename = fmt.Sprintf("%s_tables", database)
	}
	outputFile := m.buildOutputPath(config, filename)

	if config.Overwrite {
		if _, err := os.Stat(outputFile); err == nil {
			log.Info("Removing existing output file for overwrite", "path", outputFile)
			if err := os.Remove(outputFile); err != nil {
				log.Error("Failed to remove existing file", "path", outputFile, "error", err)
				return fmt.Errorf("删除已存在文件失败: %s", err.Error())
			}
		}
	}

	log.Info("Exporting tables", "database", database, "tables", tables, "output", outputFile)

	args := m.buildExportArgs(conn, config)
	args = append(args, database)
	args = append(args, tables...)

	return m.executeExport(args, outputFile, config.Compression)
}

func (m *MySQLDumpTool) ExportInstance(conn DatabaseConnection, config MySQLDumpConfig) error {
	if err := m.checkMySQLDumpExists(); err != nil {
		return err
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	outputFile := m.buildOutputPath(config, "all_databases")

	if config.Overwrite {
		if _, err := os.Stat(outputFile); err == nil {
			log.Info("Removing existing output file for overwrite", "path", outputFile)
			if err := os.Remove(outputFile); err != nil {
				log.Error("Failed to remove existing file", "path", outputFile, "error", err)
				return fmt.Errorf("删除已存在文件失败: %s", err.Error())
			}
		}
	}

	log.Info("Exporting entire instance", "host", conn.Host, "port", conn.Port, "output", outputFile)

	args := m.buildExportArgs(conn, config)
	args = append(args, "--all-databases")

	return m.executeExport(args, outputFile, config.Compression)
}

func (m *MySQLDumpTool) executeExport(args []string, outputFile string, compression string) error {
	return m.common.ExecuteCommandWithOutputFile("mysqldump", outputFile, compression, args...)
}

func (m *MySQLDumpTool) buildExportArgs(conn DatabaseConnection, config MySQLDumpConfig) []string {
	args := []string{
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.User,
		fmt.Sprintf("-p%s", conn.Password),
	}

	if config.SingleTransaction {
		args = append(args, "--single-transaction")
	}

	if config.Routines {
		args = append(args, "--routines")
	}

	if config.Events {
		args = append(args, "--events")
	}

	args = append(args, "--triggers")

	return args
}

func (m *MySQLDumpTool) ImportDump(conn DatabaseConnection, config MySQLImportConfig) error {
	if err := m.checkMySQLExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	if config.InputFile == "" {
		return fmt.Errorf("输入文件不能为空")
	}

	if !m.common.FileExists(config.InputFile) {
		return fmt.Errorf("输入文件不存在: %s", config.InputFile)
	}

	log.Info("Importing dump file", "input", config.InputFile, "database", config.Database)

	args := []string{
		"-h", conn.Host,
		"-P", fmt.Sprintf("%d", conn.Port),
		"-u", conn.User,
		fmt.Sprintf("-p%s", conn.Password),
	}

	if config.Database != "" {
		args = append(args, config.Database)
	}

	_, err := m.common.ExecuteCommandWithInputFile("mysql", config.InputFile, args...)
	if err != nil {
		log.Error("Error importing dump file", "input", config.InputFile, "error", err)
		return fmt.Errorf("导入失败: %s", err.Error())
	}

	log.Info("Imported dump file successfully", "input", config.InputFile, "database", config.Database)
	return nil
}
