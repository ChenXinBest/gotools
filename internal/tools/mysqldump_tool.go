package tools

import (
	"compress/gzip"
	"fmt"
	"gotools/internal/log"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	ErrMySQLDumpNotFound = "MYSQLDUMP_NOT_FOUND"
	ErrMySQLNotFound     = "MYSQL_NOT_FOUND"
	MySQLDownloadURL     = "https://dev.mysql.com/downloads/"
)

type MySQLDumpTool struct{}

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
	return &MySQLDumpTool{}
}

func (m *MySQLDumpTool) checkMySQLDumpExists() error {
	_, err := exec.LookPath("mysqldump")
	if err != nil {
		log.Error("mysqldump command not found", "error", err)
		return fmt.Errorf("%s: 未找到 mysqldump 命令，请先安装 MySQL 客户端工具", ErrMySQLDumpNotFound)
	}
	return nil
}

func (m *MySQLDumpTool) checkMySQLExists() error {
	_, err := exec.LookPath("mysql")
	if err != nil {
		log.Error("mysql command not found", "error", err)
		return fmt.Errorf("%s: 未找到 mysql 命令，请先安装 MySQL 客户端工具", ErrMySQLNotFound)
	}
	return nil
}

func (m *MySQLDumpTool) ConnectDatabase(conn DatabaseConnection) error {
	if err := m.checkMySQLExists(); err != nil {
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

	cmd := exec.Command("mysql", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error connecting to database", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", string(output))
		return fmt.Errorf("连接失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Connected to database successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	return nil
}

func (m *MySQLDumpTool) ListDatabases(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLExists(); err != nil {
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

	cmd := exec.Command("mysql", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "port", conn.Port, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	databases := m.parseListOutput(string(output))
	log.Info("Listed databases successfully", "host", conn.Host, "port", conn.Port, "count", len(databases))
	return databases, nil
}

func (m *MySQLDumpTool) ListTables(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLExists(); err != nil {
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

	cmd := exec.Command("mysql", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	tables := m.parseListOutput(string(output))
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

	if database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	outputFile := m.buildOutputPath(config, database)

	if config.Overwrite {
		if _, err := os.Stat(outputFile); err == nil {
			log.Info("Removing existing output file for overwrite", "path", outputFile)
			if err := os.Remove(outputFile); err != nil {
				log.Error("Failed to remove existing file", "path", outputFile, "error", err)
				return fmt.Errorf("删除已存在文件失败: %s", err.Error())
			}
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
	cmd := exec.Command("mysqldump", args...)

	var outFile *os.File
	var gzipWriter *gzip.Writer
	var err error

	if strings.ToLower(compression) == "gzip" {
		outFile, err = os.Create(outputFile)
		if err != nil {
			log.Error("Failed to create output file", "path", outputFile, "error", err)
			return fmt.Errorf("创建输出文件失败: %s", err.Error())
		}
		defer outFile.Close()

		gzipWriter = gzip.NewWriter(outFile)
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
		log.Error("Failed to start mysqldump", "error", err)
		return fmt.Errorf("启动 mysqldump 失败: %s", err.Error())
	}

	stderrOutput, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		log.Error("Error executing export", "error", err, "stderr", string(stderrOutput))
		return fmt.Errorf("导出失败: %s, 错误信息: %s", err.Error(), string(stderrOutput))
	}

	log.Info("Export completed successfully", "output", outputFile)
	return nil
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

	if config.InputFile == "" {
		return fmt.Errorf("输入文件不能为空")
	}

	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
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

	cmd := exec.Command("mysql", args...)

	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		log.Error("Failed to open input file", "path", config.InputFile, "error", err)
		return fmt.Errorf("打开输入文件失败: %s", err.Error())
	}
	defer inputFile.Close()

	isGzip := strings.HasSuffix(strings.ToLower(config.InputFile), ".gz")

	if isGzip {
		gzipReader, err := gzip.NewReader(inputFile)
		if err != nil {
			log.Error("Failed to create gzip reader", "path", config.InputFile, "error", err)
			return fmt.Errorf("创建解压读取器失败: %s", err.Error())
		}
		defer gzipReader.Close()

		cmd.Stdin = gzipReader
	} else {
		cmd.Stdin = inputFile
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Error importing dump file", "input", config.InputFile, "error", err, "output", string(output))
		return fmt.Errorf("导入失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Imported dump file successfully", "input", config.InputFile, "database", config.Database)
	return nil
}
