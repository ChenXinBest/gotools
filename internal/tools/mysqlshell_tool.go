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

type DumpMetadata struct {
	Dumper  string   `json:"dumper"`
	Version string   `json:"version"`
	Origin  string   `json:"origin"`
	Schemas []string `json:"schemas"`
	Options struct {
		Compression string `json:"compression"`
		Threads     int    `json:"threads"`
	} `json:"options"`
}

type SchemaMetadata struct {
	Schema     string   `json:"schema"`
	Tables     []string `json:"tables"`
	Views      []string `json:"views"`
	Events     []string `json:"events"`
	Functions  []string `json:"functions"`
	Procedures []string `json:"procedures"`
}

type ImportConflict struct {
	Schema     string   `json:"schema"`
	Tables     []string `json:"tables"`
	Views      []string `json:"views"`
	Events     []string `json:"events"`
	Functions  []string `json:"functions"`
	Procedures []string `json:"procedures"`
}

type ImportConflictCheckResult struct {
	HasConflicts bool             `json:"has_conflicts"`
	Conflicts    []ImportConflict `json:"conflicts"`
}

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

type ImportConfig struct {
	InputDir       string
	Threads        int
	Schema         string
	IncludeSchemas []string
	ExcludeSchemas []string
	IncludeTables  []string
	ExcludeTables  []string
	ResetProgress  bool
	WaitTimeout    int
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

func (m *MySQLShellTool) enableLocalInfile(conn DatabaseConnection) error {
	log.Info("Enabling local_infile on server", "host", conn.Host, "port", conn.Port)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--sql", "--execute", "SET GLOBAL local_infile = ON;",
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error enabling local_infile", "error", err, "output", string(output))
		return fmt.Errorf("启用 local_infile 失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Enabled local_infile successfully", "host", conn.Host, "port", conn.Port)
	return nil
}

func (m *MySQLShellTool) ImportDatabases(conn DatabaseConnection, config ImportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if config.InputDir == "" {
		return fmt.Errorf("导入目录不能为空")
	}

	if _, err := os.Stat(config.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("导入目录不存在: %s", config.InputDir)
	}

	normalizedInputDir := filepath.ToSlash(config.InputDir)

	log.Info("Importing databases", "input", config.InputDir, "threads", config.Threads)

	options := map[string]any{}

	if config.Threads > 0 {
		options["threads"] = config.Threads
	}

	if config.Schema != "" {
		options["schema"] = config.Schema
	}

	if len(config.IncludeSchemas) > 0 {
		options["includeSchemas"] = config.IncludeSchemas
	}

	if len(config.ExcludeSchemas) > 0 {
		options["excludeSchemas"] = config.ExcludeSchemas
	}

	if len(config.IncludeTables) > 0 {
		options["includeTables"] = config.IncludeTables
	}

	if len(config.ExcludeTables) > 0 {
		options["excludeTables"] = config.ExcludeTables
	}

	if config.ResetProgress {
		options["resetProgress"] = true
	}

	if config.WaitTimeout > 0 {
		options["waitTimeout"] = config.WaitTimeout
	}

	optionsJSON, _ := json.Marshal(options)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--js", "--execute", fmt.Sprintf("util.loadDump('%s', %s)", normalizedInputDir, string(optionsJSON)),
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "local_infile disabled in server") || strings.Contains(outputStr, "local_infile' global system variable must be set to ON") {
			log.Info("Detected local_infile disabled, attempting to enable it", "host", conn.Host, "port", conn.Port)
			if enableErr := m.enableLocalInfile(conn); enableErr != nil {
				log.Error("Failed to enable local_infile", "error", enableErr)
				return fmt.Errorf("导入数据库失败: 服务器 local_infile 已禁用，尝试启用失败: %s", enableErr.Error())
			}
			log.Info("Retrying import after enabling local_infile", "input", config.InputDir)
			cmd = exec.Command("mysqlsh", args...)
			output, err = cmd.CombinedOutput()
			if err != nil {
				log.Error("Error importing databases after enabling local_infile", "input", config.InputDir, "error", err, "output", string(output))
				return fmt.Errorf("导入数据库失败: %s, 错误信息: %s", err.Error(), string(output))
			}
		} else {
			log.Error("Error importing databases", "input", config.InputDir, "error", err, "output", outputStr)
			return fmt.Errorf("导入数据库失败: %s, 错误信息: %s", err.Error(), outputStr)
		}
	}

	log.Info("Imported databases successfully", "input", config.InputDir)
	return nil
}

func (m *MySQLShellTool) ImportTables(conn DatabaseConnection, database string, config ImportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if config.InputDir == "" {
		return fmt.Errorf("导入目录不能为空")
	}

	if _, err := os.Stat(config.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("导入目录不存在: %s", config.InputDir)
	}

	if database == "" {
		return fmt.Errorf("目标数据库不能为空")
	}

	normalizedInputDir := filepath.ToSlash(config.InputDir)

	log.Info("Importing tables", "database", database, "input", config.InputDir, "threads", config.Threads)

	options := map[string]any{}

	if config.Threads > 0 {
		options["threads"] = config.Threads
	}

	options["schema"] = database

	if len(config.IncludeTables) > 0 {
		options["includeTables"] = config.IncludeTables
	}

	if len(config.ExcludeTables) > 0 {
		options["excludeTables"] = config.ExcludeTables
	}

	if config.ResetProgress {
		options["resetProgress"] = true
	}

	if config.WaitTimeout > 0 {
		options["waitTimeout"] = config.WaitTimeout
	}

	optionsJSON, _ := json.Marshal(options)

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--js", "--execute", fmt.Sprintf("util.loadDump('%s', %s)", normalizedInputDir, string(optionsJSON)),
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "local_infile disabled in server") || strings.Contains(outputStr, "local_infile' global system variable must be set to ON") {
			log.Info("Detected local_infile disabled, attempting to enable it", "host", conn.Host, "port", conn.Port)
			if enableErr := m.enableLocalInfile(conn); enableErr != nil {
				log.Error("Failed to enable local_infile", "error", enableErr)
				return fmt.Errorf("导入表失败: 服务器 local_infile 已禁用，尝试启用失败: %s", enableErr.Error())
			}
			log.Info("Retrying import after enabling local_infile", "database", database, "input", config.InputDir)
			cmd = exec.Command("mysqlsh", args...)
			output, err = cmd.CombinedOutput()
			if err != nil {
				log.Error("Error importing tables after enabling local_infile", "database", database, "input", config.InputDir, "error", err, "output", string(output))
				return fmt.Errorf("导入表失败: %s, 错误信息: %s", err.Error(), string(output))
			}
		} else {
			log.Error("Error importing tables", "database", database, "input", config.InputDir, "error", err, "output", outputStr)
			return fmt.Errorf("导入表失败: %s, 错误信息: %s", err.Error(), outputStr)
		}
	}

	log.Info("Imported tables successfully", "database", database, "input", config.InputDir)
	return nil
}

func (m *MySQLShellTool) ParseDumpMetadata(inputDir string) (*DumpMetadata, map[string]*SchemaMetadata, error) {
	metadataFile := filepath.Join(inputDir, "@.json")
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		log.Error("Failed to read dump metadata file", "path", metadataFile, "error", err)
		return nil, nil, fmt.Errorf("读取元数据文件失败: %s", err.Error())
	}

	var dumpMeta DumpMetadata
	if err := json.Unmarshal(data, &dumpMeta); err != nil {
		log.Error("Failed to parse dump metadata", "error", err)
		return nil, nil, fmt.Errorf("解析元数据失败: %s", err.Error())
	}

	log.Info("Parsed dump metadata", "schemas", dumpMeta.Schemas, "origin", dumpMeta.Origin)

	schemaMetas := make(map[string]*SchemaMetadata)
	for _, schema := range dumpMeta.Schemas {
		schemaFile := filepath.Join(inputDir, schema+".json")
		schemaData, err := os.ReadFile(schemaFile)
		if err != nil {
			log.Warn("Schema metadata file not found, skipping", "schema", schema, "path", schemaFile)
			continue
		}

		var schemaMeta SchemaMetadata
		if err := json.Unmarshal(schemaData, &schemaMeta); err != nil {
			log.Warn("Failed to parse schema metadata", "schema", schema, "error", err)
			continue
		}

		schemaMetas[schema] = &schemaMeta
		log.Info("Parsed schema metadata", "schema", schema, "tables", len(schemaMeta.Tables))
	}

	return &dumpMeta, schemaMetas, nil
}

func (m *MySQLShellTool) CheckExistingObjects(conn DatabaseConnection, schemaMetas map[string]*SchemaMetadata) (*ImportConflictCheckResult, error) {
	if err := m.checkMySQLShellExists(); err != nil {
		return nil, err
	}

	result := &ImportConflictCheckResult{
		HasConflicts: false,
		Conflicts:    []ImportConflict{},
	}

	for schemaName, schemaMeta := range schemaMetas {
		conflict := ImportConflict{
			Schema:     schemaName,
			Tables:     []string{},
			Views:      []string{},
			Events:     []string{},
			Functions:  []string{},
			Procedures: []string{},
		}

		if len(schemaMeta.Tables) > 0 {
			existingTables, err := m.getExistingTables(conn, schemaName, schemaMeta.Tables)
			if err != nil {
				log.Warn("Failed to check existing tables", "schema", schemaName, "error", err)
			} else {
				conflict.Tables = existingTables
			}
		}

		if len(schemaMeta.Views) > 0 {
			existingViews, err := m.getExistingViews(conn, schemaName, schemaMeta.Views)
			if err != nil {
				log.Warn("Failed to check existing views", "schema", schemaName, "error", err)
			} else {
				conflict.Views = existingViews
			}
		}

		if len(conflict.Tables) > 0 || len(conflict.Views) > 0 || len(conflict.Events) > 0 || len(conflict.Functions) > 0 || len(conflict.Procedures) > 0 {
			result.HasConflicts = true
			result.Conflicts = append(result.Conflicts, conflict)
		}
	}

	log.Info("Checked existing objects", "has_conflicts", result.HasConflicts, "conflict_count", len(result.Conflicts))
	return result, nil
}

func (m *MySQLShellTool) getExistingTables(conn DatabaseConnection, schema string, tables []string) ([]string, error) {
	if len(tables) == 0 {
		return []string{}, nil
	}

	tableList := make([]string, len(tables))
	for i, t := range tables {
		tableList[i] = fmt.Sprintf("'%s'", t)
	}

	query := fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME IN (%s) AND TABLE_TYPE = 'BASE TABLE'", schema, strings.Join(tableList, ","))

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--sql", "--execute", query,
		"--json",
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to query existing tables", "schema", schema, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询已存在表失败: %s", err.Error())
	}

	return m.parseExistingObjectsOutput(string(output), "TABLE_NAME"), nil
}

func (m *MySQLShellTool) getExistingViews(conn DatabaseConnection, schema string, views []string) ([]string, error) {
	if len(views) == 0 {
		return []string{}, nil
	}

	viewList := make([]string, len(views))
	for i, v := range views {
		viewList[i] = fmt.Sprintf("'%s'", v)
	}

	query := fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME IN (%s)", schema, strings.Join(viewList, ","))

	args := []string{
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--sql", "--execute", query,
		"--json",
	}

	cmd := exec.Command("mysqlsh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to query existing views", "schema", schema, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询已存在视图失败: %s", err.Error())
	}

	return m.parseExistingObjectsOutput(string(output), "TABLE_NAME"), nil
}

func (m *MySQLShellTool) parseExistingObjectsOutput(output, fieldName string) []string {
	var results []string
	rows := strings.SplitSeq(output, "\n")

	for v := range rows {
		if strings.Contains(v, fmt.Sprintf(`"%s"`, fieldName)) {
			parts := strings.Split(v, `"`)
			if len(parts) >= 4 {
				results = append(results, parts[3])
			}
		}
	}

	return results
}

func (m *MySQLShellTool) DropObjects(conn DatabaseConnection, conflicts []ImportConflict) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	for _, conflict := range conflicts {
		var dropStatements []string

		for _, table := range conflict.Tables {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP TABLE IF EXISTS `%s`.`%s`;", conflict.Schema, table))
		}

		for _, view := range conflict.Views {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP VIEW IF EXISTS `%s`.`%s`;", conflict.Schema, view))
		}

		for _, event := range conflict.Events {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP EVENT IF EXISTS `%s`.`%s`;", conflict.Schema, event))
		}

		for _, fn := range conflict.Functions {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP FUNCTION IF EXISTS `%s`.`%s`;", conflict.Schema, fn))
		}

		for _, proc := range conflict.Procedures {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP PROCEDURE IF EXISTS `%s`.`%s`;", conflict.Schema, proc))
		}

		if len(dropStatements) == 0 {
			continue
		}

		dropSQL := "SET FOREIGN_KEY_CHECKS = 0; " + strings.Join(dropStatements, " ") + " SET FOREIGN_KEY_CHECKS = 1;"

		log.Info("Dropping conflicting objects", "schema", conflict.Schema, "tables", len(conflict.Tables), "views", len(conflict.Views))

		args := []string{
			"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
			"--sql", "--execute", dropSQL,
		}

		cmd := exec.Command("mysqlsh", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("Failed to drop objects", "schema", conflict.Schema, "error", err, "output", string(output))
			return fmt.Errorf("删除冲突对象失败: %s, 错误信息: %s", err.Error(), string(output))
		}

		log.Info("Dropped conflicting objects successfully", "schema", conflict.Schema)
	}

	return nil
}
