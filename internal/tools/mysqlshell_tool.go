package tools

import (
	"fmt"
	"gotools/internal/log"
	"os/exec"
	"regexp"
	"strings"
)

// MySQLShellTool 提供 mysqlshell 相关功能
type MySQLShellTool struct{}

// NewMySQLShellTool 创建一个新的 MySQLShellTool 实例
func NewMySQLShellTool() *MySQLShellTool {
	log.Info("Creating new MySQLShellTool instance")
	return &MySQLShellTool{}
}

// ConnectDatabase 使用 mysqlshell 连接数据库
func (m *MySQLShellTool) ConnectDatabase(conn DatabaseConnection) error {
	log.Info("Connecting to database using mysqlshell", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	// 构建连接命令
	cmd := exec.Command("mysqlsh",
		"--uri", fmt.Sprintf("%s:%s@%s:%d/%s", conn.User, conn.Password, conn.Host, conn.Port, conn.Database),
		"--execute", "SELECT 1;")

	// 执行命令
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error connecting to database", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", string(output))
		return fmt.Errorf("连接失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	log.Info("Connected to database successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	return nil
}

// ListDatabases 查询数据库列表
func (m *MySQLShellTool) ListDatabases(conn DatabaseConnection) ([]string, error) {
	log.Info("Listing databases using mysqlshell", "host", conn.Host, "port", conn.Port)
	// 构建查询命令
	cmd := exec.Command("mysqlsh",
		"--uri", fmt.Sprintf("%s:%s@%s:%d", conn.User, conn.Password, conn.Host, conn.Port),
		"--execute", "SHOW DATABASES;",
		"--json")

	// 执行命令
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "port", conn.Port, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	// 解析输出
	databases, err := m.parseSchemasOutput(string(output))
	if err != nil {
		log.Error("Error parsing database list output", "error", err)
		return nil, err
	}

	log.Info("Listed databases successfully", "host", conn.Host, "port", conn.Port, "count", len(databases))
	return databases, nil
}

// ListTables 查询指定数据库下的表列表
func (m *MySQLShellTool) ListTables(conn DatabaseConnection) ([]string, error) {
	log.Info("Listing tables using mysqlshell", "host", conn.Host, "port", conn.Port, "database", conn.Database)
	// 构建查询命令
	cmd := exec.Command("mysqlsh",
		"--uri", fmt.Sprintf("%s:%s@%s:%d/%s", conn.User, conn.Password, conn.Host, conn.Port, conn.Database),
		"--execute", "SHOW TABLES;",
		"--json")

	// 执行命令
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "port", conn.Port, "database", conn.Database, "error", err, "output", string(output))
		return nil, fmt.Errorf("查询失败: %s, 错误信息: %s", err.Error(), string(output))
	}

	// 解析输出
	tables, err := m.parseJSONOutput(string(output))
	if err != nil {
		log.Error("Error parsing table list output", "error", err)
		return nil, err
	}

	log.Info("Listed tables successfully", "host", conn.Host, "port", conn.Port, "database", conn.Database, "count", len(tables))
	return tables, nil
}

// 爬取数据库信息
func (m *MySQLShellTool) parseSchemasOutput(output string) ([]string, error) {
	log.Info("开始提取：", output)
	reExp := "{*?hasData.*}"
	regex := regexp.MustCompile(reExp)

	matches := regex.FindStringSubmatch(output)
	// 匹配正则
	// if len(matches) < 1 {
	// 	return nil, fmt.Errorf("提取结果集出错！")
	// }

	log.Info("提取到结果集：", matches)

	return []string{}, nil
}

// parseJSONOutput 解析 mysqlshell 的 JSON 输出
func (m *MySQLShellTool) parseJSONOutput(output string) ([]string, error) {
	// 忽略警告信息，只处理 JSON 部分
	// 找到 JSON 开始的位置
	jsonStart := strings.Index(output, "{")
	if jsonStart == -1 {
		// 没有找到 JSON，返回空列表
		log.Info("No JSON found in output")
		return []string{}, nil
	}

	// 提取 JSON 部分
	jsonOutput := output[jsonStart:]

	// 简化解析，实际项目中可能需要使用 JSON 库
	// 这里假设输出格式为包含结果集的 JSON
	// 实际实现可能需要根据 mysqlshell 的具体输出格式进行调整

	// 这里只是一个简单的示例实现
	// 实际项目中应该使用更健壮的解析方法
	result := []string{}

	// 查找结果集开始的位置
	rowsStart := strings.Index(jsonOutput, "\"rows\":")
	if rowsStart == -1 {
		// 没有找到结果集，返回空列表
		log.Info("No rows found in JSON output")
		return []string{}, nil
	}

	// 查找结果集结束的位置
	rowsEnd := strings.Index(jsonOutput[rowsStart:], "}")
	if rowsEnd == -1 {
		// 没有找到结果集结束，返回空列表
		log.Info("No end of rows found in JSON output")
		return []string{}, nil
	}

	// 提取结果集部分
	rowsOutput := jsonOutput[rowsStart : rowsStart+rowsEnd+1]

	// 简单处理：提取引号中的内容
	parts := strings.Split(rowsOutput, "\"")
	for i := 1; i < len(parts)-1; i += 2 {
		value := parts[i]
		// 过滤掉系统字段和空值
		if value != "rows" && value != "Database" && value != "Tables_in_"+parts[0] && value != "" &&
			value != "executionTime" && value != "affectedItemsCount" && value != "warningsCount" &&
			value != "warnings" && value != "info" && value != "autoIncrementValue" {
			// 过滤掉时间值（如 "0.0021 sec"）
			if !strings.Contains(value, "sec") {
				result = append(result, value)
			}
		}
	}

	log.Info("Parsed JSON output successfully", "count", len(result))
	return result, nil
}
