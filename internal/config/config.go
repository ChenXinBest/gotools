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

type Config struct {
	DatabaseConnections []DatabaseConnection `json:"database_connections"`
}

func GetConfigPath() string {
	execDir, _ := os.Executable()
	return filepath.Join(filepath.Dir(execDir), "config.json")
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	
	// If config file doesn't exist, create a default one
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{
			DatabaseConnections: []DatabaseConnection{},
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
