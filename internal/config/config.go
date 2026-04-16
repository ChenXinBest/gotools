package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DatabaseConnection 数据库连接配置
type DatabaseConnection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// MySQLShellConfig MySQL Shell 导出配置
type MySQLShellConfig struct {
	Threads     int    `json:"threads"`
	Compression string `json:"compression"`
	ChunkSize   string `json:"chunk_size"`
	SkipDefiner bool   `json:"skip_definer"`
	SkipBinlog  bool   `json:"skip_binlog"`
	Overwrite   bool   `json:"overwrite"`
}

// MySQLDumpConfig mysqldump 导出配置
type MySQLDumpConfig struct {
	Compression       string `json:"compression"`
	SingleTransaction bool   `json:"single_transaction"`
	Routines          bool   `json:"routines"`
	Events            bool   `json:"events"`
	Overwrite         bool   `json:"overwrite"`
}

// ExportSettings 导出设置
type ExportSettings struct {
	ExportTool       string           `json:"export_tool"`
	ExportPath       string           `json:"export_path"`
	LastConnectionID string           `json:"last_connection_id"`
	LastDatabases    []string         `json:"last_databases"`
	LastDatabase     string           `json:"last_database"`
	LastTables       []string         `json:"last_tables"`
	MySQLShell       MySQLShellConfig `json:"mysql_shell"`
	MySQLDump        MySQLDumpConfig  `json:"mysql_dump"`
}

// Config 应用配置结构
type Config struct {
	DatabaseConnections []DatabaseConnection `json:"database_connections"`
	ExportSettings      ExportSettings       `json:"export_settings"`
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	execDir, _ := os.Executable()
	return filepath.Join(filepath.Dir(execDir), "config.json")
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{
			DatabaseConnections: []DatabaseConnection{},
			ExportSettings: ExportSettings{
				ExportTool: "mysql-shell",
				MySQLShell: MySQLShellConfig{
					Threads:     4,
					Compression: "gzip",
					ChunkSize:   "64M",
					SkipDefiner: true,
					SkipBinlog:  false,
					Overwrite:   true,
				},
				MySQLDump: MySQLDumpConfig{
					Compression:       "gzip",
					SingleTransaction: true,
					Routines:          true,
					Events:            true,
					Overwrite:         true,
				},
			},
		}
		err = SaveConfig(defaultConfig)
		if err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	if config.ExportSettings.ExportTool == "" {
		config.ExportSettings.ExportTool = "mysql-shell"
	}

	// 确保 MySQLShell 配置有默认值
	if config.ExportSettings.MySQLShell.Threads == 0 {
		config.ExportSettings.MySQLShell.Threads = 4
	}
	if config.ExportSettings.MySQLShell.Compression == "" {
		config.ExportSettings.MySQLShell.Compression = "gzip"
	}
	if config.ExportSettings.MySQLShell.ChunkSize == "" {
		config.ExportSettings.MySQLShell.ChunkSize = "64M"
	}

	// 确保 MySQLDump 配置有默认值
	if config.ExportSettings.MySQLDump.Compression == "" {
		config.ExportSettings.MySQLDump.Compression = "gzip"
	}

	return &config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config) error {
	configPath := GetConfigPath()

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return err
	}

	return nil
}
