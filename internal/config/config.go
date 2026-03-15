package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DatabaseConnection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ExportSettings struct {
	ExportTool       string   `json:"export_tool"`
	ExportPath       string   `json:"export_path"`
	LastConnectionID string   `json:"last_connection_id"`
	LastDatabases    []string `json:"last_databases"`
	LastDatabase     string   `json:"last_database"`
	LastTables       []string `json:"last_tables"`
	Threads          int      `json:"threads"`
	SkipDefiner      bool     `json:"skip_definer"`
	SkipBinlog       bool     `json:"skip_binlog"`
	Compression      string   `json:"compression"`
	ChunkSize        string   `json:"chunk_size"`
	ExportScope      string   `json:"export_scope"`
	IncludeSchemas   string   `json:"include_schemas"`
	ExcludeSchemas   string   `json:"exclude_schemas"`
	IncludeTables    string   `json:"include_tables"`
	ExcludeTables    string   `json:"exclude_tables"`
	Overwrite        bool     `json:"overwrite"`
}

type Config struct {
	DatabaseConnections []DatabaseConnection `json:"database_connections"`
	ExportSettings      ExportSettings       `json:"export_settings"`
}

func GetConfigPath() string {
	execDir, _ := os.Executable()
	return filepath.Join(filepath.Dir(execDir), "config.json")
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{
			DatabaseConnections: []DatabaseConnection{},
			ExportSettings: ExportSettings{
				ExportTool:  "mysql-shell",
				Threads:     4,
				Compression: "gzip",
				SkipDefiner: true,
				Overwrite:   true,
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

	if config.ExportSettings.Threads == 0 {
		config.ExportSettings.Threads = 4
	}
	if config.ExportSettings.ExportTool == "" {
		config.ExportSettings.ExportTool = "mysql-shell"
	}
	if config.ExportSettings.Compression == "" {
		config.ExportSettings.Compression = "gzip"
	}

	return &config, nil
}

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
