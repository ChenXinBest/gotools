package tools

import (
	"errors"
	"gotools/internal/config"
	"gotools/internal/log"
	"time"
)

// DatabaseConnection 与 config 包中的结构保持一致
type DatabaseConnection = config.DatabaseConnection

// GetDatabaseConnections 获取所有数据库连接信息
func GetDatabaseConnections() ([]DatabaseConnection, error) {
	log.Info("Loading all database connections from config")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", "error", err)
		return nil, err
	}
	log.Info("Loaded database connections", "count", len(cfg.DatabaseConnections))
	return cfg.DatabaseConnections, nil
}

// GetDatabaseConnection 根据 ID 获取数据库连接信息
func GetDatabaseConnection(id string) (DatabaseConnection, error) {
	log.Info("Loading database connection by ID", "id", id)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", "error", err)
		return DatabaseConnection{}, err
	}

	for _, conn := range cfg.DatabaseConnections {
		if conn.ID == id {
			log.Info("Found database connection", "id", id, "name", conn.Name)
			return conn, nil
		}
	}

	log.Error("Database connection not found", "id", id)
	return DatabaseConnection{}, errors.New("connection not found")
}

// AddDatabaseConnection 添加数据库连接信息
func AddDatabaseConnection(conn DatabaseConnection) error {
	log.Info("Adding database connection", "name", conn.Name)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", "error", err)
		return err
	}

	// 生成唯一 ID
	if conn.ID == "" {
		conn.ID = generateID()
		log.Info("Generated ID for database connection", "id", conn.ID, "name", conn.Name)
	}

	// 检查是否已存在相同 ID 的连接
	for _, c := range cfg.DatabaseConnections {
		if c.ID == conn.ID {
			log.Error("Database connection with this ID already exists", "id", conn.ID)
			return errors.New("connection with this ID already exists")
		}
	}

	// 添加到配置中
	cfg.DatabaseConnections = append(cfg.DatabaseConnections, conn)

	// 保存配置
	err = config.SaveConfig(cfg)
	if err != nil {
		log.Error("Error saving config", "error", err)
		return err
	}
	log.Info("Added database connection successfully", "id", conn.ID, "name", conn.Name)
	return nil
}

// UpdateDatabaseConnection 更新数据库连接信息
func UpdateDatabaseConnection(conn DatabaseConnection) error {
	log.Info("Updating database connection", "id", conn.ID, "name", conn.Name)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", "error", err)
		return err
	}

	// 查找并更新连接
	found := false
	for i, c := range cfg.DatabaseConnections {
		if c.ID == conn.ID {
			cfg.DatabaseConnections[i] = conn
			found = true
			break
		}
	}

	if !found {
		log.Error("Database connection not found", "id", conn.ID)
		return errors.New("connection not found")
	}

	// 保存配置
	err = config.SaveConfig(cfg)
	if err != nil {
		log.Error("Error saving config", "error", err)
		return err
	}
	log.Info("Updated database connection successfully", "id", conn.ID, "name", conn.Name)
	return nil
}

// DeleteDatabaseConnection 删除数据库连接信息
func DeleteDatabaseConnection(id string) error {
	log.Info("Deleting database connection", "id", id)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", "error", err)
		return err
	}

	// 查找并删除连接
	newConnections := []DatabaseConnection{}
	for _, conn := range cfg.DatabaseConnections {
		if conn.ID != id {
			newConnections = append(newConnections, conn)
		}
	}

	// 如果连接数没有变化，说明没有找到要删除的连接
	if len(newConnections) == len(cfg.DatabaseConnections) {
		log.Error("Database connection not found", "id", id)
		return errors.New("connection not found")
	}

	// 更新配置
	cfg.DatabaseConnections = newConnections

	// 保存配置
	err = config.SaveConfig(cfg)
	if err != nil {
		log.Error("Error saving config", "error", err)
		return err
	}
	log.Info("Deleted database connection successfully", "id", id)
	return nil
}

// generateID 生成唯一 ID
func generateID() string {
	return time.Now().Format("20060102150405")
}
