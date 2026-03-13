package main

import (
	"context"
	"fmt"

	"gotools/internal/tools"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetSystemProcessInfos() ([]tools.ProcessInfo, error) {
	return tools.GetSystemProcessInfos()
}

func (a *App) SearchPidByKeyWord(keyword string) (tools.ProcessInfo, error) {
	return tools.SearchPidByKeyWord(keyword)
}

func (a *App) KillProcessByPID(pid int32) error {
	return tools.KillProcessByPID(pid)
}
