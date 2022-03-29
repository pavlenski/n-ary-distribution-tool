package app

import "time"

type config struct {
	FileInputSleepTime time.Duration
	//Disks - no need to use we can just pass them as an arg to the cli
	CounterDataLimit  uint64
	SortProgressLimit uint64
}

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() {
	//init configs

	a.run()
}
