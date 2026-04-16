package tools

import (
	"encoding/json"
	"fmt"
	"gotools/internal/log"
	"os"
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

// escapeMySQLIdentifier 转义 MySQL 标识符（数据库名、表名、列名等）
// 使用反引号包裹并转义内部反引号
func escapeMySQLIdentifier(identifier string) string {
	// 先将内部的反引号替换为两个反引号
	escaped := strings.ReplaceAll(identifier, "`", "``")
	return "`" + escaped + "`"
}

type MySQLShellTool struct {
	common *CommonTool
}

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
	return &MySQLShellTool{
		common: NewCommonTool(),
	}
}

func (m *MySQLShellTool) checkMySQLShellExists() error {
	err := m.common.CheckCommandExists("mysqlsh")
	if err != nil {
		return fmt.Errorf("%s: 未找到 mysqlsh 命令，请先安装 MySQL Shell", ErrMySQLShellNotFound)
	}
	return nil
}

func (m *MySQLShellTool) ConnectDatabase(conn DatabaseConnection) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	log.Info("Connecting to database using mysqlshell", "host", conn.Host, "port", conn.Port, "database", conn.Database)

	uri := m.common.BuildDatabaseURI(conn, true)
	output, err := m.common.ExecuteCommand("mysqlsh",
		"--uri", uri,
		"--execute", "SELECT 1;")

	if err != nil {
		log.Error("Error connecting to database", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", output)
		return fmt.Errorf("连接失败: %s, 错误信息: %s", err.Error(), output)
	}

	log.Info("Connected to database successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	return nil
}

func (m *MySQLShellTool) ListDatabases(conn DatabaseConnection) ([]string, error) {
	if err := m.checkMySQLShellExists(); err != nil {
		return nil, err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return nil, err
	}

	log.Info("Listing databases using mysqlshell", "host", conn.Host, "port", conn.Port)

	uri := m.common.BuildDatabaseURI(conn, false)
	output, err := m.common.ExecuteCommand("mysqlsh",
		"--uri", uri,
		"--execute", "SHOW DATABASES;",
		"--json")

	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "port", conn.Port, "error", err, "output", output)
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), output)
	}

	databases, err := m.parseSchemasOutput(output)
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

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return nil, err
	}

	if conn.Database == "" {
		return nil, fmt.Errorf("数据库名称不能为空")
	}

	log.Info("Listing tables using mysqlshell", "host", conn.Host, "port", conn.Port, "database", conn.Database)

	uri := m.common.BuildDatabaseURI(conn, true)
	output, err := m.common.ExecuteCommand("mysqlsh",
		"--uri", uri,
		"--execute", "SHOW TABLES;",
		"--json")

	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", output)
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), output)
	}

	tables, err := m.parseTablesOutput(output)
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

// buildExportOptions 构建导出选项
func (m *MySQLShellTool) buildExportOptions(config ExportConfig) map[string]any {
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

	if len(config.IncludeTables) > 0 {
		options["includeTables"] = config.IncludeTables
	}

	if len(config.ExcludeTables) > 0 {
		options["excludeTables"] = config.ExcludeTables
	}

	return options
}

func (m *MySQLShellTool) ExportDatabase(conn DatabaseConnection, database string, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
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

	outputPath := filepath.Join(config.OutputDir, database)
	normalizedPath := m.common.NormalizePath(outputPath)

	if config.Overwrite {
		if err := m.common.RemoveIfExists(outputPath); err != nil {
			return err
		}
	}

	log.Info("Exporting database", "database", database, "output", outputPath, "threads", config.Threads)

	options := m.buildExportOptions(config)
	optionsJSON, _ := json.Marshal(options)

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--js", "--execute", fmt.Sprintf("util.dumpSchemas(['%s'], '%s', %s)", database, normalizedPath, string(optionsJSON)),
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Error exporting database", "database", database, "error", err, "output", output)
		return fmt.Errorf("导出数据库 %s 失败: %s, 错误信息: %s", database, err.Error(), output)
	}

	log.Info("Exported database successfully", "database", database, "output", outputPath)
	return nil
}

func (m *MySQLShellTool) ExportDatabases(conn DatabaseConnection, databases []string, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	if err := m.common.ValidateExportConfig(config.OutputDir); err != nil {
		return err
	}

	if len(databases) == 0 {
		return fmt.Errorf("数据库列表不能为空")
	}

	if err := m.common.EnsureDirExists(config.OutputDir); err != nil {
		return err
	}

	if config.Overwrite {
		if err := m.common.RemoveIfExists(config.OutputDir); err != nil {
			return err
		}
		// 重新创建目录
		if err := m.common.EnsureDirExists(config.OutputDir); err != nil {
			return err
		}
	}

	normalizedOutputDir := m.common.NormalizePath(config.OutputDir)

	log.Info("Exporting multiple databases", "databases", databases, "output", config.OutputDir, "threads", config.Threads)

	options := m.buildExportOptions(config)
	databasesJSON, _ := json.Marshal(databases)
	optionsJSON, _ := json.Marshal(options)

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--js", "--execute", fmt.Sprintf("util.dumpSchemas(%s, '%s', %s)", string(databasesJSON), normalizedOutputDir, string(optionsJSON)),
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Error exporting databases", "databases", databases, "error", err, "output", output)
		return fmt.Errorf("批量导出数据库失败: %s, 错误信息: %s", err.Error(), output)
	}

	log.Info("Exported databases successfully", "databases", databases, "output", config.OutputDir)
	return nil
}

func (m *MySQLShellTool) ExportTables(conn DatabaseConnection, database string, tables []string, config ExportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
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

	if len(tables) == 0 {
		return fmt.Errorf("表列表不能为空")
	}

	if err := m.common.EnsureDirExists(config.OutputDir); err != nil {
		return err
	}

	outputPath := filepath.Join(config.OutputDir, database)

	if config.Overwrite {
		if err := m.common.RemoveIfExists(outputPath); err != nil {
			return err
		}
	}

	log.Info("Exporting tables", "database", database, "tables", tables, "output", outputPath, "threads", config.Threads)

	options := m.buildExportOptions(config)
	tablesJSON, _ := json.Marshal(tables)
	optionsJSON, _ := json.Marshal(options)

	normalizedPath := m.common.NormalizePath(outputPath)

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--js", "--execute", fmt.Sprintf("util.dumpTables('%s', %s, '%s', %s)", database, string(tablesJSON), normalizedPath, string(optionsJSON)),
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Error exporting tables", "database", database, "tables", tables, "error", err, "output", output)
		return fmt.Errorf("导出表失败: %s, 错误信息: %s", err.Error(), output)
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

	normalizedOutputDir := m.common.NormalizePath(config.OutputDir)

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--js", "--execute", fmt.Sprintf("util.dumpInstance('%s', %s)", normalizedOutputDir, string(optionsJSON)),
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Error exporting instance", "host", conn.Host, "port", conn.Port, "error", err, "output", output)
		return fmt.Errorf("导出实例失败: %s, 错误信息: %s", err.Error(), output)
	}

	log.Info("Exported instance successfully", "host", conn.Host, "port", conn.Port, "output", config.OutputDir)
	return nil
}

func (m *MySQLShellTool) enableLocalInfile(conn DatabaseConnection) error {
	log.Info("Enabling local_infile on server", "host", conn.Host, "port", conn.Port)

	uri := m.common.BuildDatabaseURI(conn, false)

	globalArgs := []string{
		"--uri", uri,
		"--sql", "--execute", "SET GLOBAL local_infile = ON;",
	}

	output, err := m.common.ExecuteCommand("mysqlsh", globalArgs...)
	if err != nil {
		log.Warn("Failed to set GLOBAL local_infile, trying SESSION", "error", err, "output", output)
	}

	sessionArgs := []string{
		"--uri", uri,
		"--sql", "--execute", "SET SESSION local_infile = ON;",
	}

	sessionOutput, sessionErr := m.common.ExecuteCommand("mysqlsh", sessionArgs...)
	if sessionErr != nil {
		log.Error("Error enabling local_infile", "error", sessionErr, "output", sessionOutput)
		return fmt.Errorf("启用 local_infile 失败: %s, 错误信息: %s", sessionErr.Error(), sessionOutput)
	}

	log.Info("Enabled local_infile successfully", "host", conn.Host, "port", conn.Port)
	return nil
}

func (m *MySQLShellTool) ImportDatabases(conn DatabaseConnection, config ImportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	if config.InputDir == "" {
		return fmt.Errorf("导入目录不能为空")
	}

	if !m.common.FileExists(config.InputDir) {
		return fmt.Errorf("导入目录不存在: %s", config.InputDir)
	}

	normalizedInputDir := m.common.NormalizePath(config.InputDir)

	log.Info("Importing databases", "input", config.InputDir, "threads", config.Threads)

	if err := m.enableLocalInfile(conn); err != nil {
		log.Warn("Failed to enable local_infile, continuing anyway", "error", err)
	}

	options := m.buildImportOptions(config)
	optionsJSON, _ := json.Marshal(options)

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--js", "--execute", fmt.Sprintf("util.loadDump('%s', %s)", normalizedInputDir, string(optionsJSON)),
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Error importing databases", "input", config.InputDir, "error", err, "output", output)
		return fmt.Errorf("导入数据库失败: %s, 错误信息: %s", err.Error(), output)
	}

	log.Info("Imported databases successfully", "input", config.InputDir)
	return nil
}

// buildImportOptions 构建导入选项
func (m *MySQLShellTool) buildImportOptions(config ImportConfig) map[string]any {
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

	return options
}

func (m *MySQLShellTool) ImportTables(conn DatabaseConnection, database string, config ImportConfig) error {
	if err := m.checkMySQLShellExists(); err != nil {
		return err
	}

	if err := m.common.ValidateDatabaseConnection(conn); err != nil {
		return err
	}

	if config.InputDir == "" {
		return fmt.Errorf("导入目录不能为空")
	}

	if !m.common.FileExists(config.InputDir) {
		return fmt.Errorf("导入目录不存在: %s", config.InputDir)
	}

	if database == "" {
		return fmt.Errorf("目标数据库不能为空")
	}

	normalizedInputDir := m.common.NormalizePath(config.InputDir)

	log.Info("Importing tables", "database", database, "input", config.InputDir, "threads", config.Threads)

	if err := m.enableLocalInfile(conn); err != nil {
		log.Warn("Failed to enable local_infile, continuing anyway", "error", err)
	}

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

	optionsJSON, _ := json.Marshal(options)

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--js", "--execute", fmt.Sprintf("util.loadDump('%s', %s)", normalizedInputDir, string(optionsJSON)),
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Error importing tables", "database", database, "input", config.InputDir, "error", err, "output", output)
		return fmt.Errorf("导入表失败: %s, 错误信息: %s", err.Error(), output)
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
		tableList[i] = escapeMySQLIdentifier(t)
	}

	query := fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = %s AND TABLE_NAME IN (%s) AND TABLE_TYPE = 'BASE TABLE'", escapeMySQLIdentifier(schema), strings.Join(tableList, ","))

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--sql", "--execute", query,
		"--json",
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Failed to query existing tables", "schema", schema, "error", err, "output", output)
		return nil, fmt.Errorf("查询已存在表失败: %s", err.Error())
	}

	return m.parseExistingObjectsOutput(output, "TABLE_NAME"), nil
}

func (m *MySQLShellTool) getExistingViews(conn DatabaseConnection, schema string, views []string) ([]string, error) {
	if len(views) == 0 {
		return []string{}, nil
	}

	viewList := make([]string, len(views))
	for i, v := range views {
		viewList[i] = escapeMySQLIdentifier(v)
	}

	query := fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.VIEWS WHERE TABLE_SCHEMA = %s AND TABLE_NAME IN (%s)", escapeMySQLIdentifier(schema), strings.Join(viewList, ","))

	uri := m.common.BuildDatabaseURI(conn, false)
	args := []string{
		"--uri", uri,
		"--sql", "--execute", query,
		"--json",
	}

	output, err := m.common.ExecuteCommand("mysqlsh", args...)
	if err != nil {
		log.Error("Failed to query existing views", "schema", schema, "error", err, "output", output)
		return nil, fmt.Errorf("查询已存在视图失败: %s", err.Error())
	}

	return m.parseExistingObjectsOutput(output, "TABLE_NAME"), nil
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

		schemaIdentifier := escapeMySQLIdentifier(conflict.Schema)

		for _, table := range conflict.Tables {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP TABLE IF EXISTS %s.%s;", schemaIdentifier, escapeMySQLIdentifier(table)))
		}

		for _, view := range conflict.Views {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP VIEW IF EXISTS %s.%s;", schemaIdentifier, escapeMySQLIdentifier(view)))
		}

		for _, event := range conflict.Events {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP EVENT IF EXISTS %s.%s;", schemaIdentifier, escapeMySQLIdentifier(event)))
		}

		for _, fn := range conflict.Functions {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP FUNCTION IF EXISTS %s.%s;", schemaIdentifier, escapeMySQLIdentifier(fn)))
		}

		for _, proc := range conflict.Procedures {
			dropStatements = append(dropStatements, fmt.Sprintf("DROP PROCEDURE IF EXISTS %s.%s;", schemaIdentifier, escapeMySQLIdentifier(proc)))
		}

		if len(dropStatements) == 0 {
			continue
		}

		dropSQL := "SET FOREIGN_KEY_CHECKS = 0; " + strings.Join(dropStatements, " ") + " SET FOREIGN_KEY_CHECKS = 1;"

		log.Info("Dropping conflicting objects", "schema", conflict.Schema, "tables", len(conflict.Tables), "views", len(conflict.Views))

		uri := m.common.BuildDatabaseURI(conn, false)
		args := []string{
			"--uri", uri,
			"--sql", "--execute", dropSQL,
		}

		output, err := m.common.ExecuteCommand("mysqlsh", args...)
		if err != nil {
			log.Error("Failed to drop objects", "schema", conflict.Schema, "error", err, "output", output)
			return fmt.Errorf("删除冲突对象失败: %s, 错误信息: %s", err.Error(), output)
		}

		log.Info("Dropped conflicting objects successfully", "schema", conflict.Schema)
	}

	return nil
}
