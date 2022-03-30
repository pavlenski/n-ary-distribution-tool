package app

import (
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"time"
)

type config struct {
	FileInputSleepTime time.Duration
	//Disks - no need to use we can just pass them as an arg to the cli
	CounterDataLimit  uint64
	SortProgressLimit uint64
}

type App struct {
	inputHandlers map[string]*input.FileInput
}

func NewApp() *App {
	return &App{inputHandlers: make(map[string]*input.FileInput)}
}

func (a *App) Start() {
	//init configs

	a.run()
}
