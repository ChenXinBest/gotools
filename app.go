package main

import (
	"context"
	"fmt"

	"gotools/internal/log"
	"gotools/internal/tools"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	log.Info("Creating new app instance")
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Info("App startup")
}

func (a *App) Greet(name string) string {
	log.Info("Greeting user", "name", name)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetSystemProcessInfos() ([]tools.ProcessInfo, error) {
	log.Info("Getting system process infos")
	processes, err := tools.GetSystemProcessInfos()
	if err != nil {
		log.Error("Error getting system process infos", "error", err)
		return nil, err
	}
	log.Info("Got system process infos", "count", len(processes))
	return processes, nil
}

func (a *App) SearchPidByKeyWord(keyword string) (tools.ProcessInfo, error) {
	log.Info("Searching process by keyword", "keyword", keyword)
	process, err := tools.SearchPidByKeyWord(keyword)
	if err != nil {
		log.Error("Error searching process by keyword", "keyword", keyword, "error", err)
		return tools.ProcessInfo{}, err
	}
	log.Info("Found process by keyword", "keyword", keyword, "pid", process.PID)
	return process, nil
}

func (a *App) KillProcessByPID(pid int32) error {
	log.Info("Killing process by PID", "pid", pid)
	err := tools.KillProcessByPID(pid)
	if err != nil {
		log.Error("Error killing process", "pid", pid, "error", err)
		return err
	}
	log.Info("Killed process successfully", "pid", pid)
	return nil
}

// GetDatabaseConnections 获取所有数据库连接信息
func (a *App) GetDatabaseConnections() ([]tools.DatabaseConnection, error) {
	log.Info("Getting all database connections")
	connections, err := tools.GetDatabaseConnections()
	if err != nil {
		log.Error("Error getting database connections", "error", err)
		return nil, err
	}
	log.Info("Got database connections", "count", len(connections))
	return connections, nil
}

// GetDatabaseConnection 根据 ID 获取数据库连接信息
func (a *App) GetDatabaseConnection(id string) (tools.DatabaseConnection, error) {
	log.Info("Getting database connection by ID", "id", id)
	connection, err := tools.GetDatabaseConnection(id)
	if err != nil {
		log.Error("Error getting database connection", "id", id, "error", err)
		return tools.DatabaseConnection{}, err
	}
	log.Info("Got database connection", "id", id, "name", connection.Name)
	return connection, nil
}

// AddDatabaseConnection 添加数据库连接信息
func (a *App) AddDatabaseConnection(conn tools.DatabaseConnection) error {
	log.Info("Adding database connection", "name", conn.Name)
	err := tools.AddDatabaseConnection(conn)
	if err != nil {
		log.Error("Error adding database connection", "name", conn.Name, "error", err)
		return err
	}
	log.Info("Added database connection successfully", "name", conn.Name)
	return nil
}

// UpdateDatabaseConnection 更新数据库连接信息
func (a *App) UpdateDatabaseConnection(conn tools.DatabaseConnection) error {
	log.Info("Updating database connection", "id", conn.ID, "name", conn.Name)
	err := tools.UpdateDatabaseConnection(conn)
	if err != nil {
		log.Error("Error updating database connection", "id", conn.ID, "name", conn.Name, "error", err)
		return err
	}
	log.Info("Updated database connection successfully", "id", conn.ID, "name", conn.Name)
	return nil
}

// DeleteDatabaseConnection 删除数据库连接信息
func (a *App) DeleteDatabaseConnection(id string) error {
	log.Info("Deleting database connection", "id", id)
	err := tools.DeleteDatabaseConnection(id)
	if err != nil {
		log.Error("Error deleting database connection", "id", id, "error", err)
		return err
	}
	log.Info("Deleted database connection successfully", "id", id)
	return nil
}

// ConnectDatabase 使用 mysqlshell 连接数据库
func (a *App) ConnectDatabase(conn tools.DatabaseConnection) error {
	log.Info("Connecting to database", "name", conn.Name, "host", conn.Host)
	mysqlShellTool := tools.NewMySQLShellTool()
	err := mysqlShellTool.ConnectDatabase(conn)
	if err != nil {
		log.Error("Error connecting to database", "name", conn.Name, "host", conn.Host, "error", err)
		return err
	}
	log.Info("Connected to database successfully", "name", conn.Name, "host", conn.Host)
	return nil
}

// ListDatabases 查询数据库列表
func (a *App) ListDatabases(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing databases", "host", conn.Host)
	mysqlShellTool := tools.NewMySQLShellTool()
	databases, err := mysqlShellTool.ListDatabases(conn)
	if err != nil {
		log.Error("Error listing databases", "host", conn.Host, "error", err)
		return nil, err
	}
	log.Info("Listed databases successfully", "host", conn.Host, "count", len(databases))
	return databases, nil
}

// ListTables 查询指定数据库下的表列表
func (a *App) ListTables(conn tools.DatabaseConnection) ([]string, error) {
	log.Info("Listing tables", "host", conn.Host, "database", conn.Database)
	mysqlShellTool := tools.NewMySQLShellTool()
	tables, err := mysqlShellTool.ListTables(conn)
	if err != nil {
		log.Error("Error listing tables", "host", conn.Host, "database", conn.Database, "error", err)
		return nil, err
	}
	log.Info("Listed tables successfully", "host", conn.Host, "database", conn.Database, "count", len(tables))
	return tables, nil
}

// SelectFolder 打开系统文件夹选择对话框
func (a *App) SelectFolder() (string, error) {
	return tools.SelectFolder()
}
