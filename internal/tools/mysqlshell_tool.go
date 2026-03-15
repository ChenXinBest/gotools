package tools

import (
	"encoding/json"
	"fmt"
	"gotools/internal/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	ErrMySQLShellNotFound = "MYSQLSHELL_NOT_FOUND"
	MySQLShellDownloadURL = "https://dev.mysql.com/downloads/shell/"
)

type MySQLShellTool struct{}

type ExportConfig struct {
	OutputDir      string
	Threads        int
	ChunkSize      string
	Compression    string
	IncludeSchemas []string
	ExcludeSchemas []string
	IncludeTables  []string
	ExcludeTables  []string
	Overwrite      bool
}

func NewMySQLShellTool() *MySQLShellTool {
	log.Info("Creating new MySQLShellTool instance")
	return &MySQLShellTool{}
}

func (m *MySQLShellTool) checkMySQLShellExists() error {
	_, err := exec.LookPath("mysqlsh")
	if err != nil {
		log.Error("mysqlsh command not found", "error", err)
		return fmt.Errorf("%s: 未找到 mysqlsh 命令，请先安装 MySQL Shell", ErrMySQLShellNotFound)
	}
	return nil
}

func (m *MySQLShellTool) ConnectDatabase(conn DatabaseConnection) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}
	log.Info("Connecting to database using mysqlshell", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	cmd := exec.Command("mysqlsh",
		"--uri", fmt.Sprintf("%s:%s@%s:%d/%s", conn.User, conn.Password, conn.Host, conn.Port, conn.Database),
		"--execute", "SELECT 1;")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error connecting to database", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", string(output))
		return fmt.Errorf("连接失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Connected to database successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	return nil
}

func (m *MySQLShellTool) ListDatabases(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLShellExists(); err != nil {
		return nil, err
	}
	log.Info("Listing databases using mysqlshell", "host", conn.Host, "port", conn.Port)
	cmd := exec.Command("mysqlsh",
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--execute", "SHOW DATABASES;",
		"--json")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "port", conn.Port, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	databases, err := m.parseSchemasOutput(string(output))
	if err != nil {
		log.Error("Error parsing database list output", "error", err)
		return nil, err
	}

	log.Info("Listed databases successfully", "host", conn.Host, "port", conn.Port, "count", len(databases))
	return databases, nil
}

func (m *MySQLShellTool) ListTables(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLShellExists(); err != nil {
		return nil, err
	}
	log.Info("Listing tables using mysqlshell", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	cmd := exec.Command("mysqlsh",
		"--uri", fmt.Sprintf("%s:%s@%s:%d/%s", conn.User, conn.Password, conn.Host, conn.Port, conn.Database),
		"--execute", "SHOW TABLES;",
		"--json")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	tables, err := m.parseTablesOutput(string(output))
	if err != nil {
		log.Error("Error parsing table list output", "error", err)
		return nil, err
	}

	log.Info("Listed tables successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database, "count", len(tables))
	return tables, nil
}

// 爬取数据库信息
func (m *MySQLShellTool) parseSchemasOutput(output string) ([]string, error) {
	var results []string
	if !strings.Contains(output, `"hasData": true`) {
		log.Error("查询结果集为空：", output)
	}

	rows := strings.SplitSeq(output, "\n")

	for v := range rows {
		if strings.Contains(v, `"Database"`) {
			parts := strings.Split(v, `"`)
			results = append(results, parts[3])
		}

	}

	return results, nil
}

// 爬取数据库表
func (m *MySQLShellTool) parseTablesOutput(output string) ([]string, error) {
	var results []string
	if !strings.Contains(output, `"hasData": true`) {
		log.Error("查询结果集为空：", output)
	}

	rows := strings.SplitSeq(output, "\n")

	for v := range rows {
		fmt.Println(v)
		if strings.Contains(v, `"Tables_in_`) {
			parts := strings.Split(v, `"`)
			results = append(results, parts[3])
		}

	}

	return results, nil
}

func (m *MySQLShellTool) buildExportArgs(conn DatabaseConnection, config ExportConfig, exportType string) []string {
	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--", "util", "export-instance",
	}

	if config.OutputDir != "" {
		args = append(args, config.OutputDir)
	}

	optionArgs := []string{}

	if config.Threads > 0 {
		optionArgs = append(optionArgs, fmt.Sprintf("threads=%d", config.Threads))
	}

	if config.ChunkSize != "" {
		optionArgs = append(optionArgs, fmt.Sprintf("chunkSize=%s", config.ChunkSize))
	}

	switch strings.ToLower(config.Compression) {
	case "gzip", "gz":
		optionArgs = append(optionArgs, "compression=gzip")
	case "zstd":
		optionArgs = append(optionArgs, "compression=zstd")
	case "none":
		optionArgs = append(optionArgs, "compression=none")
	}

	if len(config.IncludeSchemas) > 0 {
		optionArgs = append(optionArgs, fmt.Sprintf("includeSchemas=[\"%s\"]", strings.Join(config.IncludeSchemas, "\",\"")))
	}

	if len(config.ExcludeSchemas) > 0 {
		optionArgs = append(optionArgs, fmt.Sprintf("excludeSchemas=[\"%s\"]", strings.Join(config.ExcludeSchemas, "\",\"")))
	}

	if len(config.IncludeTables) > 0 {
		optionArgs = append(optionArgs, fmt.Sprintf("includeTables=[\"%s\"]", strings.Join(config.IncludeTables, "\",\"")))
	}

	if len(config.ExcludeTables) > 0 {
		optionArgs = append(optionArgs, fmt.Sprintf("excludeTables=[\"%s\"]", strings.Join(config.ExcludeTables, "\",\"")))
	}

	if len(optionArgs) > 0 {
		args = append(args, "--outputUrl")
		for _, opt := range optionArgs {
			args = append(args, opt)
		}
	}

	return args
}

func (m *MySQLShellTool) ExportDatabase(conn DatabaseConnection, database string, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	outputPath := filepath.Join(config.OutputDir, database)
	normalizedPath := filepath.ToSlash(outputPath)

	if config.Overwrite {
		if _, err := os.Stat(outputPath); err == nil {
			log.Info("Removing existing output directory for overwrite", "path", outputPath)
			if err := os.RemoveAll(outputPath); err != nil {
				log.Error("Failed to remove existing directory", "path", outputPath, "error", err)
				return fmt.Errorf("删除已存在目录失败: %s", err.Error())
			}
		}
	}

	log.Info("Exporting database", "database", database, "output", outputPath, "threads", config.Threads)

	options := map[string]any{}

	if config.Threads > 0 {
		options["threads"] = config.Threads
	}

	if config.ChunkSize != "" {
		options["chunkSize"] = config.ChunkSize
	}

	switch strings.ToLower(config.Compression) {
	case "gzip", "gz":
		options["compression"] = "gzip"
	case "zstd":
		options["compression"] = "zstd"
	case "none":
		options["compression"] = "none"
	}

	optionsJSON, _ := json.Marshal(options)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--js", "--execute", fmt.Sprintf(`util.dumpSchemas(['%s'], '%s', %s)`, database, normalizedPath, string(optionsJSON)),
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error exporting database", "database", database, "error", err, "output", string(output))
		return fmt.Errorf("导出数据库 %s 失败: %s, 错误信息: %s", database, err.Error(), string(output))
	}

	log.Info("Exported database successfully", "database", database, "output", outputPath)
	return nil
}

func (m *MySQLShellTool) ExportDatabases(conn DatabaseConnection, databases []string, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if config.Overwrite {
		if _, err := os.Stat(config.OutputDir); err == nil {
			log.Info("Removing existing output directory for overwrite", "path", config.OutputDir)
			if err := os.RemoveAll(config.OutputDir); err != nil {
				log.Error("Failed to remove existing directory", "path", config.OutputDir, "error", err)
				return fmt.Errorf("删除已存在目录失败: %s", err.Error())
			}
		}
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	normalizedOutputDir := filepath.ToSlash(config.OutputDir)

	log.Info("Exporting multiple databases", "databases", databases, "output", config.OutputDir, "threads", config.Threads)

	options := map[string]any{}

	if config.Threads > 0 {
		options["threads"] = config.Threads
	}

	if config.ChunkSize != "" {
		options["chunkSize"] = config.ChunkSize
	}

	switch strings.ToLower(config.Compression) {
	case "gzip", "gz":
		options["compression"] = "gzip"
	case "zstd":
		options["compression"] = "zstd"
	case "none":
		options["compression"] = "none"
	}

	databasesJSON, _ := json.Marshal(databases)
	optionsJSON, _ := json.Marshal(options)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--js", "--execute", fmt.Sprintf("util.dumpSchemas(%s, '%s', %s)", string(databasesJSON), normalizedOutputDir, string(optionsJSON)),
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error exporting databases", "databases", databases, "error", err, "output", string(output))
		return fmt.Errorf("批量导出数据库失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Exported databases successfully", "databases", databases, "output", config.OutputDir)
	return nil
}

func (m *MySQLShellTool) ExportTables(conn DatabaseConnection, database string, tables []string, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	outputPath := filepath.Join(config.OutputDir, database)

	if config.Overwrite {
		if _, err := os.Stat(outputPath); err == nil {
			log.Info("Removing existing output directory for overwrite", "path", outputPath)
			if err := os.RemoveAll(outputPath); err != nil {
				log.Error("Failed to remove existing directory", "path", outputPath, "error", err)
				return fmt.Errorf("删除已存在目录失败: %s", err.Error())
			}
		}
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	log.Info("Exporting tables", "database", database, "tables", tables, "output", outputPath, "threads", config.Threads)

	options := map[string]any{}

	if config.Threads > 0 {
		options["threads"] = config.Threads
	}

	if config.ChunkSize != "" {
		options["chunkSize"] = config.ChunkSize
	}

	switch strings.ToLower(config.Compression) {
	case "gzip", "gz":
		options["compression"] = "gzip"
	case "zstd":
		options["compression"] = "zstd"
	case "none":
		options["compression"] = "none"
	}

	tablesJSON, _ := json.Marshal(tables)
	optionsJSON, _ := json.Marshal(options)

	normalizedPath := filepath.ToSlash(outputPath)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--js", "--execute", fmt.Sprintf("util.dumpTables('%s', %s, '%s', %s)", database, string(tablesJSON), normalizedPath, string(optionsJSON)),
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error exporting tables", "database", database, "tables", tables, "error", err, "output", string(output))
		return fmt.Errorf("导出表失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Exported tables successfully", "database", database, "tables", tables, "output", outputPath)
	return nil
}

func (m *MySQLShellTool) ExportInstance(conn DatabaseConnection, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if config.Overwrite {
		if _, err := os.Stat(config.OutputDir); err == nil {
			log.Info("Removing existing output directory for overwrite", "path", config.OutputDir)
			if err := os.RemoveAll(config.OutputDir); err != nil {
				log.Error("Failed to remove existing directory", "path", config.OutputDir, "error", err)
				return fmt.Errorf("删除已存在目录失败: %s", err.Error())
			}
		}
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Error("Failed to create output directory", "path", config.OutputDir, "error", err)
		return fmt.Errorf("创建导出目录失败: %s", err.Error())
	}

	log.Info("Exporting entire instance", "host", conn.Host, "port", conn.Port, "output", config.OutputDir, "threads", config.Threads)

	options := map[string]any{}

	if config.Threads > 0 {
		options["threads"] = config.Threads
	}

	if config.ChunkSize != "" {
		options["chunkSize"] = config.ChunkSize
	}

	switch strings.ToLower(config.Compression) {
	case "gzip", "gz":
		options["compression"] = "gzip"
	case "zstd":
		options["compression"] = "zstd"
	case "none":
		options["compression"] = "none"
	}

	if len(config.IncludeSchemas) > 0 {
		options["includeSchemas"] = config.IncludeSchemas
	}

	if len(config.ExcludeSchemas) > 0 {
		options["excludeSchemas"] = config.ExcludeSchemas
	}

	optionsJSON, _ := json.Marshal(options)

	normalizedOutputDir := filepath.ToSlash(config.OutputDir)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--js", "--execute", fmt.Sprintf("util.dumpInstance('%s', %s)", normalizedOutputDir, string(optionsJSON)),
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error exporting instance", "host", conn.Host, "port", conn.Port, "error", err, "output", string(output))
		return fmt.Errorf("导出实例失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Exported instance successfully", "host", conn.Host, "port", conn.Port, "output", config.OutputDir)
	return nil
}
