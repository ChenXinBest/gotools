package services

import (
	"gotools/internal/log"
	"gotools/internal/tools"
)

// DialogService 对话框服务
type DialogService struct{}

// NewDialogService 创建对话框服务实例
func NewDialogService() *DialogService {
	log.Info("Creating new DialogService instance")
	return &DialogService{}
}

// SelectFolder 打开系统文件夹选择对话框并返回选择的路径
func (s *DialogService) SelectFolder() (string, error) {
	log.Info("Opening folder selection dialog")
	path, err := tools.SelectFolder()
	if err != nil {
		log.Error("Error opening folder dialog", "error", err)
		return "", err
	}
	log.Info("Folder selected", "path", path)
	return path, nil
}

// SelectFile 打开系统文件选择对话框并返回选择的文件路径
func (s *DialogService) SelectFile() (string, error) {
	log.Info("Opening file selection dialog")
	path, err := tools.SelectFile()
	if err != nil {
		log.Error("Error opening file dialog", "error", err)
		return "", err
	}
	log.Info("File selected", "path", path)
	return path, nil
}

// SelectSaveFile 打开系统保存文件对话框并返回选择的文件路径
func (s *DialogService) SelectSaveFile(defaultName string) (string, error) {
	log.Info("Opening save file dialog", "defaultName", defaultName)
	path, err := tools.SelectSaveFile(defaultName)
	if err != nil {
		log.Error("Error opening save file dialog", "error", err)
		return "", err
	}
	log.Info("Save file selected", "path", path)
	return path, nil
}
